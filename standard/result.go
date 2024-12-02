package main

import muss "github.com/mus-format/mus-stream-go"

// Result implements the base.Result interface. The client will wait for more
// command results if the LastOne() method of the received result returns false.
type Result int

func (r Result) LastOne() bool {
	return true
}

func (r Result) MarshalMUS(w muss.Writer) (n int, err error) {
	return ResultDTS.Marshal(r, w)
}
