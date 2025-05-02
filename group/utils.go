package group

import (
	"errors"
	"fmt"
	"net"
	"time"

	hw "github.com/cmd-stream/cmd-stream-examples-go/hello-world"
	dcln "github.com/cmd-stream/delegate-go/client"

	"github.com/cmd-stream/base-go"
	bcln "github.com/cmd-stream/base-go/client"
	ccln "github.com/cmd-stream/cmd-stream-go/client"
)

func Exchange[T any, R interface{ String() string }](cmd base.Cmd[T],
	timeout time.Duration,
	grp ccln.Group[T],
	wantGreeting R,
) (err error) {
	// Send the Command.
	var (
		seq      base.Seq
		clientID ccln.ClientID
		results  = make(chan base.AsyncResult, 1)
		deadline = time.Now().Add(hw.CmdSendDuration)
	)
	seq, clientID, err = grp.SendWithDeadline(deadline, cmd, results)
	if err != nil {
		return
	}

	// Wait for the result with a timeout.
	var asyncResult base.AsyncResult
	select {
	case <-time.NewTimer(timeout).C:
		grp.Forget(seq, clientID) // If you are no longer interested in the results of
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

// CreateClientGroup establishes several connections to the server, configures
// and creates a client group.
func CreateClientGroup(addr string, clientCount int,
	codec ccln.Codec[hw.Greeter]) ccln.Group[hw.Greeter] {
	var (
		// codec   = cdc.NewClientCodec(CmdMUS, ResultMUS)
		factory = ccln.ConnFactoryFn(func() (net.Conn, error) {
			return net.Dial("tcp", addr)
		})
		callback bcln.UnexpectedResultCallback = func(seq base.Seq,
			result base.Result) {
			fmt.Printf("unexpected result was received: seq %v, result %v\n",
				seq, result)
		}
		clients = ccln.MustMakeClients(clientCount, codec, factory,
			ccln.WithBase(
				// UnexpectedResultCallback handles unexpected results from the
				// server. If you call Client.Forget(seq) for a Command, its results
				// will be handled by this function.
				bcln.WithUnexpectedResultCallback(callback),
			),
			ccln.WithKeepalive(
				dcln.WithKeepaliveIntvl(time.Second),
				dcln.WithKeepaliveTime(time.Second),
			),
			// Use Transport configuration to set the buffers size. If absent
			// default values from the bufio package are used.
			// ccln.WithTransport(
			//   tcom.WithWriterBufSize(...),
			//   tcom.WithReaderBufSize(...),
			// ),
		)
		strategy = ccln.NewRoundRobinStrategy(clients)
	)
	return ccln.NewGroup(strategy)
}

// CloseGroup closes the client group and waits for it to stop.
func CloseGroup[T any](grp ccln.Group[T]) (err error) {
	err = grp.Close()
	if err != nil {
		return
	}
	// The clients receive Results from the server in the background, so we
	// have to wait until the group stops.
	select {
	case <-time.NewTimer(time.Second).C:
		return errors.New("timeout exceeded")
	case <-grp.Done():
		return
	}
}
