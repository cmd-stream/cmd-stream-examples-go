package keepalive

import (
	"sync"
	"testing"
	"time"

	hw "github.com/cmd-stream/cmd-stream-examples-go/hello-world"

	cdc "github.com/cmd-stream/codec-mus-stream-go"
	asserterror "github.com/ymz-ncnk/assert/error"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

// In this example, the client attempts to keep the connection to the server
// alive even when there are no commands to send.
//
// The client is initialized with the following configuration:
//
// ccln.New(codec, conn,
//
//	 ...
//		ccln.WithKeepalive(
//			dcln.WithKeepaliveTime(...),
//			dcln.WithKeepaliveIntvl(...),
//		),
//
// )
func TestKeepalive(t *testing.T) {
	const addr = "127.0.0.1:9003"

	// Start the server.
	var (
		receiver = hw.NewGreeter("Hello", "incredible", " ")
		codec    = cdc.NewServerCodec(hw.ResultMUS, hw.CmdMUS)
		wgS      = &sync.WaitGroup{}
	)
	server, err := hw.StartServer(addr, codec, receiver, wgS)
	assertfatal.EqualError(err, nil, t)

	SendCmd(addr, t)

	// Close the server.
	err = hw.CloseServer(server, wgS)
	assertfatal.EqualError(err, nil, t)
}

func SendCmd(addr string, t *testing.T) {
	// Create the keepalive client.
	var (
		codec = cdc.NewClientCodec(hw.CmdMUS, hw.ResultMUS)
	)
	client, err := CreateKeepaliveClient(addr, codec)
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
			sayHelloCmd  = hw.NewSayHelloCmd("world")
			wantGreeting = hw.Greeting("Hello world")
		)
		err = hw.Exchange[hw.Greeter, hw.Greeting](sayHelloCmd, timeout, client,
			wantGreeting)
		asserterror.EqualError(err, nil, t)
	}()

	// Ping-Pong time... When there are no commands to send, the client sends
	// a predefined PingCmd.
	time.Sleep(2 * hw.CmdReceiveDuration)

	// Send a command again after the long delay to check if the connection is
	// still alive.
	// Send SayFancyHelloCmd command.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		var (
			sayFancyHelloCmd = hw.NewSayFancyHelloCmd("world")
			wantGreeting     = hw.Greeting("Hello incredible world")
		)
		err = hw.Exchange[hw.Greeter, hw.Greeting](sayFancyHelloCmd, timeout, client,
			wantGreeting)
		asserterror.EqualError(err, nil, t)
	}()
	wgR.Wait()

	// Close the client.
	err = hw.CloseClient(client)
	assertfatal.EqualError(err, nil, t)
}
