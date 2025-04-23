package streaming

import (
	"sync"
	"testing"
	"time"

	hw "github.com/cmd-stream/cmd-stream-examples-go/hello-world"

	cdc "github.com/cmd-stream/codec-mus-stream-go"
	asserterror "github.com/ymz-ncnk/assert/error"
	assertffatal "github.com/ymz-ncnk/assert/fatal"
)

// In this example, you'll find a Command (SayFancyHelloMultiCmd) that sends
// multiple Results back to the client.
func TestServerStreaming(t *testing.T) {
	const addr = "127.0.0.1:9006"

	// Start the server.
	var (
		receiver = hw.NewGreeter("Hello", "incredible", " ")
		codec    = cdc.NewServerCodec(ResultMUS, CmdMUS)
		wgS      = &sync.WaitGroup{}
	)
	server, err := hw.StartServer(addr, codec, receiver, wgS)
	assertffatal.EqualError(err, nil, t)

	SendMultiCmd(addr, t)

	// Close the server.
	err = hw.CloseServer(server, wgS)
	assertffatal.EqualError(err, nil, t)
}

func SendMultiCmd(addr string, t *testing.T) {
	// Create the client.
	var (
		codec = cdc.NewClientCodec(CmdMUS, ResultMUS)
	)
	client, err := hw.CreateClient(addr, codec)
	assertffatal.EqualError(err, nil, t)

	var (
		wgR     = &sync.WaitGroup{}
		timeout = time.Second
	)
	// Send SayFancyHelloMultiCmd command.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		var (
			sayFancyHelloMultiCmd = NewSayFancyHelloMultiCmd("world")
			wantStrs              = []string{"Hello", "incredible", "world"}
		)
		err = Exchange(sayFancyHelloMultiCmd, timeout, client, wantStrs)
		asserterror.EqualError(err, nil, t)
	}()
	wgR.Wait()

	// Close the client.
	err = hw.CloseClient(client)
	assertffatal.EqualError(err, nil, t)
}
