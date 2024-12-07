package main

import (
	"context"
	"time"

	"github.com/cmd-stream/base-go"
	muss "github.com/mus-format/mus-stream-go"
)

// Old version of the PrintCmd.
type PrintCmdV1 struct {
	Text string
}

func migrateV1(cmd PrintCmdV1) PrintCmd {
	return PrintCmd{
		From: "undefined",
		Text: cmd.Text,
	}
}

// Current version.
type PrintCmd struct {
	From string // New field compared to V1.
	Text string
}

func (c PrintCmd) Exec(ctx context.Context, at time.Time, seq base.Seq,
	receiver Printer,
	proxy base.Proxy,
) error {
	receiver.Print(c.From, c.Text)
	return proxy.Send(seq, OkResult(true))
}

func (c PrintCmd) MarshalMUS(w muss.Writer) (n int, err error) {
	return PrintCmdDTS.Marshal(c, w)
}
