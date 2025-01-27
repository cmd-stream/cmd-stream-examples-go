package keepalive

import (
	"fmt"
	"net"
	"time"

	"github.com/cmd-stream/base-go"
	bcln "github.com/cmd-stream/base-go/client"
	ccln "github.com/cmd-stream/cmd-stream-go/client"
	cser "github.com/cmd-stream/cmd-stream-go/server"
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
	var (
		conf = ccln.Conf{
			// Use Transport configuration to set the buffers size. If absent default
			// values from the bufio package are used.
			// Transport transport_common.Conf{
			// 	 WriterBufSize: ... ,
			// 	 ReaderBufSize: ... ,
			// }
			Delegate: dcln.Conf{
				KeepaliveTime:  200 * time.Millisecond,
				KeepaliveIntvl: 200 * time.Millisecond,
			},
		}
		// callback handles unexpected results from the server. If you call
		// Client.Forget(seq) for a command, its results will be handled by
		// this function.
		callback bcln.UnexpectedResultCallback = func(seq base.Seq, result base.Result) {
			fmt.Printf("unexpected result was received: seq %v, result %v\n", seq,
				result)
		}
	)
	return ccln.New(conf, cser.DefaultServerInfo, codec, conn, callback)
}
