package main

import (
	"net"
	"sync"

	"github.com/cmd-stream/base-go"
	base_server "github.com/cmd-stream/base-go/server"
	examples "github.com/cmd-stream/cmd-stream-examples-go"
	cs_client "github.com/cmd-stream/cmd-stream-go/client"
	cs_server "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/cmd-stream/delegate-go"
	transport_client "github.com/cmd-stream/transport-go/client"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

const Addr = "127.0.0.1:9000"

// MaxCmdLength defines a maximum command length in bytes supported by the
// server.
const MaxCmdLength = 15

// In this example, we set a limit on the maximum command size supported by the
// server.
//
// To do this, we use ServerSettings and ServerCodec. The client receives
// ServerSettings from the server when the connection is initialized, and if
// ServerSettings.MaxCmdSize != 0, it will check the command size (using the
// ClientCodec.Size method) before sending it. ServerCodec, in turn, works on
// the server side and checks the size of each incoming command.
//
// Here we have struct{} as the receiver and examples.EchoCmd as a command.
func main() {
	// First of all let's create and run the server.
	listener, err := net.Listen("tcp", Addr)
	assert.EqualError(err, nil)
	settings := delegate.ServerSettings{MaxCmdSize: MaxCmdLength} // Sets a limit
	// Server will use Calculator to execute received commands.
	server := cs_server.New[struct{}](cs_server.DefServerInfo, settings,
		cs_server.DefConf, ServerCodec{}, struct{}{}, nil)
	// Run the server.
	wgS := &sync.WaitGroup{}
	wgS.Add(1)
	go runServer(wgS, listener, server)

	// Than connect to the server and create the client.
	conn, err := net.Dial("tcp", Addr)
	assert.EqualError(err, nil)
	// The last nil parameter corresponds to the UnexpectedResultHandler. In this
	// case, unexpected results (if any) received from the server will be simply
	// ignored.
	client, err := cs_client.NewDef[struct{}](examples.ClientCodec{}, conn, nil)
	assert.EqualError(err, nil)

	// Send too large command.
	var (
		cmd     = examples.EchoCmd("very very very long string")
		results = make(chan base.AsyncResult, 1)
	)
	_, err = client.Send(cmd, results)
	assert.EqualError(err, transport_client.ErrTooLargeCmd)

	// Finally let's close the client.
	err = client.Close()
	assert.EqualError(err, nil)
	// The client receives results from the server in the background, so we have
	// to wait for it to stop.
	<-client.Done()

	// And close the server.
	err = server.Close()
	assert.EqualError(err, nil)
	wgS.Wait()
}

func runServer(wg *sync.WaitGroup, listener net.Listener,
	server *base_server.Server) {
	defer wg.Done()
	err := server.Serve(listener.(*net.TCPListener))
	assert.EqualError(err, base_server.ErrClosed)
}
