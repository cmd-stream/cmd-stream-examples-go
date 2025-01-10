package main

import (
	"sync"
	"time"

	exmpls "github.com/cmd-stream/cmd-stream-examples-go"
	assert "github.com/ymz-ncnk/assert/panic"
)

// In this example, you'll find a command (exmpls.NewSayFancyHelloMultiCmd)
// that sends multiple results back to the client.
func main() {
	const addr = "127.0.0.1:9000"

	// Start the server.
	wgS := &sync.WaitGroup{}
	receiver := exmpls.NewGreeter("Hello", "incredible", " ")
	server, err := exmpls.StartServer(addr, ServerCodec{}, receiver, wgS)
	assert.EqualError(err, nil)

	SendMultiResultCmd(addr)

	// Close the server.
	err = exmpls.CloseServer(server, wgS)
	assert.EqualError(err, nil)
}

func SendMultiResultCmd(addr string) {
	// Create the client.
	client, err := exmpls.CreateClient(addr, exmpls.ClientCodec{})
	assert.EqualError(err, nil)

	// Execute commands and wait for them to complete.
	wgC := &sync.WaitGroup{}
	wgC.Add(1)
	timeout := time.Second

	cmd := exmpls.NewSayFancyHelloMultiCmd("world")
	wantResults := []exmpls.Result{
		exmpls.NewResult("Hello", false),
		exmpls.NewResult("incredible", false),
		exmpls.NewResult("world", true)}
	go exmpls.SendCmd[exmpls.Greeter, exmpls.Result](cmd, timeout, nil,
		wantResults,
		exmpls.CompareResults,
		client,
		wgC)

	wgC.Wait()

	// Close the client.
	err = exmpls.CloseClient(client)
	assert.EqualError(err, nil)
}
