package main

import (
	"context"
	"time"

	"github.com/cmd-stream/base-go"
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
