package main

import (
	"sync"
	"time"

	"github.com/cmd-stream/base-go"
	base_client "github.com/cmd-stream/base-go/client"
	examples "github.com/cmd-stream/cmd-stream-examples-go"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

// This example demonstrates the standard use of cmd-stream-go with the MUS
// serializer.
//
// Here we have Calculator as the receiver and Eq1Cmd, Eq2Cmd as commands. The
// other files in this package also have useful comments, so check them as well.
func main() {
	const addr = "127.0.0.1:9000"

	// Start the server.
	wgS := &sync.WaitGroup{}
	server, err := examples.StartServer(addr, ServerCodec{}, Calculator{}, wgS)
	assert.EqualError(err, nil)

	// Create the client.
	client, err := examples.CreateClient(addr, ClientCodec{})
	assert.EqualError(err, nil)

	// Execute two commands and wait for both to complete.
	wgR := &sync.WaitGroup{}
	wgR.Add(2)
	go sendCmd(wgR, client)
	go sendCmdWithTimeout(wgR, client)
	wgR.Wait()

	// Close the client.
	err = examples.CloseClient(client)
	assert.EqualError(err, nil)

	// Close the server.
	err = examples.CloseServer(server, wgS)
	assert.EqualError(err, nil)
}

func sendCmd(wg *sync.WaitGroup, client *base_client.Client[Calculator]) {
	defer wg.Done()
	var (
		cmd            = Eq2Cmd{10, 2, 3}
		expectedResult = Result(5)
		asyncResults   = make(chan base.AsyncResult, 1)
	)
	_, err := client.Send(cmd, asyncResults)
	assert.EqualError(err, nil)

	asyncResult := <-asyncResults
	// asyncResult.Error != nil if something is wrong with the connection.
	assert.EqualError(asyncResult.Error, nil)

	result := asyncResult.Result.(Result)
	assert.Equal(result, expectedResult)
}

func sendCmdWithTimeout(wg *sync.WaitGroup,
	client *base_client.Client[Calculator]) {
	defer wg.Done()
	var (
		seq            base.Seq
		cmd            = Eq1Cmd{1, 2, 3}
		expectedResult = Result(6)
		results        = make(chan base.AsyncResult, 1)
	)
	seq, err := client.Send(cmd, results)
	assert.EqualError(err, nil)
	// Wait for the result.
	select {
	case <-time.NewTimer(time.Second).C:
		client.Forget(seq) // If we are no longer interested in the results of
		// this command, we should call Forget().
	case asyncResult := <-results:
		// asyncResult.Error != nil if something is wrong with the connection.
		assert.EqualError(asyncResult.Error, nil)

		result := asyncResult.Result.(Result)
		assert.Equal(result, expectedResult)
	}
}
