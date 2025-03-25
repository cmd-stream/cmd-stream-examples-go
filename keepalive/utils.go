package keepalive

import (
	"fmt"
	"net"
	"time"

	"github.com/cmd-stream/base-go"
	bcln "github.com/cmd-stream/base-go/client"
	ccln "github.com/cmd-stream/cmd-stream-go/client"
	dcln "github.com/cmd-stream/delegate-go/client"
)

// CreateKeepaliveClient creates a client that will keep a connection alive
// when there are no commands to send.
func CreateKeepaliveClient[T any](addr string, codec ccln.Codec[T]) (
	client *bcln.Client[T], err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}

	var callback bcln.UnexpectedResultCallback = func(seq base.Seq, result base.Result) {
		fmt.Printf("unexpected result was received: seq %v, result %v\n", seq,
			result)
	}
	return ccln.New(codec, conn,
		ccln.WithBase(
			bcln.WithUnexpectedResultCallback(callback),
		),
		ccln.WithKeepalive(
			dcln.WithKeepaliveTime(200*time.Millisecond),
			dcln.WithKeepaliveIntvl(200*time.Millisecond),
		),
	)

}
