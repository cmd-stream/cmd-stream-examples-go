package main

import (
	"errors"
	"fmt"

	"github.com/cmd-stream/base-go"
	exmpls "github.com/cmd-stream/cmd-stream-examples-go"
	"github.com/cmd-stream/transport-go"
	dts "github.com/mus-format/mus-stream-dts-go"
)

// One ServerCodec will be used by all server Workers, so it must be thread-safe.
// In the Decode method, it performs the migration to the current command
// version.
type ServerCodec struct{}

// Encode is used by the server to send results to the client. If Encode fails
// with an error, the server closes the connection.
func (c ServerCodec) Encode(result base.Result, w transport.Writer) (
	err error) {
	m, ok := result.(exmpls.Marshaller)
	if !ok {
		return errors.New("result doesn't implement Marshaller interface")
	}
	return m.Marshal(w)
}

// Decode is used by the server to receive commands from the client. If it
// fails with an error, the server closes the connection.
func (c ServerCodec) Decode(r transport.Reader) (cmd base.Cmd[exmpls.Greeter],
	err error) {
	// Unmarshals dtm.
	dtm, _, err := dts.UnmarshalDTM(r)
	if err != nil {
		return
	}
	// Depending on dtm, unmarshal a specific command.
	switch dtm {
	case exmpls.OldSayHelloCmdDTM:
		var oldCmd exmpls.OldSayHelloCmd
		oldCmd, _, err = exmpls.OldSayHelloCmdDTS.UnmarshalData(r)
		if err != nil {
			return
		}
		cmd = oldCmd.Migrate()
	case exmpls.SayHelloCmdDTM:
		cmd, _, err = exmpls.SayHelloCmdDTS.UnmarshalData(r)
	default:
		err = fmt.Errorf("unexpected cmd type %v", dtm)
	}
	return
}
