package server_streaming

import (
	"sync"
	"testing"
	"time"

	hw "cmd-stream-examples-go/hello-world"

	assert_error "github.com/ymz-ncnk/assert/error"
	assert_fatal "github.com/ymz-ncnk/assert/fatal"
)

// In this example, you'll find a command (SayFancyHelloMultiCmd) that sends
// multiple results back to the client.

func TestServerStreaming(t *testing.T) {
	const addr = "127.0.0.1:9006"

	// Start the server.
	wgS := &sync.WaitGroup{}
	receiver := hw.NewGreeter("Hello", "incredible", " ")
	server, err := hw.StartServer(addr, ServerCodec{}, receiver, wgS)
	assert_fatal.EqualError(err, nil, t)

	SendMultiCmd(addr, t)

	// Close the server.
	err = hw.CloseServer(server, wgS)
	assert_fatal.EqualError(err, nil, t)
}

func SendMultiCmd(addr string, t *testing.T) {
	// Create the client.
	client, err := hw.CreateClient(addr, ClientCodec{})
	assert_fatal.EqualError(err, nil, t)

	var (
		wgR     = &sync.WaitGroup{}
		timeout = time.Second
	)
	// Send SayFancyHelloMultiCmd command.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		sayFancyHelloMultiCmd := NewSayFancyHelloMultiCmd("world")
		wantStrs := []string{"Hello", "incredible", "world"}
		err = Exchange(sayFancyHelloMultiCmd, timeout, client, wantStrs)
		assert_error.EqualError(err, nil, t)
	}()
	wgR.Wait()

	// Close the client.
	err = hw.CloseClient(client)
	assert_fatal.EqualError(err, nil, t)
}
