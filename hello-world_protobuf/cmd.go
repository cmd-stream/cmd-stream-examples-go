package hwp

import (
	"context"
	"time"

	hw "github.com/cmd-stream/cmd-stream-examples-go/hello-world"
	muss "github.com/mus-format/mus-stream-go"

	"github.com/cmd-stream/base-go"
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
		result   = NewGreeting(str)
		deadline = at.Add(hw.CmdSendDuration)
	)
	return proxy.SendWithDeadline(deadline, seq, result)
}

func (c SayHelloCmd) MarshalTypedProtobuf(w muss.Writer) (n int,
	err error) {
	return SayHelloCmdDTS.Marshal(c, w)
}

func (c SayHelloCmd) SizeTypedProtobuf() (size int) {
	return SayHelloCmdDTS.Size(c)
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
		result   = NewGreeting(str)
		deadline = at.Add(hw.CmdSendDuration)
	)
	return proxy.SendWithDeadline(deadline, seq, result)
}

func (c SayFancyHelloCmd) MarshalTypedProtobuf(w muss.Writer) (n int,
	err error) {
	return SayFancyHelloCmdDTS.Marshal(c, w)
}

func (c SayFancyHelloCmd) SizeTypedProtobuf() (size int) {
	return SayFancyHelloCmdDTS.Size(c)
}
