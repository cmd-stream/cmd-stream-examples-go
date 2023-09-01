package examples

import (
	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
)

type ClientCodec struct{}

// Encode is used by the client to send commands to the server. If Encode fails
// with an error, the Client.Send method will return it.
func (c ClientCodec) Encode(cmd base.Cmd[struct{}], w transport.Writer) (
	err error) {
	_, err = MarshalEchoCmdMUS(cmd.(EchoCmd), w)
	return
}

// Decode is used by the client to receive resulsts from the server. If Decode
// fails with an error, the client will be closed.
func (c ClientCodec) Decode(r transport.Reader) (result base.Result, err error) {
	cmd, _, err := UnmarshalEchoCmdMUS(r)
	if err != nil {
		return
	}
	result = OneEchoResult(cmd)
	return
}

// Size returns the size of the command in bytes. If the server imposes any
// restrictions on the size of the command, the client will use this method to
// check the size of the command before sending it.
func (c ClientCodec) Size(cmd base.Cmd[struct{}]) (size int) {
	return SizeEchoCmdMUS(cmd.(EchoCmd))
}
