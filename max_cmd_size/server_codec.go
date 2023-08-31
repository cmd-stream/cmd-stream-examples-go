package main

import (
	"errors"

	"github.com/cmd-stream/base-go"
	examples "github.com/cmd-stream/cmd-stream-examples-go"
	"github.com/cmd-stream/transport-go"
	com "github.com/mus-format/common-go"
	muss "github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/ord"
)

// ServerCodec checks the size of the incoming command in Decode, if it is too
// big, an error is returned (i.e., the connection to the client is closed).
type ServerCodec struct {
	examples.ServerCodec
}

func (c ServerCodec) Decode(r transport.Reader) (cmd base.Cmd[struct{}],
	err error) {
	cmd, _, err = UnmarshalValidEchoCmdMUS(r)
	return
}

func UnmarshalValidEchoCmdMUS(r muss.Reader) (c examples.EchoCmd, n int,
	err error) {
	var maxLength com.ValidatorFn[int] = func(length int) (err error) {
		if length > MaxCmdLength {
			err = errors.New("command is too big")
		}
		return
	}
	str, n, err := ord.UnmarshalValidString(maxLength, false, r)
	c = examples.EchoCmd(str)
	return
}
