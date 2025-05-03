package group

import (
	"sync"
	"testing"
	"time"

	hw "github.com/cmd-stream/cmd-stream-examples-go/hello-world"
	cdc "github.com/cmd-stream/codec-mus-stream-go"
	asserterror "github.com/ymz-ncnk/assert/error"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

// This example demonstrates the use of a client group to establish a
// high-performance communication channel with the server.
//
// Differences from the hello-world example: Commands are send using the client
// group of 2 clients.
//
// Everything else remains the same.
func TestClientGroup(t *testing.T) {
	const addr = "127.0.0.1:9000"

	// Start the server.
	var (
		receiver = hw.NewGreeter("Hello", "incredible", " ")
		// Serializers for base.Cmd and base.Result interfaces allow building
		// a server codec.
		codec = cdc.NewServerCodec(hw.ResultMUS, hw.CmdMUS)
		wgS   = &sync.WaitGroup{}
	)

	server, err := hw.StartServer(addr, codec, receiver, wgS)
	assertfatal.EqualError(err, nil, t)

	SendCmds(addr, t)

	// Close the server.
	err = hw.CloseServer(server, wgS)
	assertfatal.EqualError(err, nil, t)
}

func SendCmds(addr string, t *testing.T) {
	var (
		codec   = cdc.NewClientCodec(hw.CmdMUS, hw.ResultMUS)
		grp     = CreateClientGroup(addr, 2, codec)
		wgR     = &sync.WaitGroup{}
		timeout = 3 * time.Second
	)
	// Send SayHelloCmd.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		var (
			sayHelloCmd  = hw.NewSayHelloCmd("world")
			wantGreeting = hw.Greeting("Hello world")
		)
		err := Exchange(sayHelloCmd, timeout, grp, wantGreeting)
		asserterror.EqualError(err, nil, t)
	}()
	// Send SayFancyHelloCmd.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		var (
			sayFancyHelloCmd = hw.NewSayFancyHelloCmd("world")
			wantGreeting     = hw.Greeting("Hello incredible world")
		)
		err := Exchange(sayFancyHelloCmd, timeout, grp, wantGreeting)
		asserterror.EqualError(err, nil, t)
	}()
	wgR.Wait()

	// Close the client group.
	err := CloseGroup(grp)
	assertfatal.EqualError(err, nil, t)
}
