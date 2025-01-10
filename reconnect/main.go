package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/cmd-stream/base-go"
	base_client "github.com/cmd-stream/base-go/client"
	exmpls "github.com/cmd-stream/cmd-stream-examples-go"
	cs_client "github.com/cmd-stream/cmd-stream-go/client"
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

// In this example, the client attempts to reconnect to the server when the
// connection is lost.
//
// The client is created using the cs_client.NewDefReconnect() function.
func main() {
	// Start the server.
	wgS := &sync.WaitGroup{}
	server, err := exmpls.StartServer(Addr, exmpls.ServerCodec{},
		exmpls.NewGreeter("Hello", "incredible", " "), wgS)
	assert.EqualError(err, nil)

	// Create the client.
	client, err := CreateReconnectClient(exmpls.ClientCodec{},
		connFactory{})
	assert.EqualError(err, nil)

	// Close the server.
	err = exmpls.CloseServer(server, wgS)
	assert.EqualError(err, nil)

	// Start the server again after some time.
	time.Sleep(time.Second)
	server, err = exmpls.StartServer(Addr, exmpls.ServerCodec{},
		exmpls.NewGreeter("Hello", "incredible", " "), wgS)
	assert.EqualError(err, nil)

	// Wait for the client to reconnect.
	time.Sleep(200 * time.Millisecond)

	wgC := &sync.WaitGroup{}

	// Send a command.
	wgC.Add(1)
	timeout := time.Second
	cmd := exmpls.NewSayHelloCmd("world")
	wantResults := []exmpls.Result{exmpls.NewResult("Hello world", true)}
	exmpls.SendCmd(cmd, timeout, nil, wantResults, exmpls.CompareResults, client,
		wgC)

	wgC.Wait()
	// Close the client.
	err = exmpls.CloseClient(client)
	assert.EqualError(err, nil)

	// Close the server.
	err = exmpls.CloseServer(server, wgS)
	assert.EqualError(err, nil)
}

func CreateReconnectClient[T any](codec cs_client.Codec[T],
	connFactory cs_client.ConnFactory) (c *base_client.Client[T], err error) {
	// unexpectedResultHandler processes unexpected results from the server.
	// If you call Client.Forget(seq) for a command, its results will be handled
	// by this function.
	unexpectedResultHandler := func(seq base.Seq, result base.Result) {
		fmt.Printf("unexpected result was received: seq %v, result %v\n", seq,
			result)
	}
	return cs_client.NewDefReconnect[T](codec,
		connFactory,
		unexpectedResultHandler)
}
