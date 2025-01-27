package server_streaming

import (
	"errors"
	"fmt"

	hw "cmd-stream-examples-go/hello-world"

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
	m, ok := result.(hw.Marshaller)
	if !ok {
		return errors.New("result doesn't implement Marshaller interface")
	}
	return m.Marshal(w)
}

// Decode is used by the server to receive commands from the client. If it
// fails with an error, the server closes the connection.
func (c ServerCodec) Decode(r transport.Reader) (cmd base.Cmd[hw.Greeter],
	err error) {
	// Unmarshals dtm.
	dtm, _, err := dts.UnmarshalDTM(r)
	if err != nil {
		return
	}
	// Depending on dtm, unmarshal a specific command.
	switch dtm {
	case SayFancyHelloMultiCmdDTM:
		cmd, _, err = SayFancyHelloMultiCmdDTS.UnmarshalData(r)
	default:
		err = fmt.Errorf("unexpected cmd type %v", dtm)
	}
	return
}
