package tls

import (
	"sync"
	"testing"
	"time"

	hw "cmd-stream-examples-go/hello-world"

	assert_error "github.com/ymz-ncnk/assert/error"
	assert_fatal "github.com/ymz-ncnk/assert/fatal"
)

// cmd-stream-go + TLS protocol.

func TestTLS(t *testing.T) {
	const addr = "127.0.0.1:9007"

	// Start the server.
	wgS := &sync.WaitGroup{}
	server, err := StartServer(addr, hw.ServerCodec{},
		hw.NewGreeter("Hello", "incredible", " "), wgS)
	assert_fatal.EqualError(err, nil, t)

	SendCmd(addr, t)

	// Close the server.
	err = hw.CloseServer(server, wgS)
	assert_fatal.EqualError(err, nil, t)
}

func SendCmd(addr string, t *testing.T) {
	// Create the client.
	client, err := CreateClient[hw.Greeter](addr, hw.ClientCodec{})
	assert_fatal.EqualError(err, nil, t)

	var (
		wgR     = &sync.WaitGroup{}
		timeout = time.Second
	)
	// Send SayHelloCmd.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		sayHelloCmd := hw.NewSayHelloCmd("world")
		wantGreeting := "Hello world"
		err = hw.Exchange[hw.Greeter, hw.Result](sayHelloCmd, timeout, client,
			wantGreeting)
		assert_error.EqualError(err, nil, t)
	}()
	wgR.Wait()

	// Close the client.
	err = hw.CloseClient[hw.Greeter](client)
	assert_fatal.EqualError(err, nil, t)
}
