package short

import (
	"net"
	"testing"

	"github.com/cmd-stream/base-go"
	ccln "github.com/cmd-stream/cmd-stream-go/client"
	cser "github.com/cmd-stream/cmd-stream-go/server"
	assert_fatal "github.com/ymz-ncnk/assert/fatal"
)

// A minimal example.
func TestEcho(t *testing.T) {
	const addr = "127.0.0.1:9000"

	// Start the server.
	l, err := net.Listen("tcp", addr)
	assert_fatal.EqualError(err, nil, t)
	server := cser.New[struct{}](ServerCodec{}, cser.NewInvoker(struct{}{}))
	go func() {
		server.Serve(l.(*net.TCPListener))
	}()

	// Create the client.
	conn, err := net.Dial("tcp", addr)
	assert_fatal.EqualError(err, nil, t)
	client, err := ccln.New(ClientCodec{}, conn)
	assert_fatal.EqualError(err, nil, t)

	// Send a Command and get the Result.
	var (
		results    = make(chan base.AsyncResult, 1)
		str        = "Hello world"
		cmd        = EchoCmd(str)
		wantResult = Result(str)
	)
	_, err = client.Send(cmd, results)
	assert_fatal.EqualError(err, nil, t)
	assert_fatal.Equal((<-results).Result.(Result), wantResult, t)

	// Close the client.
	err = client.Close()
	assert_fatal.EqualError(err, nil, t)

	// Close the server.
	err = server.Close()
	assert_fatal.EqualError(err, nil, t)
}
