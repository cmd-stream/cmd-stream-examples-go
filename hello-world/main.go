package main

import (
	"io"
	"sync"
	"time"

	exmpls "github.com/cmd-stream/cmd-stream-examples-go"
	assert "github.com/ymz-ncnk/assert/panic"
)

// To start using cmd-stream-go you have to:
//   - Implement the Command pattern: define the Invoker, Receiver, commands and
//     results.
//   - Create the server side and the client side codecs.
//   - Configure and start the server.
//   - Configure the client.
//
// Commands have to implement the base.Cmd interface, it is required by the
// server. Results - the base.Result interface, it is required by the client to
// determine if the received result is the final one for the command.
//
// Besides that the codec implementations (exmpls.ServerCodec and
// exmpls.ClientCodec), require commands and Results to implement the
// exempls.Marshaller interface, it is used for encoding.
//
// Serialization is performed here using the mus-stream-go serializer. Each
// command data is prefixed with a DTM (thanks to the mus-stream-dts-go library):
//
//	DTM + cmd bytes
//
// This prefix represents the data type and is used for decoding.
//
// The server uses:
//   - The default Invoker.
//   - The exmpls.Greeter type as a Receiver.
//   - exmpls.SayHelloCmd and exmpls.SayFancyHelloCmd as supported commands.
//
// The client in addition to the supported commands, also attempts to send an
// unsupported one (the server codec doesn't know nothing about it).
func main() {
	const addr = "127.0.0.1:9000"

	// Start the server.
	wgS := &sync.WaitGroup{}
	server, err := exmpls.StartServer(addr, exmpls.ServerCodec{},
		exmpls.NewGreeter("Hello", "incredible", " "), wgS)
	assert.EqualError(err, nil)

	SendCmds(addr)
	SendUnsupportedCmd(addr)

	// Close the server.
	err = exmpls.CloseServer(server, wgS)
	assert.EqualError(err, nil)
}

// SendCmds creates a client and sends commands to the server each in own
// goroutine.
func SendCmds(addr string) {
	// Create the client.
	client, err := exmpls.CreateClient(addr, exmpls.ClientCodec{})
	assert.EqualError(err, nil)

	// Execute commands and wait for them to complete.
	wgC := &sync.WaitGroup{}
	wgC.Add(2)
	timeout := time.Second

	sayHellomd := exmpls.NewSayHelloCmd("world")
	wantResults := []exmpls.Result{exmpls.NewResult("Hello world", true)}
	go exmpls.SendCmd[exmpls.Greeter, exmpls.Result](sayHellomd, timeout, nil,
		wantResults,
		exmpls.CompareResults,
		client,
		wgC)

	sayFancyHelloCmd := exmpls.NewSayFancyHelloCmd("world")
	wantResults = []exmpls.Result{exmpls.NewResult("Hello incredible world", true)}
	go exmpls.SendCmd[exmpls.Greeter, exmpls.Result](sayFancyHelloCmd, timeout,
		nil,
		wantResults,
		exmpls.CompareResults,
		client,
		wgC)

	wgC.Wait()

	// Close the client.
	err = exmpls.CloseClient(client)
	assert.EqualError(err, nil)
}

// SendUnsupportedCmd creates the client and sends an unsupported command to the
// server.
func SendUnsupportedCmd(addr string) {
	// Create the client.
	client, err := exmpls.CreateClient(addr, exmpls.ClientCodec{})
	assert.EqualError(err, nil)

	wgC := &sync.WaitGroup{}
	wgC.Add(1)
	timeout := time.Second

	cmd := exmpls.NewUnsupportedCmd("world")
	wantErr := io.EOF
	go exmpls.SendCmd[exmpls.Greeter, exmpls.Result](cmd, timeout, wantErr, nil,
		nil,
		client,
		wgC)

	wgC.Wait()

	// There is no need to close the client. If the server receives an unsupported
	// command, it will terminate the connection, which will, in turn, cause the
	// client to close.
}
