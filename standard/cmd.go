package main

import (
	"context"
	"time"

	"github.com/cmd-stream/base-go"
	muss "github.com/mus-format/mus-stream-go"
)

// Eq1Cmd represents an equation (a + b + c) which can be calculated on the
// server. Note, that each command should implement the base.Cmd interface.
// Also it implements the MarshallerMUS interface.
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

func (c Eq1Cmd) MarshalMUS(w muss.Writer) (n int, err error) {
	return Eq1DTS.Marshal(c, w)
}

// Eq2Cmd represents an equation (a - b - c) which can be calculated on the
// server. It implements the MarshallerMUS interface.
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

func (c Eq2Cmd) MarshalMUS(w muss.Writer) (n int, err error) {
	return Eq2DTS.Marshal(c, w)
}
