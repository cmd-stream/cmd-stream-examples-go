package main

import (
	"fmt"
	"io"
	"sync"
	"time"

	exmpls "github.com/cmd-stream/cmd-stream-examples-go"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	assert "github.com/ymz-ncnk/assert/panic"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

// Differences from the hello-world example, which uses the MUS serialization
// format:
//  1. Commands (SayHelloCmd and SayFancyHelloCmd) store all properties in data
//     structures that are serializable by Protobuf.
//  2. The protobuf-format.go file instead of mus-format.go.
//
// Everything else, including the codecs, remains the same.
func main() {
	const addr = "127.0.0.1:9000"

	// Start the server.
	wgS := &sync.WaitGroup{}
	server, err := exmpls.StartServer(addr, ServerCodec{},
		exmpls.NewGreeter("Hello", "incredible", " "), wgS)
	assert.EqualError(err, nil)

	SendCmds(addr)
	SendUnsupportedCmd(addr)

	// Close the server.
	err = exmpls.CloseServer(server, wgS)
	assert.EqualError(err, nil)
}

// SendCmds creates a client and sends commands to the server.
func SendCmds(addr string) {
	// Create the client.
	client, err := exmpls.CreateClient(addr, ClientCodec{})
	assert.EqualError(err, nil)

	// Execute commands and wait for them to complete.
	wgC := &sync.WaitGroup{}
	wgC.Add(2)
	timeout := time.Second

	sayHellomd := NewSayHelloCmd("world")
	wantResults := []Result{NewResult("Hello world", true)}
	go exmpls.SendCmd[exmpls.Greeter, Result](sayHellomd, timeout, nil,
		wantResults,
		compareResults,
		client,
		wgC)

	sayFancyHelloCmd := NewSayFancyHelloCmd("world")
	wantResults = []Result{NewResult("Hello incredible world", true)}
	go exmpls.SendCmd[exmpls.Greeter, Result](sayFancyHelloCmd, timeout,
		nil,
		wantResults,
		compareResults,
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

	// Execute the command and wait for it to complete.
	wgC := &sync.WaitGroup{}
	wgC.Add(1)
	timeout := time.Second

	unsupportedCmd := NewUnsupportedCmd("world")
	wantErr := io.EOF
	go exmpls.SendCmd[exmpls.Greeter, Result](unsupportedCmd, timeout,
		wantErr, nil, nil, client, wgC)

	wgC.Wait()

	// There is no need to close the client. If the server receives an unsupported
	// command, it will terminate the connection, which will, in turn, cause the
	// client to close.
}

func compareResults(result, wantResult Result) {
	ignore := cmpopts.IgnoreTypes(protoimpl.MessageState{}, protoimpl.SizeCache(0),
		protoimpl.UnknownFields{})
	if !cmp.Equal(result, wantResult, ignore) {
		panic(fmt.Sprintf("%v != %v", result, wantResult))
	}
}
