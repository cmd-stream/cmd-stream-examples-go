package hw

import (
	"sync"
	"testing"
	"time"

	assert_error "github.com/ymz-ncnk/assert/error"
	assert_fatal "github.com/ymz-ncnk/assert/fatal"
)

// This example demonstrates how to use cmd-stream-go.
//
// 1. Implement the Command Pattern:
//   - Define the Receiver, Commands, Results, and Invoker.
//   - Commands must implement base.Cmd interface.
//   - Results must implement base.Result interface.
//   - Both must implement the Marshaller interface (required for codecs).
//
// 2. Define Codecs:
//   - Client codec must implement ccln.Codec interface.
//   - Server codec must implement cser.Codec interface.
//
// 3. Create server and client:
//   - Instantiate with the appropriate codecs.
//
// Note: This example uses musgen-go to generate MUS serialization code.
func TestGreeting(t *testing.T) {
	const addr = "127.0.0.1:9001"

	// Start the server.
	wgS := &sync.WaitGroup{}
	receiver := NewGreeter("Hello", "incredible", " ")
	server, err := StartServer(addr, ServerCodec{}, receiver, wgS)
	assert_fatal.EqualError(err, nil, t)

	SendCmds(addr, t)

	// Close the server.
	err = CloseServer(server, wgS)
	assert_fatal.EqualError(err, nil, t)
}

// SendCmds sends two commands concurrently using a single client.
func SendCmds(addr string, t *testing.T) {
	// Create the client.
	client, err := CreateClient[Greeter](addr, ClientCodec{})
	assert_fatal.EqualError(err, nil, t)

	var (
		wgR     = &sync.WaitGroup{}
		timeout = 3 * time.Second
	)
	// Send SayHelloCmd command.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		cmd := NewSayHelloCmd("world")
		wantGreeting := "Hello world"
		err = Exchange[Greeter, Result](cmd, timeout, client, wantGreeting)
		assert_error.EqualError(err, nil, t)
	}()
	// Send SayFancyHelloCmd command.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		cmd := NewSayFancyHelloCmd("world")
		wantGreeting := "Hello incredible world"
		err = Exchange[Greeter, Result](cmd, timeout, client, wantGreeting)
		assert_error.EqualError(err, nil, t)
	}()
	wgR.Wait()

	// Close the client.
	err = CloseClient[Greeter](client)
	assert_fatal.EqualError(err, nil, t)
}
