package main

import (
	"net"
	"sync"

	"github.com/cmd-stream/base-go"
	examples "github.com/cmd-stream/cmd-stream-examples-go"
	cs_client "github.com/cmd-stream/cmd-stream-go/client"
	cs_server "github.com/cmd-stream/cmd-stream-go/server"
	delegate "github.com/cmd-stream/delegate-go"
	transport_client "github.com/cmd-stream/transport-go/client"
	"github.com/ymz-ncnk/assert"
)

func init() {
	assert.On = true
}

const Addr = "127.0.0.1:9000"

// MaxCmdLength defines a maximum command length in bytes supported by the
// server.
const MaxCmdLength = 15

// In this example, we're placing a limit on the maximum size of a command that
// the server can accept. We do this with help of the ServerSettings, which is
// sent to the client, and with the ServerCodec, which checks the size of all
// incoming commands.
//
// Here we have struct{} as the receiver, and examples.EchoCmd as a command.
func main() {
	listener, err := net.Listen("tcp", Addr)
	assert.EqualError(err, nil)

	// Start the server.
	var (
		wgS      = &sync.WaitGroup{}
		settings = delegate.ServerSettings{MaxCmdSize: MaxCmdLength} // The server
		// sends ServerSettings to the client when initializing the connection. The
		// client, in this case, will check the size of each command before sending
		// it.
		server = cs_server.New[struct{}](cs_server.DefServerInfo, settings,
			cs_server.DefConf,
			ServerCodec{}, // Checks the command length in Decode.
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
	client, err := cs_client.NewDef[struct{}](examples.ClientCodec{}, conn, nil)
	assert.EqualError(err, nil)

	// Stop the client.
	defer func() {
		err := client.Close()
		assert.EqualError(err, nil)
		// Wait for the client to stop.
		<-client.Done()
	}()

	// Send too large command.
	var (
		cmd     = examples.EchoCmd("very very very long string")
		results = make(chan base.AsyncResult, 1)
	)
	_, err = client.Send(cmd, results)
	assert.EqualError(err, transport_client.ErrTooBigCmd)
}
