package main

import (
	"errors"

	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
	dts "github.com/mus-format/mus-stream-dts-go"
)

type ServerCodec struct{}

// Encode is used by the server to send results to the client. If Encode fails
// with an error, the server closes the connection.
func (c ServerCodec) Encode(result base.Result, w transport.Writer) (
	err error) {
	m, ok := result.(MarshallerMUS)
	if !ok {
		return errors.New("cmd doesn't implement MarshallerMUS interface")
	}
	_, err = m.MarshalMUS(w)
	return
}

// Decode is used by the server to receive commands from the client. If Decode
// fails with an error, the server closes the connection.
func (c ServerCodec) Decode(r transport.Reader) (cmd base.Cmd[Printer],
	err error) {
	// Unmarshals dtm.
	dtm, _, err := dts.UnmarshalDTM(r)
	if err != nil {
		return
	}
	// Depending on dtm, unmarshals a specific command.
	switch dtm {
	case PrintCmdV1DTM:
		cmd, _, err = PrintCmdV1DTS.UnmarshalData(r)
	case PrintCmdV2DTM:
		cmd, _, err = PrintCmdV2DTS.UnmarshalData(r)
	default:
		err = ErrUnsupportedCmdType
	}
	return
}
