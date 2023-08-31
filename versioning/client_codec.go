package main

import (
	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
	dts "github.com/mus-format/mus-stream-dts-go"
)

type ClientCodec struct{}

// Used by the client to send commands to the server.
func (c ClientCodec) Encode(cmd base.Cmd[Printer], w transport.Writer) (
	err error) {
	// With help of type assertions, marshals a specific command.
	switch c := cmd.(type) {
	case PrintCmdV1:
		_, err = PrintCmdV1DTS.MarshalMUS(c, w)
	case PrintCmdV2:
		_, err = PrintCmdV2DTS.MarshalMUS(c, w)
	default:
		err = ErrUnsupportedCmdType
	}
	return
}

// Used by the client to receive resulsts from the server.
func (c ClientCodec) Decode(r transport.Reader) (result base.Result,
	err error) {
	// Unmarshals dtm.
	dtm, _, err := dts.UnmarshalDTMUS(r)
	if err != nil {
		return
	}
	// Depending on dtm, unmarshals a specific result.
	switch dtm {
	case OkResultDTM:
		result, _, err = OkResultDTS.UnmarshalDataMUS(r)
	default:
		err = ErrUnsupportedResultType
	}
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
