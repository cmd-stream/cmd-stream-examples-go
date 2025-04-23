package hw

import (
	"sync"
	"testing"
	"time"

	cdc "github.com/cmd-stream/codec-mus-stream-go"
	asserterror "github.com/ymz-ncnk/assert/error"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
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
	var (
		receiver = NewGreeter("Hello", "incredible", " ")
		codec    = cdc.NewServerCodec(ResultMUS, CmdMUS)
		wgS      = &sync.WaitGroup{}
	)
	server, err := StartServer(addr, codec, receiver, wgS)
	assertfatal.EqualError(err, nil, t)

	SendCmds(addr, t)

	// Close the server.
	err = CloseServer(server, wgS)
	assertfatal.EqualError(err, nil, t)
}

// SendCmds sends two commands concurrently using a single client.
func SendCmds(addr string, t *testing.T) {
	// Create the client.
	var (
		codec = cdc.NewClientCodec(CmdMUS, ResultMUS)
	)
	client, err := CreateClient[Greeter](addr, codec)
	assertfatal.EqualError(err, nil, t)

	var (
		wgR     = &sync.WaitGroup{}
		timeout = 3 * time.Second
	)
	// Send SayHelloCmd command.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		var (
			cmd                   = NewSayHelloCmd("world")
			wantGreeting Greeting = "Hello world"
		)
		err = Exchange[Greeter, Greeting](cmd, timeout, client, wantGreeting)
		asserterror.EqualError(err, nil, t)
	}()
	// Send SayFancyHelloCmd command.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		var (
			cmd                   = NewSayFancyHelloCmd("world")
			wantGreeting Greeting = "Hello incredible world"
		)
		err = Exchange[Greeter, Greeting](cmd, timeout, client, wantGreeting)
		asserterror.EqualError(err, nil, t)
	}()
	wgR.Wait()

	// Close the client.
	err = CloseClient[Greeter](client)
	assertfatal.EqualError(err, nil, t)
}
