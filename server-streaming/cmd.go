package streaming

import (
	"context"
	"time"

	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"

	hw "github.com/cmd-stream/cmd-stream-examples-go/hello-world"
)

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
		result   = NewResult(receiver.Interjection(), false)
	)
	err = proxy.SendWithDeadline(deadline, seq, result)
	if err != nil {
		return
	}
	result = NewResult(receiver.Adjective(), false)
	err = proxy.Send(seq, result)
	if err != nil {
		return
	}
	result = NewResult(c.str, true)
	return proxy.Send(seq, result)
}

func (c SayFancyHelloMultiCmd) Marshal(w transport.Writer) (err error) {
	_, err = SayFancyHelloMultiCmdDTS.Marshal(c, w)
	return
}
