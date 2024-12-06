package main

import (
	"sync"

	"github.com/cmd-stream/base-go"
	examples "github.com/cmd-stream/cmd-stream-examples-go"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

// This example shows how the server can handle different versions of a command.
//
// Changes of the receiver required to support a new version of the command may
// break backward compatibility. In this case, we need to adapt old command
// versions (only on the server) to the new realities, since there is only one
// receiver on the server.
//
// Here we have Printer as the receiver, PrintCmdV1 (old adapted version) and
// PrintCmdV2 (current version) as commands.
func main() {
	const addr = "127.0.0.1:9000"

	// Start the server.
	wgS := &sync.WaitGroup{}
	server, err := examples.StartServer(addr, ServerCodec{}, Printer{}, wgS)
	assert.EqualError(err, nil)

	// Create the client.
	client, err := examples.CreateClient(addr, ClientCodec{})
	assert.EqualError(err, nil)

	// Send the old version of the command.
	wgR := &sync.WaitGroup{}
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		var (
			cmd     = PrintCmdV1{text: "hello world"}
			results = make(chan base.AsyncResult, 1)
		)
		_, err = client.Send(cmd, results)
		assert.EqualError(err, nil)
		result := (<-results).Result.(OkResult)
		assert.Equal[OkResult](result, OkResult(true))
	}()

	// Send the current version of the command.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		var (
			cmd     = PrintCmdV2{from: "me", text: "hello world"}
			results = make(chan base.AsyncResult, 1)
		)
		_, err = client.Send(cmd, results)
		assert.EqualError(err, nil)

		result := (<-results).Result.(OkResult)
		assert.Equal[OkResult](result, OkResult(true))
	}()
	wgR.Wait()

	// Close the client.
	err = examples.CloseClient(client)
	assert.EqualError(err, nil)

	// Close the server.
	err = examples.CloseServer(server, wgS)
	assert.EqualError(err, nil)
}
