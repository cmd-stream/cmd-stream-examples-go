package versioning

import (
	"errors"

	hw "cmd-stream-examples-go/hello-world"

	"github.com/cmd-stream/base-go"

	"github.com/cmd-stream/transport-go"
	dts "github.com/mus-format/mus-stream-dts-go"
)

type ClientCodec struct{}

// Encode is used by the client to send commands to the server. If Encode fails
// with an error, the Client.Send() method will return it.
func (c ClientCodec) Encode(cmd base.Cmd[OldGreeter], w transport.Writer) (
	err error) {
	m, ok := cmd.(hw.Marshaller)
	if !ok {
		return errors.New("cmd doesn't implement Marshaller interface")
	}
	return m.Marshal(w)
}

// Decode is used by the client to receive resulsts from the server. If Decode
// fails with an error, the client will be closed.
func (c ClientCodec) Decode(r transport.Reader) (result base.Result, err error) {
	// Unmarshal dtm.
	dtm, _, err := dts.UnmarshalDTM(r)
	if err != nil {
		return
	}
	// Depending on dtm, unmarshal a specific result.
	switch dtm {
	case hw.ResultDTM:
		result, _, err = hw.ResultDTS.UnmarshalData(r)
	default:
		err = errors.New("unexpected result type")
	}
	return
}

// Size returns the size of the command in bytes. If the server imposes any
// restrictions on the command size, the client will use this method to
// check it before sending.
func (c ClientCodec) Size(cmd base.Cmd[OldGreeter]) (size int) {
	panic("not implemented")
}
