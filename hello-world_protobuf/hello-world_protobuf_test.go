package hwp

import (
	"sync"
	"testing"
	"time"

	hw "github.com/cmd-stream/cmd-stream-examples-go/hello-world"

	cdc "github.com/cmd-stream/codec-mus-stream-go"
	asserterror "github.com/ymz-ncnk/assert/error"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

// Differences from the hello-world example:
//  1. Commands (SayHelloCmd and SayFancyHelloCmd) store all properties in data
//     structures that are serializable by Protobuf.
//  2. The protobuf-format.go file instead of mus-format.gen.go.
//  3. Uses dtm-codec-go to create codecs.
//
// Everything else remains the same.
func Test(t *testing.T) {
	const addr = "127.0.0.1:9002"

	// Start the server.
	var (
		codec    = cdc.NewServerCodec(ResultProtobuf, CmdProtobuf)
		receiver = hw.NewGreeter("Hello", "incredible", " ")
		wgS      = &sync.WaitGroup{}
	)
	server, err := hw.StartServer(addr, codec, receiver, wgS)
	assertfatal.EqualError(err, nil, t)

	SendCmds(addr, t)

	// Close the server.
	err = hw.CloseServer(server, wgS)
	assertfatal.EqualError(err, nil, t)
}

func SendCmds(addr string, t *testing.T) {
	// Create the client.
	var (
		codec = cdc.NewClientCodec(CmdProtobuf, ResultProtobuf)
	)
	client, err := hw.CreateClient[hw.Greeter](addr, codec)
	assertfatal.EqualError(err, nil, t)

	var (
		wgR     = &sync.WaitGroup{}
		timeout = time.Second
	)
	// Send SayHelloCmd command.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		var (
			sayHelloCmd  = NewSayHelloCmd("world")
			wantGreeting = NewGreeting("Hello world")
		)
		err = hw.Exchange[hw.Greeter, Greeting](sayHelloCmd, timeout, client,
			wantGreeting)
		asserterror.EqualError(err, nil, t)
	}()
	// Send SayFancyHelloCmd command.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		var (
			sayFancyHelloCmd = NewSayFancyHelloCmd("world")
			wantGreeting     = NewGreeting("Hello incredible world")
		)
		err = hw.Exchange[hw.Greeter, Greeting](sayFancyHelloCmd, timeout, client,
			wantGreeting)
		asserterror.EqualError(err, nil, t)
	}()
	wgR.Wait()

	// Close the client.
	err = hw.CloseClient[hw.Greeter](client)
	assertfatal.EqualError(err, nil, t)
}
