package main

import (
	"errors"

	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
	dts "github.com/mus-format/mus-stream-dts-go"
)

type ClientCodec struct{}

// Encode is used by the client to send commands to the server. If Encode fails
// with an error, the Client.Send method will return it.
func (c ClientCodec) Encode(cmd base.Cmd[Printer], w transport.Writer) (
	err error) {
	m, ok := cmd.(MarshallerMUS)
	if !ok {
		return errors.New("cmd doesn't implement MarshallerMUS interface")
	}
	_, err = m.MarshalMUS(w)
	return
}

// Decode is used by the client to receive resulsts from the server. If Decode
// fails with an error, the client will be closed.
func (c ClientCodec) Decode(r transport.Reader) (result base.Result,
	err error) {
	// Unmarshals dtm.
	dtm, _, err := dts.UnmarshalDTM(r)
	if err != nil {
		return
	}
	// Depending on dtm, unmarshals a specific result.
	switch dtm {
	case OkResultDTM:
		result, _, err = OkResultDTS.UnmarshalData(r)
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
