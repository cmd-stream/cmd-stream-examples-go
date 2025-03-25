package rpc

import (
	"context"
	"sync"
	"testing"

	hw "cmd-stream-examples-go/hello-world"

	assert_error "github.com/ymz-ncnk/assert/error"
	assert_fatal "github.com/ymz-ncnk/assert/fatal"
)

// The rpc example demonstrates how to implement RPC using cmd-stream-go.
//
// Here, you'll find a GreetingService with a SayHello method that sends the
// corresponding SayHelloCmd to the server.
func TestRPC(t *testing.T) {
	const addr = "127.0.0.1:9005"

	// Start the server.
	wgS := &sync.WaitGroup{}
	server, err := hw.StartServer(addr, hw.ServerCodec{},
		hw.NewGreeter("Hello", "incredible", " "),
		wgS)
	assert_fatal.EqualError(err, nil, t)

	MakeRPC_Call(addr, t)

	// Close the server.
	err = hw.CloseServer(server, wgS)
	assert_fatal.Equal(err, nil, t)
}

func MakeRPC_Call(addr string, t *testing.T) {
	// Create the client.
	client, err := hw.CreateClient(addr, hw.ClientCodec{})
	assert_fatal.EqualError(err, nil, t)

	// Create the service.
	service := GreeterService{client}

	// Make an RPC call.
	str, err := service.SayHello(context.Background(), "world")
	assert_error.EqualError(err, nil, t)
	assert_error.Equal[string](str, "Hello world", t)

	// Close the client.
	err = hw.CloseClient(client)
	assert_fatal.EqualError(err, nil, t)
}
