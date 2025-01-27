package reconnect

import (
	"log"
	"net"
	"sync"
	"testing"
	"time"

	hw "cmd-stream-examples-go/hello-world"

	assert_error "github.com/ymz-ncnk/assert/error"
	assert_fatal "github.com/ymz-ncnk/assert/fatal"
)

// In this example, the client attempts to reconnect to the server when the
// connection is lost.
//
// The client is created using the ccln.NewReconnect() function.

const Addr = "127.0.0.1:9000"

type connFactory struct{}

func (f connFactory) New() (conn net.Conn, err error) {
	time.Sleep(100 * time.Millisecond)
	conn, err = net.Dial("tcp", Addr)
	if err != nil {
		log.Println("failed to reconnect")
	} else {
		log.Println("connected")
	}
	return
}

func TestReconnect(t *testing.T) {
	// Start the server.
	wgS := &sync.WaitGroup{}
	server, err := hw.StartServer(Addr, hw.ServerCodec{},
		hw.NewGreeter("Hello", "incredible", " "), wgS)
	assert_fatal.EqualError(err, nil, t)

	// Create the client.
	client, err := CreateReconnectClient(hw.ClientCodec{},
		connFactory{})
	assert_fatal.EqualError(err, nil, t)

	// Close the server.
	err = hw.CloseServer(server, wgS)
	assert_fatal.EqualError(err, nil, t)

	// Start the server again after some time.
	time.Sleep(time.Second)
	server, err = hw.StartServer(Addr, hw.ServerCodec{},
		hw.NewGreeter("Hello", "incredible", " "), wgS)
	assert_fatal.EqualError(err, nil, t)

	// Wait for the client to reconnect.
	time.Sleep(200 * time.Millisecond)

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
	wgR.Wait()

	// Close the client.
	err = hw.CloseClient(client)
	assert_fatal.EqualError(err, nil, t)

	// Close the server.
	err = hw.CloseServer(server, wgS)
	assert_fatal.EqualError(err, nil, t)
}
