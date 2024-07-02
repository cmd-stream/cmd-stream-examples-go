package main

import (
	"net"
	"sync"

	"github.com/cmd-stream/base-go"
	cs_client "github.com/cmd-stream/cmd-stream-go/client"
	cs_server "github.com/cmd-stream/cmd-stream-go/server"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

const Addr = "127.0.0.1:9000"

// This example shows how you can use different versions of the same command,
// for example, to support old clients.
//
// With the new version of the command, the receiver usually changes so that the
// old version is no longer compatible with it. To fix this, we must change the
// Exec method of the old command - place data migration in it.
//
// Here we have Printer as the receiver and PrintCmdV1 (old version),
// PrintCmdV2 (current version) as commands.
func main() {
	listener, err := net.Listen("tcp", Addr)
	assert.EqualError(err, nil)

	// Start the server.
	var (
		wgS    = &sync.WaitGroup{}
		server = cs_server.NewDef[Printer](ServerCodec{}, Printer{})
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
	client, err := cs_client.NewDef[Printer](ClientCodec{}, conn, nil)
	assert.EqualError(err, nil)

	// Stop the client.
	defer func() {
		err := client.Close()
		assert.EqualError(err, nil)
		// Wait for the client to stop.
		<-client.Done()
	}()

	wgR := &sync.WaitGroup{}

	// Send the old version of the command.
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
}
