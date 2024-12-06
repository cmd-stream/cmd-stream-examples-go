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

// In this example, we send one command and get several results. The client
// recognizes the last result using the Result.LastOne() method.
//
// Here we have struct{} as the receiver and MultiEchoCmd as a command.
func main() {
	const addr = "127.0.0.1:9000"

	// Start the server.
	wgS := &sync.WaitGroup{}
	server, err := examples.StartServer(addr, ServerCodec{}, nil, wgS)
	assert.EqualError(err, nil)

	// Create the client.
	client, err := examples.CreateClient(addr, ClientCodec{})
	assert.EqualError(err, nil)

	// Send a command.
	var (
		cmd     = MultiEchoCmd("hello world")
		results = make(chan base.AsyncResult, EchoCount)
	)
	_, err = client.Send(cmd, results)
	assert.EqualError(err, nil)

	// Receive the first result.
	result1 := (<-results).Result.(MultiEchoResult)
	assert.Equal[MultiEchoResult](result1, MultiEchoResult{cmd, false})

	// Receive the sedcond result.
	result2 := (<-results).Result.(MultiEchoResult)
	assert.Equal[MultiEchoResult](result2, MultiEchoResult{cmd, false})

	// Receive the third result.
	result3 := (<-results).Result.(MultiEchoResult)
	assert.Equal[MultiEchoResult](result3, MultiEchoResult{cmd, true})

	// Close the client.
	err = examples.CloseClient(client)
	assert.EqualError(err, nil)

	// Close the server.
	err = examples.CloseServer(server, wgS)
	assert.EqualError(err, nil)
}
