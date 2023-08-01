package main

import (
	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
)

type ClientCodec struct{}

// Used by the client to send commands to the server.
//
// We have two kinds of commands, that's why we encode a command type, than a
// command itself.
func (c ClientCodec) Encode(cmd base.Cmd[Printer], w transport.Writer) (
	err error) {
	switch c := cmd.(type) {
	case PrintCmdV1:
		_, err = MarshalCmdTypeMUS(PrintCmdV1Type, w)
		if err != nil {
			return
		}
		_, err = MarshalPrintCmdV1MUS(c, w)
	case PrintCmdV2:
		_, err = MarshalCmdTypeMUS(PrintCmdV2Type, w)
		if err != nil {
			return
		}
		_, err = MarshalPrintCmdV2MUS(c, w)
	default:
		panic("unexpected cmd type")
	}
	return
}

// Used by the client to receive resulsts from the server.
//
// We have only one type of result.
func (c ClientCodec) Decode(r transport.Reader) (result base.Result, err error) {
	result, _, err = UnmarshalOkResultMUS(r)
	return
}

// Size returns the size of the command in bytes. If the server imposes any
// restrictions on the command size, the client will use this method to
// check it before sending.
func (c ClientCodec) Size(cmd base.Cmd[Printer]) (size int) {
	switch c := cmd.(type) {
	case PrintCmdV1:
		size = SizePrintCmdV1MUS(c)
	case PrintCmdV2:
		size = SizePrintCmdV2MUS(c)
	default:
		panic("unexpected cmd type")
	}
	return
}
