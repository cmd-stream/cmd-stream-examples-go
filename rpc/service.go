package rpc

import (
	"context"

	hw "github.com/cmd-stream/cmd-stream-examples-go/hello-world"

	"github.com/cmd-stream/base-go"
	bcln "github.com/cmd-stream/base-go/client"
)

type GreeterService struct {
	client *bcln.Client[hw.Greeter]
}

func (s GreeterService) SayHello(ctx context.Context, str string) (string, error) {
	cmd := hw.NewSayHelloCmd(str)
	greeting, err := SendCmd[hw.Greeter, hw.Greeting](ctx, cmd, s.client)
	if err != nil {
		return "", err
	}
	return string(greeting), nil
}

func SendCmd[T, R any](ctx context.Context, cmd base.Cmd[T],
	client *bcln.Client[T]) (result R, err error) {
	var (
		seq     base.Seq
		results = make(chan base.AsyncResult)
	)
	seq, err = client.Send(cmd, results)
	if err != nil {
		return
	}
	select {
	case <-ctx.Done():
		client.Forget(seq)
		err = context.Canceled
		return
	case asyncResult := <-results:
		if asyncResult.Error != nil {
			err = asyncResult.Error
			return
		}
		result = asyncResult.Result.(R)
		return
	}
}
