package short

import (
	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
	"github.com/mus-format/mus-stream-go/ord"
)

// ServerCodec used to initialize the server.
type ServerCodec struct{}

func (c ServerCodec) Encode(result base.Result, w transport.Writer) (err error) {
	_, err = ord.MarshalString(string(result.(Result)), nil, w)
	return
}

func (c ServerCodec) Decode(r transport.Reader) (cmd base.Cmd[struct{}],
	err error) {
	str, _, err := ord.UnmarshalString(nil, r)
	if err != nil {
		return
	}
	cmd = EchoCmd(str)
	return
}
