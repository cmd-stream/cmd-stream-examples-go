package main

import muss "github.com/mus-format/mus-stream-go"

func NewResult(r int64) Result {
	return Result{ResultData: &ResultData{R: r}}
}

// Result implements the base.Result interface. The client will wait for more
// command results if the LastOne() method of the received result returns false.
type Result struct {
	*ResultData
}

func (r Result) Equal(ar Result) bool {
	return r.R == ar.R
}

func (r Result) LastOne() bool {
	return true
}

func (r Result) MarshalProtobuf(w muss.Writer) (n int, err error) {
	return ResultDTS.Marshal(r, w)
}
