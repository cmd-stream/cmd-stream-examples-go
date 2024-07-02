package main

import (
	"crypto/tls"
	"log"
	"net"
	"sync"
	"time"

	"github.com/cmd-stream/base-go"
	examples "github.com/cmd-stream/cmd-stream-examples-go"
	cs_client "github.com/cmd-stream/cmd-stream-go/client"
	cs_server "github.com/cmd-stream/cmd-stream-go/server"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

type listenerAdapter struct {
	net.Listener
	l *net.TCPListener
}

func (l listenerAdapter) SetDeadline(tm time.Time) error {
	return l.l.SetDeadline(tm)
}

const Addr = "127.0.0.1:9000"

// This example shows how you can use cmd-stream-go with the TLS protocol.
//
// Here we have struct{} as the receiver and examples.EchoCmd as a command.
func main() {
	cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}

	l, err := net.Listen("tcp", Addr)
	assert.EqualError(err, nil)
	listener := tls.NewListener(l, &config)
	assert.EqualError(err, nil)

	// Start the server.
	var (
		wgS    = &sync.WaitGroup{}
		server = cs_server.NewDef[struct{}](examples.ServerCodec{}, struct{}{})
	)
	wgS.Add(1)
	go func() {
		defer wgS.Done()
		server.Serve(listenerAdapter{listener, l.(*net.TCPListener)})
	}()

	// Stop the server.
	defer func() {
		err := server.Close()
		assert.EqualError(err, nil)
		wgS.Wait()
	}()

	// Connect to the server.
	cert, err = tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	config = tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", Addr, &config)
	assert.EqualError(err, nil)

	// Create the client.
	wgR := &sync.WaitGroup{}
	client, err := cs_client.NewDef[struct{}](examples.ClientCodec{}, conn, nil)
	assert.EqualError(err, nil)

	// Stop the client.
	defer func() {
		wgR.Wait()
		err := client.Close()
		assert.EqualError(err, nil)
		// Wait for the client to stop.
		<-client.Done()
	}()

	// Send a command.
	var (
		cmd     = examples.EchoCmd("hello world")
		results = make(chan base.AsyncResult, 1)
	)
	_, err = client.Send(cmd, results)
	assert.EqualError(err, nil)

	result := (<-results).Result.(examples.OneEchoResult)
	assert.Equal[examples.OneEchoResult](result,
		examples.OneEchoResult(cmd))
}
