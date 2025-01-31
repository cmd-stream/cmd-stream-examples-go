package versioning

import (
	"sync"
	"testing"
	"time"

	hw "cmd-stream-examples-go/hello-world"

	assert_error "github.com/ymz-ncnk/assert/error"
	assert_fatal "github.com/ymz-ncnk/assert/fatal"
)

// In this example, the server supports two versions of the SayHelloCmd.
// The old version works with an old Receiver and is used by the old client.
//
// After decoding, the server codec migrates the old version to the current one
// so that the server can execute it.

func TestVersionin(t *testing.T) {
	const addr = "127.0.0.1:9008"

	// Start the server.
	wgS := &sync.WaitGroup{}
	receiver := hw.NewGreeter("Hello", "incredible", " ")
	server, err := hw.StartServer(addr, ServerCodec{}, receiver, wgS)
	assert_fatal.EqualError(err, nil, t)

	SendCmd(addr, t)
	SendOldCmd(addr, t)

	// Close the server.
	err = hw.CloseServer(server, wgS)
	assert_fatal.EqualError(err, nil, t)
}

func SendCmd(addr string, t *testing.T) {
	// Create the client.
	client, err := hw.CreateClient(addr, hw.ClientCodec{})
	assert_fatal.EqualError(err, nil, t)

	wgR := &sync.WaitGroup{}
	timeout := time.Second
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
	wgR.Wait()

	// Close the client.
	err = hw.CloseClient(client)
	assert_fatal.EqualError(err, nil, t)
}

func SendOldCmd(addr string, t *testing.T) {
	// Create the client.
	client, err := hw.CreateClient(addr, ClientCodec{})
	assert_fatal.EqualError(err, nil, t)

	wgR := &sync.WaitGroup{}
	timeout := time.Second
	// Send OldSayHelloCmd command.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		oldSayHelloCmd := NewOldSayHelloCmd("world")
		wantGreeting := "Hello world"
		err = hw.Exchange[OldGreeter, hw.Result](oldSayHelloCmd, timeout, client,
			wantGreeting)
		assert_error.EqualError(err, nil, t)
	}()
	wgR.Wait()

	// Close the client.
	err = hw.CloseClient(client)
	assert_fatal.EqualError(err, nil, t)
}
