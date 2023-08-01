package main

import (
	"errors"

	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
)

type ServerCodec struct{}

// Used by the server to send results to the client.
//
// We have only one type of result.
func (c ServerCodec) Encode(result base.Result, w transport.Writer) (
	err error) {
	_, err = MarshalOkResultMUS(result.(OkResult), w) // We have only one type of
	// result.
	return
}

// Used by the server to receive commands from the client.
//
// We have two kinds of commands, that's why we decode a command type, than a
// command itself.
func (c ServerCodec) Decode(r transport.Reader) (cmd base.Cmd[Printer],
	err error) {
	tp, _, err := UnmarshalCmdTypeMUS(r)
	if err != nil {
		return
	}
	switch tp {
	case PrintCmdV1Type:
		cmd, _, err = UnmarshalPrintCmdV1MUS(r)
	case PrintCmdV2Type:
		cmd, _, err = UnmarshalPrintCmdV2MUS(r)
	default:
		err = errors.New("unexpected cmd type")
	}
	return
}
