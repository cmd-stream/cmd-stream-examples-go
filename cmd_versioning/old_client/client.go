package old_client

import (
	"net"

	base_client "github.com/cmd-stream/base-go/client"
	cs_client "github.com/cmd-stream/cmd-stream-go/client"
)

func CreateClient(addr string) (
	c *base_client.Client[Printer], err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	// The last nil parameter corresponds to the UnexpectedResultHandler. In this
	// case, unexpected results (if any) received from the server will be simply
	// ignored.
	return cs_client.NewDef[Printer](ClientCodec{}, conn, nil)
}
