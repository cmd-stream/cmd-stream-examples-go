package main

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/cmd-stream/base-go"
	examples "github.com/cmd-stream/cmd-stream-examples-go"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

const Addr = "127.0.0.1:9000"

type connFactory struct{}

func (f connFactory) New() (conn net.Conn, err error) {
	time.Sleep(100 * time.Millisecond)
	conn, err = net.Dial("tcp", Addr)
	if err != nil {
		log.Println("failed to reconnect")
	} else {
		log.Println("connected")
	}
	return
}

// This example demostrates the reconnection feature.
//
// It uses the cs_client.NewDefReconnect() function to create a client that can
// reconnect to the server if the connection is lost.
//
// Here we have struct{} as the receiver and examples.EchoCmd as a command.
func main() {
	// Start the server.
	wgS := &sync.WaitGroup{}
	server, err := examples.StartServer(Addr, examples.ServerCodec{}, struct{}{},
		wgS)
	assert.EqualError(err, nil)

	// Create the client.
	client, err := examples.CreateReconnectClient(examples.ClientCodec{},
		connFactory{})
	assert.EqualError(err, nil)

	// Close the server.
	err = examples.CloseServer(server, wgS)
	assert.EqualError(err, nil)

	// Start the server again after some time.
	time.Sleep(time.Second)
	server, err = examples.StartServer(Addr, examples.ServerCodec{}, struct{}{},
		wgS)
	assert.EqualError(err, nil)

	// Wait for the client to reconnect.
	time.Sleep(200 * time.Millisecond)

	// Send a command.
	var (
		cmd     = examples.EchoCmd("hello world")
		results = make(chan base.AsyncResult, 1)
	)
	_, err = client.Send(cmd, results)
	assert.EqualError(err, nil)

	result := (<-results).Result.(examples.OneEchoResult)
	assert.Equal[examples.OneEchoResult](result, examples.OneEchoResult(cmd))

	// Close the client.
	err = examples.CloseClient(client)
	assert.EqualError(err, nil)

	// Close the server.
	err = examples.CloseServer(server, wgS)
	assert.EqualError(err, nil)
}
