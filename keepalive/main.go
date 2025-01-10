package main

import (
	"net"
	"sync"
	"time"

	base_client "github.com/cmd-stream/base-go/client"
	exmpls "github.com/cmd-stream/cmd-stream-examples-go"
	cs_client "github.com/cmd-stream/cmd-stream-go/client"
	delegate_client "github.com/cmd-stream/delegate-go/client"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

// In this example, the client attempts to keep the connection to the server
// alive even when there are no commands to send.
//
// The client is initialized with the following configuration:
//
//	conf := cs_client.Conf {
//		Delegate: delegate_client.Conf{
//			KeepaliveTime:  200 * time.Millisecond,
//			KeepaliveIntvl: 200 * time.Millisecond,
//		},
//	}
func main() {
	const addr = "127.0.0.1:9000"

	// Start the server.
	wgS := &sync.WaitGroup{}
	server, err := exmpls.StartServer(addr, exmpls.ServerCodec{},
		exmpls.NewGreeter("Hello", "incredible", " "), wgS)
	assert.EqualError(err, nil)

	// Create the client.
	client, err := CreateKeepaliveClient(addr, exmpls.ClientCodec{})
	assert.EqualError(err, nil)

	wgC := &sync.WaitGroup{}
	timeout := time.Second

	// Send a command.
	wgC.Add(1)
	cmd := exmpls.NewSayHelloCmd("world")
	wantResults := []exmpls.Result{exmpls.NewResult("Hello world", true)}
	exmpls.SendCmd(cmd, timeout, nil, wantResults, exmpls.CompareResults,
		client, wgC)

	// Ping-Pong time... When there are no commands to send, the client sends
	// a predefined PingCmd.
	time.Sleep(2 * exmpls.ReceiveTimeout)

	// Send a command again after the long delay to check if the connection is
	// still alive.
	wgC.Add(1)
	cmd = exmpls.NewSayHelloCmd("world again")
	wantResults = []exmpls.Result{exmpls.NewResult("Hello world again", true)}
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

// CreateKeepaliveClient creates a client that will keep a connection alive
// when there are no commands to send.
func CreateKeepaliveClient[T any](addr string, codec cs_client.Codec[T]) (
	c *base_client.Client[T], err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	conf := cs_client.Conf{
		Delegate: delegate_client.Conf{
			KeepaliveTime:  200 * time.Millisecond,
			KeepaliveIntvl: 200 * time.Millisecond,
		},
	}
	return exmpls.CreateClientWith(conf, codec, conn)
}
