package main

import (
	"net"
	"sync"

	"github.com/cmd-stream/base-go"
	base_server "github.com/cmd-stream/base-go/server"
	examples "github.com/cmd-stream/cmd-stream-examples-go"
	cs_server "github.com/cmd-stream/cmd-stream-go/server"
	delegate_server "github.com/cmd-stream/delegate-go"
	transport_client "github.com/cmd-stream/transport-go/client"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

// Defines a maximum command size in bytes supported by the server.
const MaxCmdSize = 15

// In this example, we set a limit on the maximum command size supported by the
// server.
//
// This is achieved by configuring the server with
// delegate_server.ServerSettings.MaxCmdSize != 0
// and by checking the command size in ServerCodec.Decode() method.
//
// Here we have struct{} as the receiver and examples.EchoCmd as a command.
func main() {
	const addr = "127.0.0.1:9000"

	// Start the server.
	wgS := &sync.WaitGroup{}
	server, err := StartServer(addr, ServerCodec{}, struct{}{}, wgS)
	assert.EqualError(err, nil)

	// Create the client.
	client, err := examples.CreateClient(addr, examples.ClientCodec{})
	assert.EqualError(err, nil)

	// Send too large command.
	var (
		cmd     = examples.EchoCmd("very very very long string")
		results = make(chan base.AsyncResult, 1)
	)
	_, err = client.Send(cmd, results)
	assert.EqualError(err, transport_client.ErrTooLargeCmd)

	// Close the client.
	err = examples.CloseClient(client)
	assert.EqualError(err, nil)

	// Close the server.
	err = examples.CloseServer(server, wgS)
	assert.EqualError(err, nil)
}

func StartServer[T any](addr string, codec cs_server.Codec[T],
	receiver T, wg *sync.WaitGroup) (s *base_server.Server, err error) {
	l, err := net.Listen("tcp", addr)
	assert.EqualError(err, nil)

	settings := delegate_server.ServerSettings{MaxCmdSize: MaxCmdSize} // Defines a limit
	s = cs_server.New[struct{}](cs_server.DefServerInfo, settings,
		cs_server.DefConf, ServerCodec{}, struct{}{}, nil)

	wg.Add(1)
	go func(wg *sync.WaitGroup, listener net.Listener,
		server *base_server.Server) {
		defer wg.Done()
		err := server.Serve(listener.(*net.TCPListener))
		assert.EqualError(err, base_server.ErrClosed)
	}(wg, l, s)
	return
}
