package keepalive

import (
	"sync"
	"testing"
	"time"

	hw "github.com/cmd-stream/cmd-stream-examples-go/hello-world"

	assert_error "github.com/ymz-ncnk/assert/error"
	assert_fatal "github.com/ymz-ncnk/assert/fatal"
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
	wgS := &sync.WaitGroup{}
	server, err := hw.StartServer(addr, hw.ServerCodec{},
		hw.NewGreeter("Hello", "incredible", " "), wgS)
	assert_fatal.EqualError(err, nil, t)

	SendCmd(addr, t)

	// Close the server.
	err = hw.CloseServer(server, wgS)
	assert_fatal.EqualError(err, nil, t)
}

func SendCmd(addr string, t *testing.T) {
	// Create the keepalive client.
	client, err := CreateKeepaliveClient(addr, hw.ClientCodec{})
	assert_fatal.EqualError(err, nil, t)

	var (
		wgR     = &sync.WaitGroup{}
		timeout = time.Second
	)
	// Send SayHelloCmd command.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		sayHelloCmd := hw.NewSayHelloCmd("world")
		wantGreeting := "Hello world"
		err = hw.Exchange[hw.Greeter, hw.Result](sayHelloCmd, timeout, client,
			wantGreeting)
		assert_error.EqualError(err, nil, t)
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
		sayFancyHelloCmd := hw.NewSayFancyHelloCmd("world")
		wantGreeting := "Hello incredible world"
		err = hw.Exchange[hw.Greeter, hw.Result](sayFancyHelloCmd, timeout, client,
			wantGreeting)
		assert_error.EqualError(err, nil, t)
	}()
	wgR.Wait()

	// Close the client.
	err = hw.CloseClient(client)
	assert_fatal.EqualError(err, nil, t)
}
