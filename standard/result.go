package main

import (
	muss "github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/varint"
)

// Result implements the base.Result interface. The client will wait for more
// command results if the LastOne method of the received result returns false.
type Result int

func (r Result) LastOne() bool {
	return true
}

func MarshalResultMUS(result Result, w muss.Writer) (n int, err error) {
	return varint.MarshalInt(int(result), w)
}

func UnmarshalResultMUS(r muss.Reader) (result Result, n int, err error) {
	num, n, err := varint.UnmarshalInt(r)
	result = Result(num)
	return
}

func SizeResultMUS(result Result) (size int) {
	return varint.SizeInt(int(result))
}
