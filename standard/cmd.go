package main

import (
	"context"
	"time"

	"github.com/cmd-stream/base-go"
	muss "github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/raw"
	"github.com/mus-format/mus-stream-go/varint"
)

const (
	Eq1CmdType = iota + 1
	Eq2CmdType
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
// Eq1Cmd represents an equation (a + b + c) which can be calculated on the
// server. Note, that each command should implement the base.Cmd interface.
type Eq1Cmd struct {
	a int
	b int
	c int
}

func (c Eq1Cmd) Exec(ctx context.Context, at time.Time, seq base.Seq,
	receiver Calculator,
	proxy base.Proxy,
) error {
	// If the command execution has a timeout, you can use here a ctx and the
	// proxy.SendWithDeadline method, also server should be configured with
	// Conf.Handler.At == true:
	//
	// deadline := at.Add(timeout)
	// ownCtx, cancel := context.WithDeadline(ctx, deadline)
	// // Do the context-related work.
	// ...
	// return proxy.SendWithDeadline(seq, result, deadline)

	result := Result(receiver.Add(receiver.Add(c.a, c.b), c.c))
	// With help of the proxy, the command sends back results.
	return proxy.Send(seq, result)
}

func MarshalEq1CmdMUS(c Eq1Cmd, w muss.Writer) (n int, err error) {
	n, err = varint.MarshalInt(c.a, w)
	if err != nil {
		return
	}
	var n1 int
	n1, err = varint.MarshalInt(c.b, w)
	n += n1
	if err != nil {
		return
	}
	n1, err = varint.MarshalInt(c.c, w)
	n += n1
	return
}

func UnmarshalEq1CmdMUS(r muss.Reader) (c Eq1Cmd, n int, err error) {
	c.a, n, err = varint.UnmarshalInt(r)
	if err != nil {
		return
	}
	var n1 int
	c.b, n, err = varint.UnmarshalInt(r)
	n += n1
	if err != nil {
		return
	}
	c.c, n, err = varint.UnmarshalInt(r)
	n += n1
	return
}

func SizeEq1CmdMUS(c Eq1Cmd) (size int) {
	size = varint.SizeInt(c.a)
	size += varint.SizeInt(c.b)
	return size + varint.SizeInt(c.c)
}

// -----------------------------------------------------------------------------
// Eq2Cmd represents an equation (a - b - c) which can be calculated on the
// server.
type Eq2Cmd struct {
	a int
	b int
	c int
}

func (c Eq2Cmd) Exec(ctx context.Context, at time.Time, seq base.Seq,
	receiver Calculator,
	proxy base.Proxy,
) error {
	result := Result(receiver.Sub(receiver.Sub(c.a, c.b), c.c))
	return proxy.Send(seq, result)
}

func MarshalEq2CmdMUS(c Eq2Cmd, w muss.Writer) (n int, err error) {
	n, err = varint.MarshalInt(c.a, w)
	if err != nil {
		return
	}
	var n1 int
	n1, err = varint.MarshalInt(c.b, w)
	n += n1
	if err != nil {
		return
	}
	n1, err = varint.MarshalInt(c.c, w)
	n += n1
	return
}

func UnmarshalEq2CmdMUS(r muss.Reader) (c Eq2Cmd, n int, err error) {
	c.a, n, err = varint.UnmarshalInt(r)
	if err != nil {
		return
	}
	var n1 int
	c.b, n, err = varint.UnmarshalInt(r)
	n += n1
	if err != nil {
		return
	}
	c.c, n, err = varint.UnmarshalInt(r)
	n += n1
	return
}

func SizeEq2CmdMUS(c Eq2Cmd) (size int) {
	size = varint.SizeInt(c.a)
	size += varint.SizeInt(c.b)
	return size + varint.SizeInt(c.c)
}
