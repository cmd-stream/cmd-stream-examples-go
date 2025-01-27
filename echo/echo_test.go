package short

import (
	"net"
	"testing"

	"github.com/cmd-stream/base-go"
	ccln "github.com/cmd-stream/cmd-stream-go/client"
	cser "github.com/cmd-stream/cmd-stream-go/server"
	assert_fatal "github.com/ymz-ncnk/assert/fatal"
)

// This example contains an implementation of an echo server.

func TestEcho(t *testing.T) {
	const addr = "127.0.0.1:9000"

	// Start the server.
	server := cser.Default[struct{}](ServerCodec{}, struct{}{})
	l, err := net.Listen("tcp", addr)
	assert_fatal.EqualError(err, nil, t)
	go func() {
		server.Serve(l.(*net.TCPListener))
	}()

	// Create the client.
	conn, err := net.Dial("tcp", addr)
	assert_fatal.EqualError(err, nil, t)
	client, err := ccln.Default(ClientCodec{}, conn)
	assert_fatal.EqualError(err, nil, t)

	// Send a Command and get the Result.
	results := make(chan base.AsyncResult, 1)
	_, err = client.Send(EchoCmd("Hello world"), results)
	assert_fatal.EqualError(err, nil, t)
	assert_fatal.Equal((<-results).Result.(Result), "Hello world", t)

	// Close the client.
	err = client.Close()
	assert_fatal.EqualError(err, nil, t)

	// Close the server.
	err = server.Close()
	assert_fatal.EqualError(err, nil, t)
}
