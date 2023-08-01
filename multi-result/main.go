package main

import (
	"net"
	"sync"

	"github.com/cmd-stream/base-go"
	cs_client "github.com/cmd-stream/cmd-stream-go/client"
	cs_server "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/ymz-ncnk/assert"
)

func init() {
	assert.On = true
}

const Addr = "127.0.0.1:9000"

// In this example, we send one command and get several results. The client can
// recognize the last command result by the Result.LastOne method.
//
// Here we have struct{} as the receiver, and MultiEchoCmd as a command.
func main() {
	listener, err := net.Listen("tcp", Addr)
	assert.EqualError(err, nil)

	// Start the server.
	var (
		wgS    = &sync.WaitGroup{}
		server = cs_server.NewDef[any](ServerCodec{}, nil)
	)
	wgS.Add(1)
	go func() {
		defer wgS.Done()
		server.Serve(listener.(*net.TCPListener))
	}()

	// Stop the server.
	defer func() {
		err := server.Close()
		assert.EqualError(err, nil)
		wgS.Wait()
	}()

	// Connect to the server.
	conn, err := net.Dial("tcp", Addr)
	assert.EqualError(err, nil)

	// Create the client.
	wgR := &sync.WaitGroup{}
	client, err := cs_client.NewDef[any](ClientCodec{}, conn, nil)
	assert.EqualError(err, nil)

	// Stop the client.
	defer func() {
		wgR.Wait()
		err := client.Close()
		assert.EqualError(err, nil)
		// Wait for the client to stop.
		<-client.Done()
	}()

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
}
