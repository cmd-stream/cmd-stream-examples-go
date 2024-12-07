package old_client

import (
	"context"
	"time"

	"github.com/cmd-stream/base-go"
	muss "github.com/mus-format/mus-stream-go"
)

// -----------------------------------------------------------------------------

type PrintCmd struct {
	Text string
}

func (c PrintCmd) Exec(ctx context.Context, at time.Time, seq base.Seq,
	receiver Printer,
	proxy base.Proxy,
) error {
	receiver.Print(c.Text)
	return proxy.Send(seq, OkResult(true))
}

func (c PrintCmd) MarshalMUS(w muss.Writer) (n int, err error) {
	return PrintCmdDTS.Marshal(c, w)
}
