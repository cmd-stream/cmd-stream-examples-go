package main

import (
	"fmt"
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

// This example demonstrates the standard use of cmd-stream with Protobuf. The
// other files in this package also have useful comments, so check them as well.
//
// In general, this example is the same as the standard one. The difference is
// in the mus-format.go file.
//
// Here we have Calculator as the receiver and Eq1Cmd, Eq2Cmd as commands.
func main() {
	// First of all let's create and run the server.
	listener, err := net.Listen("tcp", Addr)
	assert.EqualError(err, nil)
	// Server will use Calculator to execute received commands.
	server := cs_server.NewDef[Calculator](ServerCodec{}, Calculator{})
	// Run the server.
	wgS := &sync.WaitGroup{}
	wgS.Add(1)
	go runServer(wgS, listener, server)

	// Than connect to the server and create the client.
	conn, err := net.Dial("tcp", Addr)
	assert.EqualError(err, nil)
	// The last nil parameter corresponds to the UnexpectedResultHandler. In this
	// case, unexpected results (if any) received from the server will be simply
	// ignored.
	client, err := cs_client.NewDef[Calculator](ClientCodec{}, conn, nil)
	assert.EqualError(err, nil)

	// And now we will execute two commands.
	wgR := &sync.WaitGroup{}
	wgR.Add(2)
	go sendCmd(wgR, client)
	go sendCmdWithTimeout(wgR, client)

	// Wait while all commands are executed.
	wgR.Wait()

	// Finally let's close the client.
	err = client.Close()
	assert.EqualError(err, nil)
	// The client receives results from the server in the background, so we have
	// to wait for it to stop.
	<-client.Done()

	// And close the server.
	err = server.Close()
	assert.EqualError(err, nil)
	wgS.Wait()
}

func runServer(wg *sync.WaitGroup, listener net.Listener,
	server *base_server.Server) {
	defer wg.Done()
	err := server.Serve(listener.(*net.TCPListener))
	assert.EqualError(err, base_server.ErrClosed)
}

func sendCmd(wg *sync.WaitGroup, client *base_client.Client[Calculator]) {
	defer wg.Done()
	var (
		cmd          = NewEq2Cmd(10, 2, 3)
		want         = NewResult(5)
		asyncResults = make(chan base.AsyncResult, 1)
	)
	_, err := client.Send(cmd, asyncResults)
	assert.EqualError(err, nil)

	asyncResult := <-asyncResults
	// asyncResult.Error != nil if something is wrong with the connection.
	assert.EqualError(asyncResult.Error, nil)
	// The result sent by the command.
	result := asyncResult.Result.(Result)

	if !result.Equal(want) {
		panic(fmt.Sprintf("unexpected result, want %v actual %v", want, result))
	}
}

func sendCmdWithTimeout(wg *sync.WaitGroup,
	client *base_client.Client[Calculator]) {
	defer wg.Done()
	var (
		seq     base.Seq
		cmd     = NewEq1Cmd(1, 2, 3)
		want    = NewResult(6)
		results = make(chan base.AsyncResult, 1)
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
		// The result sent by the command.
		result := asyncResult.Result.(Result)

		if !result.Equal(want) {
			panic(fmt.Sprintf("unexpected result, want %v actual %v", want, result))
		}
	}
}
