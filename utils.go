package exmpls

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/delegate-go"
	"github.com/cmd-stream/handler-go"
	assert "github.com/ymz-ncnk/assert/panic"

	base_client "github.com/cmd-stream/base-go/client"
	base_server "github.com/cmd-stream/base-go/server"
	cs_client "github.com/cmd-stream/cmd-stream-go/client"
	cs_server "github.com/cmd-stream/cmd-stream-go/server"
)

// SendCmdDeadline is a deadline used to send commands.
const SendCmdDeadline = time.Second

// ReceiveTimeout determines how long the server will wait for a command from
// the client before closing the connection.
const ReceiveTimeout = 2 * time.Second

// SendCmd sends a command using the specified client.
//
// It compares the received results with the expected ones.
func SendCmd[T any, R comparable](cmd base.Cmd[T],
	timeout time.Duration,
	wantErr error,
	wantResults []R,
	compareResultsFn func(result, wantResult R),
	client *base_client.Client[T],
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	var (
		seq      base.Seq
		results  = make(chan base.AsyncResult, 1)
		deadline = time.Now().Add(SendCmdDeadline)
	)
	seq, err := client.SendWithDeadline(deadline, cmd, results)
	assert.EqualError(err, nil)

	receiveResults[T, R](seq, results, timeout, wantErr, wantResults,
		compareResultsFn, client, 0)
}

func receiveResults[T any, R comparable](seq base.Seq,
	results chan base.AsyncResult,
	timeout time.Duration,
	wantErr error,
	wantResults []R,
	compareResultsFn func(result, wantResult R),
	client *base_client.Client[T],
	i int,
) {
	select {
	case <-time.NewTimer(timeout).C:
		client.Forget(seq) // If you are no longer interested in the results of
		// this command, call Forget().
	case asyncResult := <-results:
		var result R
		if asyncResult.Error != nil {
			assert.EqualError(asyncResult.Error, wantErr)
			return
		}
		result = asyncResult.Result.(R)
		compareResultsFn(result, wantResults[i])

		if !asyncResult.Result.LastOne() {
			i += 1
			receiveResults[T, R](seq, results, timeout, wantErr, wantResults,
				compareResultsFn, client, i)
		}
	}
}

// StartServer creates and launches a server with the specified address, codec,
// receiver, and the default Invoker.
func StartServer[T any](addr string, codec cs_server.Codec[T], receiver T,
	wg *sync.WaitGroup) (s *base_server.Server, err error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}
	return StartServerWith(addr, codec, receiver, l.(*net.TCPListener),
		delegate.ServerSettings{}, wg)
}

func StartServerWith[T any](addr string, codec cs_server.Codec[T], receiver T,
	l base.Listener,
	settings delegate.ServerSettings,
	wg *sync.WaitGroup,
) (s *base_server.Server, err error) {
	conf := cs_server.Conf{
		// Transport: transport_common.Conf{ // Use it to set buffers size.
		// 	WriterBufSize: ...,
		//  ReaderBufSize: ...,
		// },
		Handler: handler.Conf{
			ReceiveTimeout: ReceiveTimeout, // In a production environment, always
			// set ReceiveTimeout - it allows the server to close inactive client
			// connections.
			At: true, // Commands will receive not nill 'at' parameter.
		},
		Base: base_server.Conf{
			WorkersCount: 8, // Determines the number of simultaneous connections to
			// the server.
			LostConnCallback: func(addr net.Addr, err error) { // LostConnCallback is
				// useful for debugging, it is called by the server when the connection
				// to the client is lost.
				fmt.Printf("lost connection to %v, cause %v\n", addr, err)
			},
		},
	}
	s = cs_server.New(cs_server.DefServerInfo,
		settings,
		conf,
		codec,
		receiver,
		nil) // If nil the default invoker will be used.

	wg.Add(1)
	go func(wg *sync.WaitGroup, listener base.Listener,
		server *base_server.Server) {
		defer wg.Done()
		err := server.Serve(listener)
		assert.EqualError(err, base_server.ErrClosed)
	}(wg, l, s)
	return
}

// CreateClient creates a client with the specified server addr and codec.
func CreateClient[T any](addr string, codec cs_client.Codec[T]) (
	c *base_client.Client[T], err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	conf := cs_client.Conf{
		//  Transport transport_common.Conf{ // Use it to set buffers size.
		//  	 WriterBufSize: ... ,
		//  	 ReaderBufSize: ... ,
		//  }
	}
	return CreateClientWith(conf, codec, conn)
}

func CreateClientWith[T any](conf cs_client.Conf,
	codec cs_client.Codec[T],
	conn net.Conn,
) (c *base_client.Client[T], err error) {
	// unexpectedResultHandler handles unexpected results from the server. If you
	// call Client.Forget(seq) for a command, its results will be handled by
	// this function.
	unexpectedResultHandler := func(seq base.Seq, result base.Result) {
		fmt.Printf("unexpected result was received: seq %v, result %v\n", seq,
			result)
	}
	return cs_client.New(cs_server.DefServerInfo, conf, codec, conn,
		unexpectedResultHandler)
}

// CloseServer closes the server.
func CloseServer(s *base_server.Server, wg *sync.WaitGroup) (err error) {
	err = s.Close()
	if err != nil {
		return
	}
	wg.Wait()
	return
}

// CloseClient closes the client.
func CloseClient[T any](c *base_client.Client[T]) (err error) {
	err = c.Close()
	if err != nil {
		return
	}
	// The client receives results from the server in the background, so we have
	// to wait for it to stop.
	select {
	case <-time.NewTimer(time.Second).C:
		return errors.New("can't close the client")
	case <-c.Done():
		// To get a client's connection error, use Client.Err(). If the client has
		// already been closed, all uncompleted commands will receive result with
		// AsyncResult.Error() == Client.Err().
		return
	}
}

func CompareResults[R comparable](result, wantResult R) {
	assert.Equal[R](result, wantResult)
}
