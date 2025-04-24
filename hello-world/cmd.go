// cmd.go

package hw

import (
	"context"
	"errors"
	"time"

	"github.com/cmd-stream/base-go"
	muss "github.com/mus-format/mus-stream-go"
)

// CmdExecDuration defines the duration of Command execution.
const CmdExecDuration = time.Second

// NewSayHelloCmd creates a new SayHelloCmd.
func NewSayHelloCmd(str string) SayHelloCmd {
	return SayHelloCmd{str}
}

// SayHelloCmd implements the base.Cmd[Greeter] interface and produces
// greetings like "Hello world".
type SayHelloCmd struct {
	str string
}

func (c SayHelloCmd) Exec(ctx context.Context, at time.Time, seq base.Seq,
	receiver Greeter, proxy base.Proxy) error {
	var (
		result = Greeting(
			receiver.Join(receiver.Interjection(), c.str),
		)
		// Limiting the execution time of a Command on the server is
		// considered a good practice that can be achieved with a deadline.
		deadline = at.Add(CmdExecDuration)
	)

	// A Command can behave in various ways:
	// 1. It can send back only one Result:
	//
	//      return proxy.SendWithDeadline(deadline, seq, result)
	//
	//    Note: The deadline is applied to the entire connection. This means
	//    it will also affect subsequent commands unless they update it with
	//    their own value.
	//
	//    So if one Command uses the Proxy.SendWithDeadline() method, all
	//    others should do the same. Mixing Proxy.Send() and
	//    Proxy.SendWithDeadline() can result in unpredictable behavior
	//    due to unintended deadline propagation.
	//
	//    To cancel the deadline, use time.Time{}:
	//
	//      return proxy.SendWithDeadline(time.Time{}, seq, result)
	//
	// 2. It can perform context-related tasks:
	//
	//      ownCtx, cancel := context.WithDeadline(ctx, deadline)
	//      // Use ownCtx to perform a context-related task.
	//      ...
	//      return proxy.SendWithDeadline(deadline, seq, result)
	//
	// 3. It can send back multiple results (server streaming):
	//
	//      err = proxy.SendWithDeadline(deadline, seq, result1)
	//      if err != nil {
	//         return
	//      }
	//      return proxy.Send(seq, result2)
	//
	//    All results except the first one are sent back using the
	//    Proxy.Send() method.
	//
	// Regardless of the case, the final Result should have
	// Result.LastOne() == true.

	// As you can see, the current Command sends back only one Result.
	return proxy.SendWithDeadline(deadline, seq, result)
}

// NewSayFancyHelloCmd creates a new SayFancyHelloCmd.
func NewSayFancyHelloCmd(str string) SayFancyHelloCmd {
	return SayFancyHelloCmd{str}
}

// SayFancyHelloCmd implements the base.Cmd[Greeter] interface and produces
// greetings like "Hello incredible world".
type SayFancyHelloCmd struct {
	str string
}

func (c SayFancyHelloCmd) Exec(ctx context.Context, at time.Time,
	seq base.Seq, receiver Greeter, proxy base.Proxy) error {
	// SayFancyHelloCmd differs from SayHelloCmd in the way it uses the
	// Receiver.
	var (
		result = Greeting(
			receiver.Join(receiver.Interjection(), receiver.Adjective(), c.str),
		)
		deadline = at.Add(CmdExecDuration)
	)
	return proxy.SendWithDeadline(deadline, seq, result)
}

// ValidateLength is used on the server to validate a Command's content.
func ValidateLength(length int) (err error) {
	if length > 10 {
		return errors.New("too large")
	}
	return
}

func (c SayHelloCmd) MarshalTypedMUS(w muss.Writer) (n int, err error) {
	return SayHelloCmdDTS.Marshal(c, w) // The Command will be marshalled as
	// 'DTM + command data'.
}

func (c SayHelloCmd) SizeTypedMUS() (size int) {
	return SayHelloCmdDTS.Size(c)
}

func (c SayFancyHelloCmd) MarshalTypedMUS(w muss.Writer) (n int, err error) {
	return SayFancyHelloCmdDTS.Marshal(c, w)
}

func (c SayFancyHelloCmd) SizeTypedMUS() (size int) {
	return SayFancyHelloCmdDTS.Size(c)
}
