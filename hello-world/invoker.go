package hw

import (
	"context"
	"time"

	"github.com/cmd-stream/base-go"
)

// NewInvoker creates a new Invoker.
func NewInvoker[T any](receiver T) Invoker[T] {
	return Invoker[T]{receiver}
}

// Invoker does nothing but execute the Command.
type Invoker[T any] struct {
	receiver T
}

func (i Invoker[T]) Invoke(ctx context.Context, at time.Time, seq base.Seq,
	cmd base.Cmd[T],
	proxy base.Proxy,
) error {
	return cmd.Exec(ctx, at, seq, i.receiver, proxy)
}
