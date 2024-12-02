package main

import muss "github.com/mus-format/mus-stream-go"

type OkResult bool

func (r OkResult) LastOne() bool {
	return true
}

func (r OkResult) MarshalMUS(w muss.Writer) (n int, err error) {
	return OkResultDTS.Marshal(r, w)
}
