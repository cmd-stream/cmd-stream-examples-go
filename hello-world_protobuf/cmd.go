package main

import (
	"context"
	"time"

	"github.com/cmd-stream/base-go"
	exmpls "github.com/cmd-stream/cmd-stream-examples-go"
	"github.com/cmd-stream/transport-go"
	com "github.com/mus-format/common-go"
)

// DTMs help distinguish one command/result from another on the server/client
// side.
const (
	SayHelloCmdDTM com.DTM = iota
	SayFancyHelloCmdDTM
	UnsupportedCmdDTM
	ResultDTM
)

// NewSayHelloCmd creates a new SayHelloCmd.
func NewSayHelloCmd(str string) SayHelloCmd {
	return SayHelloCmd{
		SayHelloData: &SayHelloData{Str: str},
	}
}

// SayHelloCmd implements base.Cmd and Marshaller interfaces.
type SayHelloCmd struct {
	*SayHelloData
}

func (c SayHelloCmd) Exec(ctx context.Context, at time.Time, seq base.Seq,
	receiver exmpls.Greeter,
	proxy base.Proxy,
) error {
	var (
		str      = receiver.Join(receiver.Interjection(), c.Str)
		result   = NewResult(str, true)
		deadline = at.Add(exmpls.SendCmdDeadline)
	)
	return proxy.SendWithDeadline(deadline, seq, result)
}

func (c SayHelloCmd) Marshal(w transport.Writer) (err error) {
	_, err = SayHelloCmdDTS.Marshal(c, w)
	return
}

// NewSayFancyHelloCmd creates a new SayFancyHelloCmd.
func NewSayFancyHelloCmd(str string) SayFancyHelloCmd {
	return SayFancyHelloCmd{
		SayFancyHelloData: &SayFancyHelloData{Str: str},
	}
}

// SayFancyHelloCmd implements base.Cmd and Marshaller interfaces.
type SayFancyHelloCmd struct {
	*SayFancyHelloData
}

func (c SayFancyHelloCmd) Exec(ctx context.Context, at time.Time, seq base.Seq,
	receiver exmpls.Greeter,
	proxy base.Proxy,
) error {
	var (
		str      = receiver.Join(receiver.Interjection(), receiver.Adjective(), c.Str)
		result   = NewResult(str, true)
		deadline = at.Add(exmpls.SendCmdDeadline)
	)
	return proxy.SendWithDeadline(deadline, seq, result)
}

func (c SayFancyHelloCmd) Marshal(w transport.Writer) (err error) {
	_, err = SayFancyHelloCmdDTS.Marshal(c, w)
	return
}

// NewUnsupportedCmd creates a new UnsupportedCmd.
func NewUnsupportedCmd(str string) UnsupportedCmd {
	return UnsupportedCmd{
		UnsupportedData: &UnsupportedData{Str: str},
	}
}

// UnsupportedCmd implements base.Cmd and Marshaller interfaces.
type UnsupportedCmd struct {
	*UnsupportedData
}

func (c UnsupportedCmd) Exec(ctx context.Context, at time.Time, seq base.Seq,
	receiver exmpls.Greeter,
	proxy base.Proxy,
) error {
	var (
		str      = receiver.Join(receiver.Adjective(), receiver.Interjection(), c.Str)
		result   = NewResult(str, true)
		deadline = at.Add(exmpls.SendCmdDeadline)
	)
	return proxy.SendWithDeadline(deadline, seq, result)
}

func (c UnsupportedCmd) Marshal(w transport.Writer) (err error) {
	_, err = UnsupportedCmdDTS.Marshal(c, w)
	return
}
