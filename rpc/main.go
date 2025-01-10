package main

import (
	"context"
	"sync"

	exmpls "github.com/cmd-stream/cmd-stream-examples-go"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

// The rpc example demonstrates how to implement RPC using cmd-stream-go.
//
// Here, you'll find a GreetingService with a SayHello method that sends the
// corresponding SayHelloCmd command to the server.
func main() {
	const addr = "127.0.0.1:9000"

	// Start the server.
	wgS := &sync.WaitGroup{}
	server, err := exmpls.StartServer(addr, exmpls.ServerCodec{},
		exmpls.NewGreeter("Hello", "incredible", " "),
		wgS)
	assert.EqualError(err, nil)

	MakeRPC_Call(addr)

	// Close the server.
	err = exmpls.CloseServer(server, wgS)
	assert.Equal(err, nil)
}

func MakeRPC_Call(addr string) {
	// Create the client.
	client, err := exmpls.CreateClient(addr, exmpls.ClientCodec{})
	assert.EqualError(err, nil)

	// Create the service.
	service := GreeterService{client}

	// Make an RPC call.
	str, err := service.SayHello(context.Background(), "world")
	assert.EqualError(err, nil)
	assert.Equal[string](str, "Hello world")

	// Close the client.
	err = exmpls.CloseClient(client)
	assert.EqualError(err, nil)
}
