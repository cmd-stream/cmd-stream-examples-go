package group

import (
	"errors"
	"fmt"
	"time"

	hw "github.com/cmd-stream/cmd-stream-examples-go/hello-world"

	"github.com/cmd-stream/base-go"
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
