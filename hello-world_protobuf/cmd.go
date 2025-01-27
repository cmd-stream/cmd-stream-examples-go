package hwp

import (
	hw "cmd-stream-examples-go/hello-world"
	"context"
	"time"

	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
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
	receiver hw.Greeter,
	proxy base.Proxy,
) error {
	var (
		str      = receiver.Join(receiver.Interjection(), c.Str)
		result   = NewResult(str)
		deadline = at.Add(hw.CmdSendDuration)
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
	receiver hw.Greeter,
	proxy base.Proxy,
) error {
	var (
		str      = receiver.Join(receiver.Interjection(), receiver.Adjective(), c.Str)
		result   = NewResult(str)
		deadline = at.Add(hw.CmdSendDuration)
	)
	return proxy.SendWithDeadline(deadline, seq, result)
}

func (c SayFancyHelloCmd) Marshal(w transport.Writer) (err error) {
	_, err = SayFancyHelloCmdDTS.Marshal(c, w)
	return
}
