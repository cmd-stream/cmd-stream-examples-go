package exmpls

import (
	"github.com/cmd-stream/transport-go"
)

// NewResult creates a new Result.
func NewResult(str string, lastOne bool) Result {
	return Result{str, lastOne}
}

// Result represents the outcome of the command's execution, implements
// base.Result and Marshaller interfaces.
type Result struct {
	str     string
	lastOne bool
}

func (r Result) Str() string {
	return r.str
}

func (r Result) LastOne() bool {
	return r.lastOne
}

func (r Result) Marshal(w transport.Writer) (err error) {
	_, err = ResultDTS.Marshal(r, w)
	return
}
