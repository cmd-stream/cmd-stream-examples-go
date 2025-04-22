package tls

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	hw "github.com/cmd-stream/cmd-stream-examples-go/hello-world"

	"github.com/cmd-stream/base-go"
	bcln "github.com/cmd-stream/base-go/client"
	bser "github.com/cmd-stream/base-go/server"
	ccln "github.com/cmd-stream/cmd-stream-go/client"
	cser "github.com/cmd-stream/cmd-stream-go/server"
)

type listenerAdapter struct {
	net.Listener
	l *net.TCPListener
}

func (l listenerAdapter) SetDeadline(tm time.Time) error {
	return l.l.SetDeadline(tm)
}

func StartServer[T any](addr string, codec cser.Codec[T], receiver T,
	wg *sync.WaitGroup) (server *bser.Server, err error) {
	cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
	if err != nil {
		return
	}
	tlsConf := tls.Config{Certificates: []tls.Certificate{cert}}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}
	adapter := listenerAdapter{tls.NewListener(l, &tlsConf), l.(*net.TCPListener)}
	return hw.StartServerWith(addr, codec, cser.NewInvoker(receiver), adapter, wg)
}

func CreateClient[T any](addr string, codec ccln.Codec[T]) (
	c *bcln.Client[T], err error) {
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

	var callback bcln.UnexpectedResultCallback = func(seq base.Seq, result base.Result) {
		fmt.Printf("unexpected result was received: seq %v, result %v\n", seq,
			result)
	}

	return ccln.New(codec, conn,
		ccln.WithBase(
			bcln.WithUnexpectedResultCallback(callback),
		),
	)
}
