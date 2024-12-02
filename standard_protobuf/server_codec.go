package main

import (
	"errors"

	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
	dts "github.com/mus-format/mus-stream-dts-go"
)

// One ServerCodec will be used by all server Workers, so it must be thread-safe.
type ServerCodec struct{}

// Encode is used by the server to send results to the client. If Encode fails
// with an error, the server closes the connection.
func (c ServerCodec) Encode(result base.Result, w transport.Writer) (
	err error) {
	m, ok := result.(MarshallerProtobuf)
	if !ok {
		return errors.New("result doesn't implement MarshallerProtobuf interface")
	}
	_, err = m.MarshalProtobuf(w)
	return
}

// Decode is used by the server to receive commands from the client. If Decode
// fails with an error, the server closes the connection.
func (c ServerCodec) Decode(r transport.Reader) (cmd base.Cmd[Calculator],
	err error) {
	// Unmarshals dtm.
	dtm, _, err := dts.UnmarshalDTM(r)
	if err != nil {
		return
	}
	// Depending on dtm, unmarshals a specific command.
	switch dtm {
	case Eq1DTM:
		cmd, _, err = Eq1DTS.UnmarshalData((r))
	case Eq2DTM:
		cmd, _, err = Eq2DTS.UnmarshalData((r))
	default:
		err = errors.New("unexpected cmd type")
	}
	return
}
