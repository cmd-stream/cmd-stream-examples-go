package main

import "context"

type EchoService interface {
	Echo(ctx context.Context, str string) (string, error)
}
