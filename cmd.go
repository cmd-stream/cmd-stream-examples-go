package examples

import (
	"context"
	"time"

	"github.com/cmd-stream/base-go"
	muss "github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/ord"
)

type EchoCmd string

func (c EchoCmd) Exec(ctx context.Context, at time.Time, seq base.Seq,
	receiver struct{},
	proxy base.Proxy,
) (err error) {
	return proxy.Send(seq, c)
}

func (e EchoCmd) LastOne() bool {
	return true
}

func MarshalEchoCmdMUS(c EchoCmd, w muss.Writer) (n int, err error) {
	return ord.MarshalString(string(c), nil, w)
}

func UnmarshalEchoCmdMUS(r muss.Reader) (c EchoCmd, n int, err error) {
	str, n, err := ord.UnmarshalString(nil, r)
	c = EchoCmd(str)
	return
}

func SizeEchoCmdMUS(c EchoCmd) (size int) {
	return ord.SizeString(string(c), nil)
}
