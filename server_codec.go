package examples

import (
	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
)

type ServerCodec struct{}

// Encode is used by the server to send results to the client. If Encode fails
// with an error, the server closes the connection.
func (c ServerCodec) Encode(result base.Result, w transport.Writer) (err error) {
	_, err = MarshalEchoCmdMUS(result.(EchoCmd), w)
	return
}

// Decode is used by the server to receive commands from the client. If Decode
// fails with an error, the server closes the connection.
func (c ServerCodec) Decode(r transport.Reader) (cmd base.Cmd[struct{}], err error) {
	cmd, _, err = UnmarshalEchoCmdMUS(r)
	return
}
