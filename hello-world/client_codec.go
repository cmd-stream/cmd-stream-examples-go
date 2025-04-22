package hw

import (
	"errors"
	"fmt"

	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
	dts "github.com/mus-format/dts-stream-go"
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
	// Unmarshal DTM using the dts-stream-go library.
	dtm, _, err := dts.DTMSer.Unmarshal(r)
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
