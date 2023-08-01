package main

import (
	"net"
	"sync"
	"time"

	"github.com/cmd-stream/base-go"
	base_server "github.com/cmd-stream/base-go/server"
	examples "github.com/cmd-stream/cmd-stream-examples-go"
	cs_client "github.com/cmd-stream/cmd-stream-go/client"
	cs_server "github.com/cmd-stream/cmd-stream-go/server"
	delegate "github.com/cmd-stream/delegate-go"
	delegate_client "github.com/cmd-stream/delegate-go/client"
	"github.com/cmd-stream/handler-go"
	"github.com/ymz-ncnk/assert"
)

func init() {
	assert.On = true
}

const Addr = "127.0.0.1:9000"

// In this example, the client is trying to keep alive the connection to the
// server, when there are no commands to send.
//
// Here we have struct{} as the receiver, and examples.EchoCmd as a command.
func main() {
	listener, err := net.Listen("tcp", Addr)
	assert.EqualError(err, nil)

	// Start the server.
	var (
		wgS        = &sync.WaitGroup{}
		confServer = cs_server.Conf{
			Base: base_server.Conf{
				WorkersCount: 1,
			},
			Handler: handler.Conf{
				ReceiveTimeout: 2 * time.Second, // If no commands are sent from the
				// client within 2 seconds, the server will close the connection.
			},
		}
		server = cs_server.New[struct{}](cs_server.DefServerInfo,
			delegate.ServerSettings{},
			confServer,
			examples.ServerCodec{},
			struct{}{},
			nil)
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
	// If both KeepaliveTime and KeepaliveIntvl != 0, the keeapalive mode will be
	// enabled.
	confClient := cs_client.Conf{
		Delegate: delegate_client.Conf{
			KeepaliveTime:  time.Second,
			KeepaliveIntvl: time.Second,
		},
	}
	client, err := cs_client.New[struct{}](cs_server.DefServerInfo, confClient,
		examples.ClientCodec{},
		conn,
		nil)
	assert.EqualError(err, nil)

	// Stop the client.
	defer func() {
		err := client.Close()
		assert.EqualError(err, nil)
		// Wait for the client to stop.
		<-client.Done()
	}()

	// Send a command.
	var (
		cmd     = examples.EchoCmd("hello world")
		results = make(chan base.AsyncResult, 1)
	)
	_, err = client.Send(cmd, results)
	assert.EqualError(err, nil)

	result := (<-results).Result.(examples.OneEchoResult)
	assert.Equal[examples.OneEchoResult](result, examples.OneEchoResult(cmd))

	// Ping-Pong time...
	time.Sleep(confServer.Handler.ReceiveTimeout + 5*time.Second)

	// Send a command again after a long delay.
	cmd = examples.EchoCmd("hello world again")
	_, err = client.Send(cmd, results)
	assert.EqualError(err, nil)

	result = (<-results).Result.(examples.OneEchoResult)
	assert.Equal[examples.OneEchoResult](result, examples.OneEchoResult(cmd))
}
