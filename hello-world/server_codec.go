package hw

import (
	"errors"
	"fmt"

	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
	dts "github.com/mus-format/dts-stream-go"
)

// A single ServerCodec instance is shared by all server Workers (each Worker
// handles one client connection at a time), so it must be thread-safe.
type ServerCodec struct{}

func (c ServerCodec) Encode(result base.Result, w transport.Writer) (
	err error) {
	m, ok := result.(Marshaller) // Marshaller interface is used again.
	if !ok {
		return errors.New("result doesn't implement Marshaller interface")
	}
	return m.Marshal(w)
}

func (c ServerCodec) Decode(r transport.Reader) (cmd base.Cmd[Greeter],
	err error) {
	// Using dts-stream-go library unmarshal dtm.
	dtm, _, err := dts.DTMSer.Unmarshal(r)
	if err != nil {
		return
	}
	// Unmarshal a specific Command based on dtm. This ensures the server only
	// executes expected Commands, preventing unintended behavior.
	switch dtm {
	case SayHelloCmdDTM:
		cmd, _, err = SayHelloCmdDTS.UnmarshalData(r)
	case SayFancyHelloCmdDTM:
		cmd, _, err = SayFancyHelloCmdDTS.UnmarshalData(r)
	default:
		err = fmt.Errorf("unexpected Command type %v", dtm)
	}
	return
}
