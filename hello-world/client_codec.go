package hw

import (
	"errors"
	"fmt"

	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
	dts "github.com/mus-format/mus-stream-dts-go"
)

type ClientCodec struct{}

func (c ClientCodec) Encode(cmd base.Cmd[Greeter], w transport.Writer) (
	err error) {
	m, ok := cmd.(Marshaller) // Use defined previously Marshaller interface.
	if !ok {
		return errors.New("cmd doesn't implement the Marshaller interface")
	}
	return m.Marshal(w)
}

func (c ClientCodec) Decode(r transport.Reader) (result base.Result, err error) {
	// Unmarshal DTM using the mus-stream-dts-go library.
	dtm, _, err := dts.UnmarshalDTM(r)
	if err != nil {
		return
	}
	// Depending on dtm, unmarshal a specific Result.
	switch dtm {
	case ResultDTM:
		result, _, err = ResultDTS.UnmarshalData(r)
	default:
		err = fmt.Errorf("unexpected Result type %v", dtm)
	}
	return
}

func (c ClientCodec) Size(cmd base.Cmd[Greeter]) (size int) {
	// Implementation is unnecessary as ServerSettings.MaxCmdSize == 0.
	panic("not implemented")
}
