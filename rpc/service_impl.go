package main

import (
	"context"

	"github.com/cmd-stream/base-go"
	base_client "github.com/cmd-stream/base-go/client"
	examples "github.com/cmd-stream/cmd-stream-examples-go"
)

type EchoServiceImpl struct {
	client *base_client.Client[struct{}]
}

func (s EchoServiceImpl) Echo(ctx context.Context, str string) (string, error) {
	var (
		seq     base.Seq
		cmd     = examples.EchoCmd(str)
		results = make(chan base.AsyncResult)
		err     error
	)
	if seq, err = s.client.Send(cmd, results); err != nil {
		return "", err
	}
	select {
	case <-ctx.Done():
		// The command was not completed in the allotted time.
		s.client.Forget(seq)
		return "", context.Canceled
	case result := <-results:
		if result.Error != nil {
			return "", result.Error
		}
		return string(result.Result.(examples.OneEchoResult)), nil
	}
}
