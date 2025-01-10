//go:generate go run gen/main.go
package exmpls

import (
	"context"
	"time"

	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
	com "github.com/mus-format/common-go"
)

const CmdExecDuration = time.Second

// DTMs help distinguish one command/result from another on the server/client
// side.
const (
	SayHelloCmdDTM com.DTM = iota
	SayFancyHelloCmdDTM
	UnsupportedCmdDTM
	SayFancyHelloMultiCmdDTM
	OldSayHelloCmdDTM
	ResultDTM
)

// NewSayHelloCmd reates new SayHelloCmd.
func NewSayHelloCmd(str string) SayHelloCmd {
	return SayHelloCmd{str}
}

// SayHelloCmd implements base.Cmd and Marshaller interfaces.
type SayHelloCmd struct {
	str string
}

func (c SayHelloCmd) Exec(ctx context.Context, at time.Time, seq base.Seq,
	receiver Greeter,
	proxy base.Proxy,
) error {
	deadline := at.Add(CmdExecDuration)
	// If the command needs to perform some context-related work:
	//
	//	 ownCtx, cancel := context.WithDeadline(ctx, deadline)
	//	 ...

	result := Result{
		str:     receiver.Join(receiver.Interjection(), c.str),
		lastOne: true,
	}
	return proxy.SendWithDeadline(deadline, seq, result)
}

func (c SayHelloCmd) Marshal(w transport.Writer) (err error) {
	_, err = SayHelloCmdDTS.Marshal(c, w)
	return
}

// NewSayFancyHelloCmd creates new SayFancyHelloCmd.
func NewSayFancyHelloCmd(str string) SayFancyHelloCmd {
	return SayFancyHelloCmd{str}
}

// SayFancyHelloCmd implements base.Cmd and Marshaller interfaces.
type SayFancyHelloCmd struct {
	str string
}

func (c SayFancyHelloCmd) Exec(ctx context.Context, at time.Time, seq base.Seq,
	receiver Greeter,
	proxy base.Proxy,
) error {
	var (
		result = Result{
			str: receiver.Join(receiver.Interjection(), receiver.Adjective(),
				c.str),
			lastOne: true,
		}
		deadline = at.Add(CmdExecDuration)
	)
	return proxy.SendWithDeadline(deadline, seq, result)
}

func (c SayFancyHelloCmd) Marshal(w transport.Writer) (err error) {
	_, err = SayFancyHelloCmdDTS.Marshal(c, w)
	return
}

// NewUnsupportedCmd creates a new UnsupportedCmd.
func NewUnsupportedCmd(str string) UnsupportedCmd {
	return UnsupportedCmd{str}
}

// UnsupportedCmd implements base.Cmd and Marshaller interfaces.
type UnsupportedCmd struct {
	str string
}

func (c UnsupportedCmd) Exec(ctx context.Context, at time.Time, seq base.Seq,
	receiver Greeter,
	proxy base.Proxy,
) error {
	var (
		result = Result{
			str: receiver.Join(receiver.Interjection(), receiver.Interjection(),
				c.str),
			lastOne: true,
		}
		deadline = at.Add(CmdExecDuration)
	)
	return proxy.SendWithDeadline(deadline, seq, result)
}

func (c UnsupportedCmd) Marshal(w transport.Writer) (err error) {
	_, err = UnsupportedCmdDTS.Marshal(c, w)
	return
}

// NewSayFancyHelloMultiCmd creates a new SayFancyHelloMultiCmd.
func NewSayFancyHelloMultiCmd(str string) SayFancyHelloMultiCmd {
	return SayFancyHelloMultiCmd{str}
}

// SayFancyHelloMultiCmd implements base.Cmd and Marshaller interfaces.
type SayFancyHelloMultiCmd struct {
	str string
}

func (c SayFancyHelloMultiCmd) Exec(ctx context.Context, at time.Time,
	seq base.Seq,
	receiver Greeter,
	proxy base.Proxy,
) (err error) {
	result := NewResult(receiver.Interjection(), false)
	err = proxy.Send(seq, result)
	if err != nil {
		return
	}
	result = NewResult(receiver.Adjective(), false)
	err = proxy.Send(seq, result)
	if err != nil {
		return
	}
	result = NewResult(c.str, true)
	deadline := at.Add(SendCmdDeadline)
	return proxy.SendWithDeadline(deadline, seq, result)
}

func (c SayFancyHelloMultiCmd) Marshal(w transport.Writer) (err error) {
	_, err = SayFancyHelloMultiCmdDTS.Marshal(c, w)
	return
}

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
	result := NewResult(receiver.SayHello(c.str), true)
	return proxy.Send(seq, result)
}

func (c OldSayHelloCmd) Marshal(w transport.Writer) (err error) {
	_, err = OldSayHelloCmdDTS.Marshal(c, w)
	return
}

// Migrate is used by the server Codec.
func (c OldSayHelloCmd) Migrate() SayHelloCmd {
	return NewSayHelloCmd(c.str)
}
