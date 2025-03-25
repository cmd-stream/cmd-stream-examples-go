package reconnect

import (
	"fmt"

	"github.com/cmd-stream/base-go"
	bcln "github.com/cmd-stream/base-go/client"
	ccln "github.com/cmd-stream/cmd-stream-go/client"
)

func CreateReconnectClient[T any](codec ccln.Codec[T],
	connFactory ccln.ConnFactory) (c *bcln.Client[T], err error) {

	var callback bcln.UnexpectedResultCallback = func(seq base.Seq, result base.Result) {
		fmt.Printf("unexpected result was received: seq %v, result %v\n", seq,
			result)
	}

	return ccln.NewReconnect[T](codec, connFactory,
		ccln.WithBase(
			bcln.WithUnexpectedResultCallback(callback),
		),
	)

}
