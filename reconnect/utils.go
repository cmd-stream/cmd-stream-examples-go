package reconnect

import (
	"fmt"

	"github.com/cmd-stream/base-go"
	bcln "github.com/cmd-stream/base-go/client"
	ccln "github.com/cmd-stream/cmd-stream-go/client"
	cser "github.com/cmd-stream/cmd-stream-go/server"
)

func CreateReconnectClient[T any](codec ccln.Codec[T],
	connFactory ccln.ConnFactory) (c *bcln.Client[T], err error) {
	// unexpectedResultHandler processes unexpected results from the server.
	// If you call Client.Forget(seq) for a command, its results will be handled
	// by this function.
	unexpectedResultCallback := func(seq base.Seq, result base.Result) {
		fmt.Printf("unexpected result was received: seq %v, result %v\n", seq,
			result)
	}
	return ccln.NewReconnect[T](ccln.Conf{}, cser.DefaultServerInfo, codec,
		connFactory, unexpectedResultCallback)
	// return ccln.NewDefReconnect[T](codec,
	// 	connFactory,
	// 	unexpectedResultHandler)
}
