// hello-world_test.go

package hw

import (
	"sync"
	"testing"
	"time"

	cdc "github.com/cmd-stream/codec-mus-stream-go"
	asserterror "github.com/ymz-ncnk/assert/error"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestHelloWorld(t *testing.T) {
	const addr = "127.0.0.1:9000"

	// Start the server.
	var (
		receiver = NewGreeter("Hello", "incredible", " ")
		// Serializers for base.Cmd and base.Result interfaces allow building
		// a server codec.
		codec = cdc.NewServerCodec(ResultMUS, CmdMUS)
		wgS   = &sync.WaitGroup{}
	)

	server, err := StartServer(addr, codec, receiver, wgS)
	assertfatal.EqualError(err, nil, t)

	SendCmds(addr, t)

	// Close the server.
	err = CloseServer(server, wgS)
	assertfatal.EqualError(err, nil, t)
}

// SendCmds sends two Commands concurrently using a single client.
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
	// Send SayHelloCmd.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		var (
			sayHelloCmd  = NewSayHelloCmd("world")
			wantGreeting = Greeting("Hello world")
		)
		err = Exchange[Greeter, Greeting](sayHelloCmd, timeout, client, wantGreeting)
		asserterror.EqualError(err, nil, t)
	}()
	// Send SayFancyHelloCmd.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		var (
			sayFancyHelloCmd = NewSayFancyHelloCmd("world")
			wantGreeting     = Greeting("Hello incredible world")
		)
		err = Exchange[Greeter, Greeting](sayFancyHelloCmd, timeout, client,
			wantGreeting)
		asserterror.EqualError(err, nil, t)
	}()
	wgR.Wait()

	// Close the client.
	err = CloseClient[Greeter](client)
	assertfatal.EqualError(err, nil, t)
}
