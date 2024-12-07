package main

import (
	"sync"

	"github.com/cmd-stream/base-go"
	examples "github.com/cmd-stream/cmd-stream-examples-go"
	"github.com/cmd-stream/cmd-stream-examples-go/cmd_versioning/old_client"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

// In this example there are two clients:
//   - Old client - sends a command that is current for this client and old for
//     the server.
//   - Current client - sends a command that is current for both this client and
//     the server.
//
// To process an older version of a command on the server, ServerCodec.Decode()
// migrates it to the current version.
//
// Here we have Printer as the receiver, PrintCmdV1 (old version) and PrintCmdV2
// (current version) as commands.
func main() {
	const addr = "127.0.0.1:9000"

	// Start the server.
	wgS := &sync.WaitGroup{}
	server, err := examples.StartServer(addr, ServerCodec{}, Printer{}, wgS)
	assert.EqualError(err, nil)

	// Create the old client.
	oldClient, err := old_client.CreateClient(addr)
	assert.EqualError(err, nil)

	// Create the current client.
	client, err := examples.CreateClient(addr, ClientCodec{})
	assert.EqualError(err, nil)

	// Old client sends the old version of the command.
	wgR := &sync.WaitGroup{}
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		results := make(chan base.AsyncResult, 1)

		_, err = oldClient.Send(old_client.PrintCmd{Text: "old"}, results)
		assert.EqualError(err, nil)

		result := (<-results).Result.(old_client.OkResult)
		assert.Equal[old_client.OkResult](result, old_client.OkResult(true))
	}()

	// Current client sends the current version of the command.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		results := make(chan base.AsyncResult, 1)

		_, err = client.Send(PrintCmd{From: "me", Text: "hello"}, results)
		assert.EqualError(err, nil)

		result := (<-results).Result.(OkResult)
		assert.Equal[OkResult](result, OkResult(true))
	}()
	wgR.Wait()

	// Close the old client.
	err = examples.CloseClient(oldClient)
	assert.EqualError(err, nil)

	// Close the current client.
	err = examples.CloseClient(client)
	assert.EqualError(err, nil)

	// Close the server.
	err = examples.CloseServer(server, wgS)
	assert.EqualError(err, nil)
}
