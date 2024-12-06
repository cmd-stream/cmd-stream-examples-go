package main

import (
	"context"
	"sync"

	examples "github.com/cmd-stream/cmd-stream-examples-go"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

// This example shows how RPC can be implemented using cmd-stream-go.
//
// It initialize EchoService with the cmd-stream-go client and then makes an RPC
// call - EchoService.Echo().
//
// Here we have struct{} as the receiver and examples.EchoCmd as a command.
func main() {
	const addr = "127.0.0.1:9000"

	// Start the server.
	wgS := &sync.WaitGroup{}
	server, err := examples.StartServer(addr, examples.ServerCodec{}, struct{}{},
		wgS)
	assert.EqualError(err, nil)

	// Create the client.
	client, err := examples.CreateClient(addr, examples.ClientCodec{})
	assert.EqualError(err, nil)

	// Create the service.
	service := EchoServiceImpl{client}

	// Make an RPC call.
	str, err := service.Echo(context.Background(), "hello world")
	assert.EqualError(err, nil)
	assert.Equal[string](str, "hello world")

	// Close the client.
	err = examples.CloseClient(client)
	assert.EqualError(err, nil)

	// Close the server.
	err = examples.CloseServer(server, wgS)
	assert.Equal(err, nil)
}
