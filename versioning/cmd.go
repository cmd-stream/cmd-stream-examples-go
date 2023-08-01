package main

import (
	"context"
	"time"

	"github.com/cmd-stream/base-go"
	muss "github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/ord"
	"github.com/mus-format/mus-stream-go/raw"
)

const (
	PrintCmdV1Type = iota + 1
	PrintCmdV2Type
)

type CmdType byte

func MarshalCmdTypeMUS(ct CmdType, w muss.Writer) (n int, err error) {
	return raw.MarshalByte(byte(ct), w)
}

func UnmarshalCmdTypeMUS(r muss.Reader) (ct CmdType, n int, err error) {
	b, n, err := raw.UnmarshalByte(r)
	ct = CmdType(b)
	return
}

func SizeCmdTypeMUS(ct CmdType) (size int) {
	return raw.SizeByte(byte(ct))
}

// -----------------------------------------------------------------------------
type PrintCmdV1 struct {
	text string
}

func (c PrintCmdV1) Exec(ctx context.Context, at time.Time, seq base.Seq,
	receiver Printer,
	proxy base.Proxy,
) error {
	// PrintCmdV1 does not have a "from" field, so it uses the "undefined" value.
	receiver.Print("undefined", c.text)
	return proxy.Send(seq, OkResult(true))
}

func MarshalPrintCmdV1MUS(c PrintCmdV1, w muss.Writer) (n int, err error) {
	return ord.MarshalString(c.text, w)
}

func UnmarshalPrintCmdV1MUS(r muss.Reader) (c PrintCmdV1, n int, err error) {
	c.text, n, err = ord.UnmarshalString(r)
	return
}

func SizePrintCmdV1MUS(c PrintCmdV1) (size int) {
	return ord.SizeString(c.text)
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

func MarshalPrintCmdV2MUS(c PrintCmdV2, w muss.Writer) (n int, err error) {
	n, err = ord.MarshalString(c.from, w)
	if err != nil {
		return
	}
	var n1 int
	n1, err = ord.MarshalString(c.text, w)
	n += n1
	return
}

func UnmarshalPrintCmdV2MUS(r muss.Reader) (c PrintCmdV2, n int, err error) {
	c.from, n, err = ord.UnmarshalString(r)
	if err != nil {
		return
	}
	var n1 int
	c.text, n1, err = ord.UnmarshalString(r)
	n += n1
	return
}

func SizePrintCmdV2MUS(c PrintCmdV2) (size int) {
	size = ord.SizeString(c.from)
	return size + ord.SizeString(c.text)
}
