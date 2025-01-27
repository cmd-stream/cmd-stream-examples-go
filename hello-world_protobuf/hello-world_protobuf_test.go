package hwp

import (
	"sync"
	"testing"
	"time"

	hw "cmd-stream-examples-go/hello-world"

	assert_error "github.com/ymz-ncnk/assert/error"
	assert_fatal "github.com/ymz-ncnk/assert/fatal"
)

// Differences from the hello-world example:
//  1. Commands (SayHelloCmd and SayFancyHelloCmd) store all properties in data
//     structures that are serializable by Protobuf.
//  2. The protobuf-format.go file instead of mus-format.go.
//
// Everything else, including the codecs, remains the same.

func Test(t *testing.T) {
	const addr = "127.0.0.1:9000"

	// Start the server.
	wgS := &sync.WaitGroup{}
	server, err := hw.StartServer(addr, ServerCodec{},
		hw.NewGreeter("Hello", "incredible", " "),
		wgS)
	assert_fatal.EqualError(err, nil, t)

	SendCmds(addr, t)

	// Close the server.
	err = hw.CloseServer(server, wgS)
	assert_fatal.EqualError(err, nil, t)
}

func SendCmds(addr string, t *testing.T) {
	// Create the client.
	client, err := hw.CreateClient[hw.Greeter](addr, ClientCodec{})
	assert_fatal.EqualError(err, nil, t)

	var (
		wgR     = &sync.WaitGroup{}
		timeout = time.Second
	)
	// Send SayHelloCmd command.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		sayHelloCmd := NewSayHelloCmd("world")
		wantGreeting := "Hello world"
		err = hw.Exchange[hw.Greeter, Result](sayHelloCmd, timeout, client, wantGreeting)
		assert_error.EqualError(err, nil, t)
	}()
	// Send SayFancyHelloCmd command.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		sayFancyHelloCmd := NewSayFancyHelloCmd("world")
		wantGreeting := "Hello incredible world"
		err = hw.Exchange[hw.Greeter, Result](sayFancyHelloCmd, timeout, client,
			wantGreeting)
		assert_error.EqualError(err, nil, t)
	}()
	wgR.Wait()

	// Close the client.
	err = hw.CloseClient[hw.Greeter](client)
	assert_fatal.EqualError(err, nil, t)
}
