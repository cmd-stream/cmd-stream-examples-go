package main

import (
	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
)

type ClientCodec struct{}

func (c ClientCodec) Encode(cmd base.Cmd[any], w transport.Writer) (
	err error) {
	_, err = MarshalMultiEchoCmdMUS(cmd.(MultiEchoCmd), w)
	return
}

func (c ClientCodec) Decode(r transport.Reader) (result base.Result, err error) {
	result, _, err = UnmarshalMultiEchoResultMUS(r)
	return
}

// Size returns the size of the command in bytes. If the server imposes any
// restrictions on the size of the command, the client will use this method to
// check the size of the command before sending it.
func (c ClientCodec) Size(cmd base.Cmd[any]) (size int) {
	switch c := cmd.(type) {
	case MultiEchoCmd:
		size = SizeMultiEchoCmdMUS(c)
	default:
		panic("unexpected cmd type")
	}
	return
}
