package main

import (
	"context"
	"time"

	"github.com/cmd-stream/base-go"
)

func NewEq1Cmd(a, b, c int64) Eq1Cmd {
	return Eq1Cmd{&Eq1Data{A: a, B: b, C: c}}
}

// Eq1Cmd represents an equation (a + b + c) which can be calculated on the
// server. Note, that each command should implement the base.Cmd interface.
type Eq1Cmd struct {
	*Eq1Data
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

	result := Result{
		ResultData: &ResultData{
			R: receiver.Add(receiver.Add(c.Eq1Data.A, c.Eq1Data.B), c.Eq1Data.C),
		},
	}

	// With help of the proxy, the command sends back results.
	return proxy.Send(seq, result)
}

func NewEq2Cmd(a, b, c int64) Eq2Cmd {
	return Eq2Cmd{&Eq2Data{A: a, B: b, C: c}}
}

// Eq2Cmd represents an equation (a - b - c) which can be calculated on the
// server.
type Eq2Cmd struct {
	*Eq2Data
}

func (c Eq2Cmd) Exec(ctx context.Context, at time.Time, seq base.Seq,
	receiver Calculator,
	proxy base.Proxy,
) error {
	result := Result{
		ResultData: &ResultData{
			R: receiver.Sub(receiver.Sub(c.Eq2Data.A, c.Eq2Data.B), c.Eq2Data.C),
		},
	}
	return proxy.Send(seq, result)
}
