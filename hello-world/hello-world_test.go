package hw

import (
	"sync"
	"testing"
	"time"

	assert_error "github.com/ymz-ncnk/assert/error"
	assert_fatal "github.com/ymz-ncnk/assert/fatal"
)

// To start using cmd-stream-go you have to:
//   1. Implement the Command Pattern: define the Receiver, Commands, Results,
//      and Invoker.
//   2. Create server and client codecs.
//   3. Configure the server.
//   4. Configure the client.
//
// Commands must implement the base.Cmd interface, Results - base.Result.
// Besides that the codec implementations require them both to implement the
// Marshaller interface, that is used for encoding.
//
// The Invoker must implement the handler.Invoker interface.
//
// Client and server codecs - the ccln.Codec and cser.Codec interfaces,
// respectively.
//
// In this example, musgen-go is used to generate serialization code. Refer
// to the gen/main.go file for details.

func TestGreeting(t *testing.T) {
	const addr = "127.0.0.1:9000"

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
		timeout = time.Second
	)
	// Send SayHelloCmd command.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		sayHelloCmd := NewSayHelloCmd("world")
		wantGreeting := "Hello world"
		err = Exchange[Greeter, Result](sayHelloCmd, timeout, client, wantGreeting)
		assert_error.EqualError(err, nil, t)
	}()
	// Send SayFancyHelloCmd command.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		sayFancyHelloCmd := NewSayFancyHelloCmd("world")
		wantGreeting := "Hello incredible world"
		err = Exchange[Greeter, Result](sayFancyHelloCmd, timeout, client,
			wantGreeting)
		assert_error.EqualError(err, nil, t)
	}()
	wgR.Wait()

	// Close the client.
	err = CloseClient[Greeter](client)
	assert_fatal.EqualError(err, nil, t)
}
