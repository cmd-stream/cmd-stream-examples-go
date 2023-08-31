package main

import (
	"errors"

	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
	dts "github.com/mus-format/mus-stream-dts-go"
)

type ClientCodec struct{}

// Used by the client to send commands to the server.
func (c ClientCodec) Encode(cmd base.Cmd[Calculator], w transport.Writer) (
	err error) {
	// With help of type assertions, marshals a specific command.
	switch c := cmd.(type) {
	case Eq1Cmd:
		_, err = Eq1DTS.MarshalMUS(c, w)
	case Eq2Cmd:
		_, err = Eq2DTS.MarshalMUS(c, w)
	default:
		panic("unexpected cmd type")
	}
	return
}

// Used by the client to receive resulsts from the server.
func (c ClientCodec) Decode(r transport.Reader) (result base.Result, err error) {
	// Unmarshals dtm.
	dtm, _, err := dts.UnmarshalDTMUS(r)
	if err != nil {
		return
	}
	// Depending on dtm, unmarshals a specific result.
	switch dtm {
	case ResultDTM:
		result, _, err = ResultDTS.UnmarshalDataMUS(r)
	default:
		err = errors.New("unexpected result type")
	}
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
