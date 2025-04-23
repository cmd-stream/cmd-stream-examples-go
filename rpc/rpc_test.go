package rpc

import (
	"context"
	"sync"
	"testing"

	hw "github.com/cmd-stream/cmd-stream-examples-go/hello-world"

	cdc "github.com/cmd-stream/codec-mus-stream-go"
	asserterror "github.com/ymz-ncnk/assert/error"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

// The rpc example demonstrates how to implement RPC using cmd-stream-go.
//
// Here, you'll find a GreetingService with a SayHello method that sends the
// corresponding SayHelloCmd to the server.
func TestRPC(t *testing.T) {
	const addr = "127.0.0.1:9005"

	// Start the server.
	var (
		receiver = hw.NewGreeter("Hello", "incredible", " ")
		codec    = cdc.NewServerCodec(hw.ResultMUS, hw.CmdMUS)
		wgS      = &sync.WaitGroup{}
	)
	server, err := hw.StartServer(addr, codec, receiver, wgS)
	assertfatal.EqualError(err, nil, t)

	MakeRPC_Call(addr, t)

	// Close the server.
	err = hw.CloseServer(server, wgS)
	assertfatal.Equal(err, nil, t)
}

func MakeRPC_Call(addr string, t *testing.T) {
	// Create the client.
	var (
		codec = cdc.NewClientCodec(hw.CmdMUS, hw.ResultMUS)
	)
	client, err := hw.CreateClient(addr, codec)
	assertfatal.EqualError(err, nil, t)

	// Create the service.
	service := GreeterService{client}

	// Make an RPC call.
	str, err := service.SayHello(context.Background(), "world")
	asserterror.EqualError(err, nil, t)
	asserterror.Equal[string](str, "Hello world", t)

	// Close the client.
	err = hw.CloseClient(client)
	assertfatal.EqualError(err, nil, t)
}
