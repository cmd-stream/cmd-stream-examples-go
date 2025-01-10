package main

import (
	"sync"
	"time"

	exmpls "github.com/cmd-stream/cmd-stream-examples-go"
	assert "github.com/ymz-ncnk/assert/panic"
)

// In this example, the server supports two versions of the exmpls.SayHelloCmd.
// The old version works with the old Receiver and is used by the old client.
//
// After decoding, the server codec migrates the old version to the current one
// so that the server can execute it.
func main() {
	const addr = "127.0.0.1:9000"

	// Start the server.
	wgS := &sync.WaitGroup{}
	receiver := exmpls.NewGreeter("Hello", "incredible", " ")
	server, err := exmpls.StartServer(addr, ServerCodec{}, receiver, wgS)
	assert.EqualError(err, nil)

	SendCmd(addr)
	SendOldCmd(addr)

	// Close the server.
	err = exmpls.CloseServer(server, wgS)
	assert.EqualError(err, nil)
}

func SendCmd(addr string) {
	// Create the client.
	client, err := exmpls.CreateClient(addr, exmpls.ClientCodec{})
	assert.EqualError(err, nil)

	// Execute command and wait for it to complete.
	wgC := &sync.WaitGroup{}
	wgC.Add(1)
	timeout := time.Second

	sayHelloCmd := exmpls.NewSayHelloCmd("world")
	wantResults := []exmpls.Result{exmpls.NewResult("Hello world", true)}
	go exmpls.SendCmd[exmpls.Greeter, exmpls.Result](sayHelloCmd, timeout,
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

func SendOldCmd(addr string) {
	// Create an old client.
	oldClient, err := exmpls.CreateClient(addr, ClientCodec{})
	assert.EqualError(err, nil)

	// Execute command and wait for it to complete.
	wgC := &sync.WaitGroup{}
	wgC.Add(1)
	timeout := time.Second

	oldSayHelloCmd := exmpls.NewOldSayHelloCmd("world")
	wantResults := []exmpls.Result{exmpls.NewResult("Hello world", true)}
	go exmpls.SendCmd[exmpls.OldGreeter, exmpls.Result](oldSayHelloCmd, timeout,
		nil,
		wantResults,
		exmpls.CompareResults,
		oldClient,
		wgC)

	wgC.Wait()

	// Close the client.
	err = exmpls.CloseClient(oldClient)
	assert.EqualError(err, nil)
}
