package main

import (
	"net"
	"sync"
	"time"

	"github.com/cmd-stream/base-go"
	base_client "github.com/cmd-stream/base-go/client"
	base_server "github.com/cmd-stream/base-go/server"
	examples "github.com/cmd-stream/cmd-stream-examples-go"
	cs_client "github.com/cmd-stream/cmd-stream-go/client"
	cs_server "github.com/cmd-stream/cmd-stream-go/server"
	delegate "github.com/cmd-stream/delegate-go"
	delegate_client "github.com/cmd-stream/delegate-go/client"
	"github.com/cmd-stream/handler-go"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

const ReceiveTimeout = 2 * time.Second

// In this example, the client is trying to keep alive the connection to the
// server, when there are no commands to send.
//
// This is achieved by client configuration - both cs_client.Conf.KeepaliveTime
// and cs_client.Conf.KeepaliveIntvl != 0.
//
// Here we have struct{} as the receiver and examples.EchoCmd as a command.
func main() {
	const addr = "127.0.0.1:9000"

	// Start the server.
	wgS := &sync.WaitGroup{}
	server, err := StartServer(addr, examples.ServerCodec{}, struct{}{}, wgS)
	assert.EqualError(err, nil)

	// Create the client.
	client, err := CreateClient(addr, examples.ClientCodec{})
	assert.EqualError(err, nil)

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
	time.Sleep(ReceiveTimeout + 5*time.Second)

	// Send a command again after a long delay.
	cmd = examples.EchoCmd("hello world again")
	_, err = client.Send(cmd, results)
	assert.EqualError(err, nil)

	result = (<-results).Result.(examples.OneEchoResult)
	assert.Equal[examples.OneEchoResult](result, examples.OneEchoResult(cmd))

	// Close the client.
	err = examples.CloseClient(client)
	assert.EqualError(err, nil)

	// Close the server.
	err = examples.CloseServer(server, wgS)
	assert.EqualError(err, nil)
}

func StartServer[T any](addr string, codec cs_server.Codec[T], receiver T,
	wg *sync.WaitGroup) (s *base_server.Server, err error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}

	// Configuration for the server.
	confServer := cs_server.Conf{
		Base: base_server.Conf{
			WorkersCount: 1,
		},
		Handler: handler.Conf{
			ReceiveTimeout: ReceiveTimeout, // If no commands are sent from the
			// client within 2 seconds, the server will close the connection.
		},
	}
	s = cs_server.New[T](cs_server.DefServerInfo, delegate.ServerSettings{},
		confServer,
		codec,
		receiver,
		nil)

	wg.Add(1)
	go func(wg *sync.WaitGroup, listener net.Listener,
		server *base_server.Server) {
		defer wg.Done()
		err := server.Serve(listener.(*net.TCPListener))
		assert.EqualError(err, base_server.ErrClosed)
	}(wg, l, s)
	return
}

func CreateClient[T any](addr string, codec cs_client.Codec[T]) (
	c *base_client.Client[T], err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}

	// Configuration for the client. If both KeepaliveTime and KeepaliveIntvl != 0,
	// the keeapalive mode will be enabled.
	confClient := cs_client.Conf{
		Delegate: delegate_client.Conf{
			KeepaliveTime:  time.Second,
			KeepaliveIntvl: time.Second,
		},
	}
	// The last nil parameter corresponds to the UnexpectedResultHandler. In this
	// case, unexpected results (if any) received from the server will be simply
	// ignored.
	return cs_client.New[T](cs_server.DefServerInfo, confClient, codec, conn, nil)
}
