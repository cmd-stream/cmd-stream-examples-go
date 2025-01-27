package short

import (
	"context"
	"time"

	"github.com/cmd-stream/base-go"
)

// EchoCmd sends its content back as Result.
type EchoCmd string

func (c EchoCmd) Exec(ctx context.Context, at time.Time, seq base.Seq,
	receiver struct{}, proxy base.Proxy) error {
	return proxy.Send(seq, Result(c))
}
