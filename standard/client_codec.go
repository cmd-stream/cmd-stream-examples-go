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
func (c ClientCodec) Encode(cmd base.Cmd[Calculator], w transport.Writer) (
	err error) {
	switch c := cmd.(type) {
	case Eq1Cmd:
		_, err = MarshalCmdTypeMUS(Eq1CmdType, w)
		if err != nil {
			return
		}
		_, err = MarshalEq1CmdMUS(c, w)
	case Eq2Cmd:
		_, err = MarshalCmdTypeMUS(Eq2CmdType, w)
		if err != nil {
			return
		}
		_, err = MarshalEq2CmdMUS(c, w)
	default:
		panic("unexpected cmd type")
	}
	return
}

// Used by the client to receive resulsts from the server.
//
// We have only one type of result.
func (c ClientCodec) Decode(r transport.Reader) (result base.Result, err error) {
	result, _, err = UnmarshalResultMUS(r)
	return
}

// Size returns the size of the command in bytes. If the server imposes any
// restrictions on the command size, the client will use this method to
// check it before sending.
func (c ClientCodec) Size(cmd base.Cmd[Calculator]) (size int) {
	switch c := cmd.(type) {
	case Eq1Cmd:
		size = SizeEq1CmdMUS(c)
	case Eq2Cmd:
		size = SizeEq2CmdMUS(c)
	default:
		panic("unexpected cmd type")
	}
	return
}
