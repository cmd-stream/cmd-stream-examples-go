package main

import (
	"crypto/tls"
	"log"
	"net"
	"sync"
	"time"

	base_client "github.com/cmd-stream/base-go/client"
	base_server "github.com/cmd-stream/base-go/server"
	exmpls "github.com/cmd-stream/cmd-stream-examples-go"
	cs_client "github.com/cmd-stream/cmd-stream-go/client"
	cs_server "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/cmd-stream/delegate-go"
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

// cmd-stream-go + TLS protocol.
func main() {
	const addr = "127.0.0.1:9000"

	// Start the server.
	wgS := &sync.WaitGroup{}
	server, err := StartServer(addr, exmpls.ServerCodec{},
		exmpls.NewGreeter("Hello", "incredible", " "), wgS)
	assert.EqualError(err, nil)

	SendCmd(addr)

	// Close the server.
	err = exmpls.CloseServer(server, wgS)
	assert.EqualError(err, nil)
}

func SendCmd(addr string) {
	// Create the client.
	client, err := CreateClient(addr, exmpls.ClientCodec{})
	assert.EqualError(err, nil)

	wgC := &sync.WaitGroup{}
	wgC.Add(1)
	timeout := time.Second

	// Send a command.
	cmd := exmpls.NewSayHelloCmd("world")
	wantResults := []exmpls.Result{exmpls.NewResult("Hello world", true)}
	exmpls.SendCmd(cmd, timeout, nil, wantResults, exmpls.CompareResults, client, wgC)

	wgC.Wait()

	// Close the client.
	err = exmpls.CloseClient(client)
	assert.EqualError(err, nil)
}

func StartServer[T any](addr string, codec cs_server.Codec[T], receiver T,
	wg *sync.WaitGroup) (s *base_server.Server, err error) {
	cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
	if err != nil {
		return
	}
	tlsConf := tls.Config{Certificates: []tls.Certificate{cert}}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}
	la := listenerAdapter{tls.NewListener(l, &tlsConf), l.(*net.TCPListener)}
	return exmpls.StartServerWith(addr, codec, receiver, la,
		delegate.ServerSettings{},
		wg)
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
	return exmpls.CreateClientWith(cs_client.Conf{}, codec, conn)
}
