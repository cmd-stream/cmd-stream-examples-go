package examples

import (
	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
)

type ServerCodec struct{}

// If the Encode method fails with an error, the server closes the connection.
func (c ServerCodec) Encode(result base.Result, w transport.Writer) (err error) {
	_, err = MarshalEchoCmdMUS(result.(EchoCmd), w)
	return
}

// If the Decode method fails with an error, the server closes the connection.
func (c ServerCodec) Decode(r transport.Reader) (cmd base.Cmd[struct{}], err error) {
	cmd, _, err = UnmarshalEchoCmdMUS(r)
	return
}
