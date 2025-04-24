// utils.go

package hw

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/cmd-stream/base-go"
	bcln "github.com/cmd-stream/base-go/client"
	bser "github.com/cmd-stream/base-go/server"
	ccln "github.com/cmd-stream/cmd-stream-go/client"
	cser "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/cmd-stream/handler-go"
)

// CmdReceiveDuration specifies the timeout after which the server will
// terminate the connection if no Command is received.
const CmdReceiveDuration = time.Second

// StartServer creates a listener, configures the server, and starts it
// in a separate goroutine.
func StartServer[T any](addr string, codec cser.Codec[T], receiver T,
	wg *sync.WaitGroup) (server *bser.Server, err error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}
	invoker := cser.NewInvoker(receiver)
	return StartServerWith(addr, codec, invoker, l.(*net.TCPListener), wg)
}

func StartServerWith[T any](addr string, codec cser.Codec[T],
	invoker handler.Invoker[T],
	l base.Listener,
	wg *sync.WaitGroup,
) (server *bser.Server, err error) {
	var callback bser.LostConnCallback = func(addr net.Addr, err error) {
		fmt.Printf("lost connection to %v, cause %v\n", addr, err)
	}
	server = cser.New(codec, invoker,
		// ServerInfo is optional and helps the client verify compatibility with the
		// server. It can identify supported commands or other server-specific
		// details. As a byte slice, it can store any arbitrary data.
		// cser.WithServerInfo(info)

		// Use Transport configuration to set the buffers size. If absent default
		// values from the bufio package are used.
		// cser.WithTransport(
		//   tcom.WithWriterBufSize(wsize),
		//   tcom.WithReaderBufSize(rsize)
		// )

		cser.WithHandler(
			// In a production environment, always set CmdReceiveTimeout. It allows
			// the server to close inactive client connections.
			handler.WithCmdReceiveDuration(CmdReceiveDuration),
			handler.WithAt(),
		),

		cser.WithBase(
			// WorkersCount determines the number of Workers, i.e., the number of
			// simultaneous connections to the server.
			bser.WithWorkersCount(8),

			// LostConnCallback is useful for debugging, it is called by the server
			// when the connection to the client is lost.
			bser.WithLostConnCallback(callback),
		),
	)

	wg.Add(1)
	go func(wg *sync.WaitGroup, listener base.Listener,
		server *bser.Server) {
		defer wg.Done()
		server.Serve(listener)
	}(wg, l, server)
	return
}

// utils.go

// CreateClient establishes a connection to the server, configures and
// creates a client.
func CreateClient[T any](addr string, codec ccln.Codec[T]) (
	client *bcln.Client[T], err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	var callback bcln.UnexpectedResultCallback = func(seq base.Seq,
		result base.Result) {
		fmt.Printf("unexpected result was received: seq %v, result %v\n",
			seq, result)
	}
	return ccln.New(codec, conn,
		ccln.WithBase(
			// UnexpectedResultCallback handles unexpected results from the
			// server. If you call Client.Forget(seq) for a Command, its results
			// will be handled by this function.
			bcln.WithUnexpectedResultCallback(callback),
		),
		// Use Transport configuration to set the buffers size. If absent
		// default values from the bufio package are used.
		// ccln.WithTransport(
		//   tcom.WithWriterBufSize(...),
		//   tcom.WithReaderBufSize(...),
		// ),
	)
}

// utils.go

// CloseServer closes the server and waits for it to stop.
func CloseServer(server *bser.Server, wg *sync.WaitGroup) (err error) {
	err = server.Close()
	if err != nil {
		return
	}

	// There is also the Server.Shutdown() method, which allows the server
	// to complete processing of already established connections without
	// accepting new ones.

	wg.Wait()
	return
}

// utils.go

// CloseClient closes the client and waits for it to stop.
func CloseClient[T any](client *bcln.Client[T]) (err error) {
	// Generally, the client will be closed if:
	// - Client.Close() is called.
	// - The server terminates the connection.
	//
	// In both cases, all uncompleted Commands will receive
	// AsyncResult.Error() == Client.Err(), where Client.Err() returns a
	// connection error.

	err = client.Close()
	if err != nil {
		return
	}
	// The client receives Results from the server in the background, so we
	// have to wait until it stops.
	select {
	case <-time.NewTimer(time.Second).C:
		return errors.New("timeout exceeded")
	case <-client.Done():
		return
	}
}

// CmdSendDuration defines the time during which the command should be
// sent.
const CmdSendDuration = time.Second

// Exchange sends a Command and checks whether the received greeting
// matches the expected value.
func Exchange[T any, R interface{ String() string }](cmd base.Cmd[T],
	timeout time.Duration,
	client *bcln.Client[T],
	wantGreeting R,
) (err error) {
	// Send the Command.
	var (
		seq      base.Seq
		results  = make(chan base.AsyncResult, 1)
		deadline = time.Now().Add(CmdSendDuration)
	)
	seq, err = client.SendWithDeadline(deadline, cmd, results)
	if err != nil {
		return
	}

	// Wait for the result with a timeout.
	var asyncResult base.AsyncResult
	select {
	case <-time.NewTimer(timeout).C:
		client.Forget(seq) // If you are no longer interested in the results of
		// this command, call Forget().
		return errors.New("timeout")
	case asyncResult = <-results:
	}

	// Check the Result.
	if asyncResult.Error != nil {
		return asyncResult.Error
	}
	greeting := asyncResult.Result.(R)
	if greeting.String() != wantGreeting.String() {
		return fmt.Errorf("unexpected greeting, want %v actual %v", wantGreeting,
			greeting)
	}
	return
}
