package hw

import "github.com/cmd-stream/transport-go"

// NewResult creates a new Result.
func NewResult(str string) Result {
	return Result{str}
}

// Result implements the base.Result interface.
type Result struct {
	str string
}

func (r Result) Greeting() string {
	return r.str
}

// All Commands in this tuttorial send back a single Result.
func (r Result) LastOne() bool {
	return true
}

func (r Result) Marshal(w transport.Writer) (err error) {
	_, err = ResultDTS.Marshal(r, w)
	return
}
