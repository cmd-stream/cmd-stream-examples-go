package short

import (
	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
	"github.com/mus-format/mus-stream-go/ord"
)

// ClientCodec used to initialize the client.
type ClientCodec struct{}

func (c ClientCodec) Encode(cmd base.Cmd[struct{}], w transport.Writer) (
	err error) {
	_, err = ord.String.Marshal(string(cmd.(EchoCmd)), w)
	return
}

func (c ClientCodec) Decode(r transport.Reader) (result base.Result, err error) {
	str, _, err := ord.String.Unmarshal(r)
	if err != nil {
		return
	}
	result = Result(str)
	return
}

func (c ClientCodec) Size(cmd base.Cmd[struct{}]) (size int) {
	panic("not implemented")
}
