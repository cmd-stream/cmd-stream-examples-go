package server_streaming

import (
	hw "cmd-stream-examples-go/hello-world"
	"errors"
	"fmt"
	"time"

	"github.com/cmd-stream/base-go"
	bcln "github.com/cmd-stream/base-go/client"
)

func Exchange[T any](cmd base.Cmd[T], timeout time.Duration,
	client *bcln.Client[T],
	wantStrs []string,
) (err error) {
	var (
		seq      base.Seq
		results  = make(chan base.AsyncResult, 1)
		deadline = time.Now().Add(hw.CmdSendDuration)
	)
	seq, err = client.SendWithDeadline(deadline, cmd, results)
	if err != nil {
		return
	}

	for i := 0; i < len(wantStrs); i++ {
		// Waiting for the result with a timeout.
		select {
		case <-time.NewTimer(timeout).C:
			client.Forget(seq) // If you are no longer interested in the results of
			// this command, call Forget().
			return errors.New("timeout")
		case asyncResult := <-results:
			if asyncResult.Error != nil {
				return asyncResult.Error
			}
			greeting := asyncResult.Result.(Result).Str()
			if greeting != wantStrs[i] {
				return fmt.Errorf("unexpected greeting, want %v actual %v", wantStrs[i],
					greeting)
			}
		}
	}
	return
}
