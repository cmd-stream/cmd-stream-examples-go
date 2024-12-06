package examples

import (
	"errors"
	"net"
	"sync"
	"time"

	assert "github.com/ymz-ncnk/assert/panic"

	base_client "github.com/cmd-stream/base-go/client"
	base_server "github.com/cmd-stream/base-go/server"
	cs_client "github.com/cmd-stream/cmd-stream-go/client"
	cs_server "github.com/cmd-stream/cmd-stream-go/server"
)

func StartServer[T any](addr string, codec cs_server.Codec[T], receiver T,
	wg *sync.WaitGroup) (s *base_server.Server, err error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}

	s = cs_server.NewDef[T](codec, receiver)

	wg.Add(1)
	go func(wg *sync.WaitGroup, listener net.Listener,
		server *base_server.Server) {
		defer wg.Done()
		err := server.Serve(listener.(*net.TCPListener))
		assert.EqualError(err, base_server.ErrClosed)
	}(wg, l, s)
	return
}

func CreateClient[T any](addr string, codec cs_client.Codec[T]) (
	c *base_client.Client[T], err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	// The last nil parameter corresponds to the UnexpectedResultHandler. In this
	// case, unexpected results (if any) received from the server will be simply
	// ignored.
	return cs_client.NewDef[T](codec, conn, nil)
}

func CreateReconnectClient[T any](codec cs_client.Codec[T],
	connFactory cs_client.ConnFactory) (c *base_client.Client[T], err error) {
	// assert.EqualError(err, nil)
	// The last nil parameter corresponds to the UnexpectedResultHandler. In this
	// case, unexpected results (if any) received from the server will be simply
	// ignored.
	return cs_client.NewDefReconnect[T](codec,
		connFactory,
		nil)
}

func CloseServer(s *base_server.Server, wg *sync.WaitGroup) (err error) {
	err = s.Close()
	if err != nil {
		return
	}
	wg.Wait()
	return
}

func CloseClient[T any](c *base_client.Client[T]) (err error) {
	err = c.Close()
	if err != nil {
		return
	}
	// The client receives results from the server in the background, so we have
	// to wait for it to stop.
	select {
	case <-time.NewTimer(time.Second).C:
		return errors.New("timeout")
	case <-c.Done():
		return
	}
}
