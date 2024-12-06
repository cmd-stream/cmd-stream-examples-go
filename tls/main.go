package main

import (
	"crypto/tls"
	"log"
	"net"
	"sync"
	"time"

	"github.com/cmd-stream/base-go"
	base_client "github.com/cmd-stream/base-go/client"
	base_server "github.com/cmd-stream/base-go/server"
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

// This example shows how you can use cmd-stream-go with the TLS protocol.
//
// Here we have struct{} as the receiver and examples.EchoCmd as a command.
func main() {
	const addr = "127.0.0.1:9000"

	// Start the server.
	wgS := &sync.WaitGroup{}
	server, err := examples.StartServer(addr, examples.ServerCodec{}, struct{}{}, wgS)
	assert.EqualError(err, nil)

	// Create the client.
	client, err := examples.CreateClient(addr, examples.ClientCodec{})
	assert.EqualError(err, nil)

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

	// Close the client.
	err = examples.CloseClient(client)
	assert.EqualError(err, nil)

	// Close the server.
	err = examples.CloseServer(server, wgS)
	assert.EqualError(err, nil)
}

func StartServer[T any](addr string, codec cs_server.Codec[T], receiver T,
	wg *sync.WaitGroup) (s *base_server.Server, err error) {
	cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
	if err != nil {
		return
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}
	tl := tls.NewListener(l, &config)

	s = cs_server.NewDef[T](codec, receiver)

	wg.Add(1)
	go func(wg *sync.WaitGroup, tl net.Listener,
		server *base_server.Server) {
		defer wg.Done()
		err := server.Serve(listenerAdapter{tl, l.(*net.TCPListener)})
		assert.EqualError(err, base_server.ErrClosed)
	}(wg, tl, s)
	return
}

func CreateClient[T any](addr string, codec cs_client.Codec[T]) (
	c *base_client.Client[T], err error) {
	cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert},
		InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", addr, &config)
	if err != nil {
		return
	}
	// The last nil parameter corresponds to the UnexpectedResultHandler. In this
	// case, unexpected results (if any) received from the server will be simply
	// ignored.
	return cs_client.NewDef[T](codec, conn, nil)
}
