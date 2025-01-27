package server_streaming

import "github.com/cmd-stream/transport-go"

// NewResult creates a new Result.
func NewResult(str string, lastOne bool) Result {
	return Result{str, lastOne}
}

// Result implements the Result interface.
type Result struct {
	str     string
	lastOne bool
}

func (r Result) Str() string {
	return r.str
}

// Command in this tutorial sends back several results.
func (r Result) LastOne() bool {
	return r.lastOne
}

func (r Result) Marshal(w transport.Writer) (err error) {
	_, err = ResultDTS.Marshal(r, w)
	return
}
