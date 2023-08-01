package main

import (
	"net"
	"sync"
	"time"

	"github.com/cmd-stream/base-go"
	cs_client "github.com/cmd-stream/cmd-stream-go/client"
	cs_server "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/ymz-ncnk/assert"
)

func init() {
	assert.On = true
}

const Addr = "127.0.0.1:9000"

// This example demonstrates the standard usage of the client and server.
//
// Here we have the Calculator as the receiver, and Eq1Cmd, Eq2Cmd as commands.
func main() {
	listener, err := net.Listen("tcp", Addr)
	assert.EqualError(err, nil)

	// Start the server.
	var (
		wgS    = &sync.WaitGroup{}
		server = cs_server.NewDef[Calculator](ServerCodec{}, Calculator{})
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
	client, err := cs_client.NewDef[Calculator](ClientCodec{}, conn, nil)
	assert.EqualError(err, nil)

	// Stop the client.
	defer func() {
		err := client.Close()
		assert.EqualError(err, nil)
		// Wait for the client to stop.
		<-client.Done()
	}()

	wgR := &sync.WaitGroup{}

	// Send the first command with a timeout.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		var (
			seq     base.Seq
			cmd     = Eq1Cmd{1, 2, 3}
			results = make(chan base.AsyncResult, 1)
		)
		seq, err = client.Send(cmd, results)
		assert.EqualError(err, nil)

		select {
		case <-time.NewTimer(time.Second).C:
			client.Forget(seq) // If we are no longer interested in the results of
			// this command, we should call Forget.
		case asyncResult := <-results:
			result := (asyncResult).Result.(Result)
			assert.Equal[Result](result, Result(6))
		}
	}()

	// Send a second command.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		var (
			cmd     = Eq2Cmd{10, 2, 3}
			results = make(chan base.AsyncResult, 1)
		)
		_, err = client.Send(cmd, results)
		assert.EqualError(err, nil)

		result := (<-results).Result.(Result)
		assert.Equal[Result](result, Result(5))
	}()

	wgR.Wait()
}
