package tls

import (
	"sync"
	"testing"
	"time"

	hw "github.com/cmd-stream/cmd-stream-examples-go/hello-world"

	cdc "github.com/cmd-stream/codec-mus-stream-go"
	asserterror "github.com/ymz-ncnk/assert/error"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

// cmd-stream-go + TLS protocol.
func TestTLS(t *testing.T) {
	const addr = "127.0.0.1:9007"

	// Start the server.
	var (
		receiver = hw.NewGreeter("Hello", "incredible", " ")
		codec    = cdc.NewServerCodec(hw.ResultMUS, hw.CmdMUS)
		wgS      = &sync.WaitGroup{}
	)
	server, err := StartServer(addr, codec, receiver, wgS)
	assertfatal.EqualError(err, nil, t)

	SendCmd(addr, t)

	// Close the server.
	err = hw.CloseServer(server, wgS)
	assertfatal.EqualError(err, nil, t)
}

func SendCmd(addr string, t *testing.T) {
	// Create the client.
	var (
		codec = cdc.NewClientCodec(hw.CmdMUS, hw.ResultMUS)
	)
	client, err := CreateClient[hw.Greeter](addr, codec)
	assertfatal.EqualError(err, nil, t)

	var (
		wgR     = &sync.WaitGroup{}
		timeout = time.Second
	)
	// Send SayHelloCmd.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		var (
			sayHelloCmd  = hw.NewSayHelloCmd("world")
			wantGreeting = hw.Greeting("Hello world")
		)
		err = hw.Exchange[hw.Greeter, hw.Greeting](sayHelloCmd, timeout, client,
			wantGreeting)
		asserterror.EqualError(err, nil, t)
	}()
	wgR.Wait()

	// Close the client.
	err = hw.CloseClient[hw.Greeter](client)
	assertfatal.EqualError(err, nil, t)
}
