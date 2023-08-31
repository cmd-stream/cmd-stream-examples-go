package main

import (
	"errors"

	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
	dts "github.com/mus-format/mus-stream-dts-go"
)

type ServerCodec struct{}

// Used by the server to send results to the client.
func (c ServerCodec) Encode(result base.Result, w transport.Writer) (
	err error) {
	// With help of type assertions, marshals a specific result.
	switch rt := result.(type) {
	case Result:
		_, err = ResultDTS.MarshalMUS(rt, w)
	default:
		err = errors.New("unexpected result type")
	}
	return
}

// Used by the server to receive commands from the client.
func (c ServerCodec) Decode(r transport.Reader) (cmd base.Cmd[Calculator],
	err error) {
	// Unmarshals dtm.
	dtm, _, err := dts.UnmarshalDTMUS(r)
	if err != nil {
		return
	}
	// Depending on dtm, unmarshals a specific command.
	switch dtm {
	case Eq1DTM:
		cmd, _, err = Eq1DTS.UnmarshalDataMUS((r))
	case Eq2DTM:
		cmd, _, err = Eq2DTS.UnmarshalDataMUS((r))
	default:
		err = errors.New("unexpected cmd type")
	}
	return
}
