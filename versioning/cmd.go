package versioning

import (
	"context"
	"time"

	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"

	hw "github.com/cmd-stream/cmd-stream-examples-go/hello-world"
)

// NewOldSayHelloCmd creates a new OldSayHelloCmd.
func NewOldSayHelloCmd(str string) OldSayHelloCmd {
	return OldSayHelloCmd{str}
}

// OldSayHelloCmd is designed for the OldGreeter, implements base.Cmd and
// Marshaller interfaces.
type OldSayHelloCmd struct {
	str string
}

func (c OldSayHelloCmd) Exec(ctx context.Context, at time.Time, seq base.Seq,
	receiver OldGreeter,
	proxy base.Proxy,
) error {
	result := hw.NewResult(receiver.SayHello(c.str))
	return proxy.Send(seq, result)
}

func (c OldSayHelloCmd) Marshal(w transport.Writer) (err error) {
	_, err = OldSayHelloCmdDTS.Marshal(c, w)
	return
}

// Migrate is used by the server Codec.
func (c OldSayHelloCmd) Migrate() hw.SayHelloCmd {
	return hw.NewSayHelloCmd(c.str)
}
