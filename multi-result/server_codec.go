package main

import (
	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
)

type ServerCodec struct{}

func (c ServerCodec) Encode(result base.Result, w transport.Writer) (
	err error) {
	_, err = MarshalMultiEchoResultMUS(result.(MultiEchoResult), w)
	return
}

func (c ServerCodec) Decode(r transport.Reader) (cmd base.Cmd[any], err error) {
	cmd, _, err = UnmarshalMultiEchoCmdMUS(r)
	return
}
