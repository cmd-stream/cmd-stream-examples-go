package main

import (
	"context"
	"time"

	"github.com/cmd-stream/base-go"
	muss "github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/ord"
)

const EchoCount = 3

// It's a command with several results.
type MultiEchoCmd string

func (c MultiEchoCmd) Exec(ctx context.Context, at time.Time, seq base.Seq,
	receiver any,
	proxy base.Proxy,
) (err error) {
	for i := 0; i < EchoCount-1; i++ {
		err = proxy.Send(seq, MultiEchoResult{c, false})
		if err != nil {
			return
		}
	}
	return proxy.Send(seq, MultiEchoResult{c, true})
}

func MarshalMultiEchoCmdMUS(c MultiEchoCmd, w muss.Writer) (n int, err error) {
	return ord.MarshalString(string(c), nil, w)
}

func UnmarshalMultiEchoCmdMUS(r muss.Reader) (c MultiEchoCmd, n int, err error) {
	str, n, err := ord.UnmarshalString(nil, r)
	c = MultiEchoCmd(str)
	return
}

func SizeMultiEchoCmdMUS(c MultiEchoCmd) (size int) {
	return ord.SizeString(string(c), nil)
}
