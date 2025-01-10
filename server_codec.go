package exmpls

import (
	"errors"
	"fmt"

	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
	dts "github.com/mus-format/mus-stream-dts-go"
)

// ServerCodec is a server Codec.
//
// A single ServerCodec will be used by all server Workers, so it must be
// thread-safe.
type ServerCodec struct{}

// Encode is used by the server to send results. If Encode fails with an error,
// the server closes the coresponding client connection.
func (c ServerCodec) Encode(result base.Result, w transport.Writer) (
	err error) {
	m, ok := result.(Marshaller)
	if !ok {
		return errors.New("result doesn't implement Marshaller interface")
	}
	return m.Marshal(w)
}

// Decode is used by the server to receive commands. If it fails with an error,
// the server closes the corresponding client connection.
func (c ServerCodec) Decode(r transport.Reader) (cmd base.Cmd[Greeter],
	err error) {
	// Unmarshal dtm.
	dtm, _, err := dts.UnmarshalDTM(r)
	if err != nil {
		return
	}
	// Depending on dtm, unmarshal a specific command.
	switch dtm {
	case SayHelloCmdDTM:
		cmd, _, err = SayHelloCmdDTS.UnmarshalData(r)
	case SayFancyHelloCmdDTM:
		cmd, _, err = SayFancyHelloCmdDTS.UnmarshalData(r)
	default:
		err = fmt.Errorf("unexpected cmd type %v", dtm)
	}
	return
}
