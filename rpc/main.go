package main

import (
	"context"
	"net"
	"sync"

	examples "github.com/cmd-stream/cmd-stream-examples-go"
	cs_client "github.com/cmd-stream/cmd-stream-go/client"
	cs_server "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/ymz-ncnk/assert"
)

func init() {
	assert.On = true
}

const Addr = "127.0.0.1:9000"

type connFactory struct{}

func (f connFactory) New() (net.Conn, error) {
	return net.Dial("tcp", Addr)
}

// This example shows how you can implement RPC using the cmd-stream-go.
//
// Here we have struct{} as the receiver, and examples.EchoCmd as a command.
func main() {
	listener, err := net.Listen("tcp", Addr)
	if err != nil {
		return
	}

	// Start the server.
	var (
		wgS    = &sync.WaitGroup{}
		server = cs_server.NewDef[struct{}](examples.ServerCodec{}, struct{}{})
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

	// Create the client.
	client, err := cs_client.NewDefReconnect[struct{}](examples.ClientCodec{},
		connFactory{},
		nil)
	assert.EqualError(err, nil)

	// Stop the client.
	defer func() {
		err := client.Close()
		assert.EqualError(err, nil)
		// Wait for the client to stop.
		<-client.Done()
	}()

	// Use the service.
	service := EchoServiceImpl{client}

	str, err := service.Echo(context.Background(), "hello world")
	assert.EqualError(err, nil)
	assert.Equal[string](str, "hello world")
}
