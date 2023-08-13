package main

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/cmd-stream/base-go"
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

// In this example, the client is trying to reconnect to the server when the
// connection has been lost.
//
// To do this, we create the client using the cs_client.NewDefReconnect method.
//
// Here we have struct{} as the receiver, and examples.EchoCmd as a command.
func main() {
	listener, err := net.Listen("tcp", Addr)
	assert.EqualError(err, nil)

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

	// Connect to the server.
	client, err := cs_client.NewDefReconnect[struct{}](examples.ClientCodec{},
		connFactory{},
		nil)
	assert.EqualError(err, nil)

	// Stop the server.
	err = server.Close()
	assert.EqualError(err, nil)
	wgS.Wait()

	// Start the server again after some time.
	time.Sleep(time.Second)
	listener, err = net.Listen("tcp", Addr)
	assert.EqualError(err, nil)

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
	assert.Equal[examples.OneEchoResult](result,
		examples.OneEchoResult(cmd))
}
