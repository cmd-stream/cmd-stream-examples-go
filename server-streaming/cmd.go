package streaming

import (
	"context"
	"time"

	"github.com/cmd-stream/base-go"
	com "github.com/mus-format/common-go"
	muss "github.com/mus-format/mus-stream-go"

	hw "github.com/cmd-stream/cmd-stream-examples-go/hello-world"
)

const SayFancyHelloMultiCmdDTM com.DTM = iota + 10

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
	receiver hw.Greeter,
	proxy base.Proxy,
) (err error) {
	var (
		deadline = at.Add(hw.CmdSendDuration)
		result   = NewGreeting(receiver.Interjection(), false)
	)
	err = proxy.SendWithDeadline(deadline, seq, result)
	if err != nil {
		return
	}
	result = NewGreeting(receiver.Adjective(), false)
	err = proxy.Send(seq, result)
	if err != nil {
		return
	}
	result = NewGreeting(c.str, true)
	return proxy.Send(seq, result)
}

func (c SayFancyHelloMultiCmd) MarshalTypedMUS(w muss.Writer) (n int, err error) {
	return SayFancyHelloMultiCmdDTS.Marshal(c, w)
}

func (c SayFancyHelloMultiCmd) SizeTypedMUS() (size int) {
	return SayFancyHelloMultiCmdDTS.Size(c)
}
