package main

import (
	"context"
	"time"

	"github.com/cmd-stream/base-go"
	muss "github.com/mus-format/mus-stream-go"
)

// -----------------------------------------------------------------------------

type PrintCmdV1 struct {
	text string
}

func (c PrintCmdV1) Exec(ctx context.Context, at time.Time, seq base.Seq,
	receiver Printer,
	proxy base.Proxy,
) error {
	// From now, PrintCmdV1 should perform data migration to use Printer.
	receiver.Print("undefined", c.text)
	return proxy.Send(seq, OkResult(true))
}

func (c PrintCmdV1) MarshalMUS(w muss.Writer) (n int, err error) {
	return PrintCmdV1DTS.Marshal(c, w)
}

// -----------------------------------------------------------------------------

type PrintCmdV2 struct {
	from string
	text string
}

func (c PrintCmdV2) Exec(ctx context.Context, at time.Time, seq base.Seq,
	receiver Printer,
	proxy base.Proxy,
) error {
	receiver.Print(c.from, c.text)
	return proxy.Send(seq, OkResult(true))
}

func (c PrintCmdV2) MarshalMUS(w muss.Writer) (n int, err error) {
	return PrintCmdV2DTS.Marshal(c, w)
}
