package main

import (
	"context"

	"github.com/cmd-stream/base-go"
	base_client "github.com/cmd-stream/base-go/client"
	exmpls "github.com/cmd-stream/cmd-stream-examples-go"
)

type GreeterService struct {
	client *base_client.Client[exmpls.Greeter]
}

func (s GreeterService) SayHello(ctx context.Context, str string) (string, error) {
	cmd := exmpls.NewSayHelloCmd(str)

	result, err := SendCmd[exmpls.Greeter, exmpls.Result](ctx, cmd, s.client)
	if err != nil {
		return "", err
	}
	return result.Str(), nil
}

func SendCmd[T, R any](ctx context.Context, cmd base.Cmd[T],
	client *base_client.Client[T]) (result R, err error) {
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
