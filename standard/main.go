package main

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/cmd-stream/base-go"
	base_client "github.com/cmd-stream/base-go/client"
	base_server "github.com/cmd-stream/base-go/server"
	cs_client "github.com/cmd-stream/cmd-stream-go/client"
	cs_server "github.com/cmd-stream/cmd-stream-go/server"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

const Addr = "127.0.0.1:9000"

// This example demonstrates the standard use of cmd-stream with MUS. Here we
// have Calculator as the receiver and Eq1Cmd, Eq2Cmd as commands.
//
// The other  files in this package also have useful comments, so check them as
// well.
func main() {
	wgS := &sync.WaitGroup{}
	// First of all let's start the server.
	server, err := startServer(wgS)
	assert.EqualError(err, nil)

	// Than create the client.
	client, err := createClient()
	assert.EqualError(err, nil)

	// Now we will execute two commands.
	wgR := &sync.WaitGroup{}
	wgR.Add(2)
	go sendCmd(wgR, client)
	go sendCmdWithTimeout(wgR, client)
	// And wait while all of them are executed.
	wgR.Wait()

	// Finally let's close the client.
	err = closeClient(client)
	assert.EqualError(err, nil)

	// And close the server.
	err = closeServer(wgS, server)
	assert.EqualError(err, nil)
}

func startServer(wg *sync.WaitGroup) (s *base_server.Server, err error) {
	// First of all let's create and run the server.
	l, err := net.Listen("tcp", Addr)
	assert.EqualError(err, nil)
	// Server will use Calculator to execute received commands.
	s = cs_server.NewDef[Calculator](ServerCodec{}, Calculator{})
	// Run the server.
	// wgS := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup, listener net.Listener,
		server *base_server.Server) {
		defer wg.Done()
		err := server.Serve(listener.(*net.TCPListener))
		assert.EqualError(err, base_server.ErrClosed)
	}(wg, l, s)
	return
}

func createClient() (c *base_client.Client[Calculator], err error) {
	conn, err := net.Dial("tcp", Addr)
	assert.EqualError(err, nil)
	// The last nil parameter corresponds to the UnexpectedResultHandler. In this
	// case, unexpected results (if any) received from the server will be simply
	// ignored.
	return cs_client.NewDef[Calculator](ClientCodec{}, conn, nil)
}

func sendCmd(wg *sync.WaitGroup, client *base_client.Client[Calculator]) {
	defer wg.Done()
	var (
		cmd            = Eq2Cmd{10, 2, 3}
		expectedResult = Result(5)
		asyncResults   = make(chan base.AsyncResult, 1)
	)
	_, err := client.Send(cmd, asyncResults)
	assert.EqualError(err, nil)
	asyncResult := <-asyncResults
	// asyncResult.Error != nil if something is wrong with the connection.
	assert.EqualError(asyncResult.Error, nil)
	result := asyncResult.Result.(Result)
	assert.Equal(result, expectedResult)
}

func sendCmdWithTimeout(wg *sync.WaitGroup,
	client *base_client.Client[Calculator]) {
	defer wg.Done()
	var (
		seq            base.Seq
		cmd            = Eq1Cmd{1, 2, 3}
		expectedResult = Result(6)
		results        = make(chan base.AsyncResult, 1)
	)
	seq, err := client.Send(cmd, results)
	assert.EqualError(err, nil)
	// Let's wait for the result.
	select {
	case <-time.NewTimer(time.Second).C:
		client.Forget(seq) // If we are no longer interested in the results of
		// this command, we should call Forget().
	case asyncResult := <-results:
		// asyncResult.Error != nil if something is wrong with the connection.
		assert.EqualError(asyncResult.Error, nil)
		result := asyncResult.Result.(Result)
		assert.Equal(result, expectedResult)
	}
}

func closeClient(c *base_client.Client[Calculator]) (err error) {
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

func closeServer(wg *sync.WaitGroup, s *base_server.Server) (err error) {
	err = s.Close()
	if err != nil {
		return
	}
	wg.Wait()
	return
}
