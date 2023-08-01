package main

import (
	muss "github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/ord"
)

type OkResult bool

func (r OkResult) LastOne() bool {
	return true
}

func MarshalOkResultMUS(result OkResult, w muss.Writer) (n int, err error) {
	return ord.MarshalBool(bool(result), w)
}

func UnmarshalOkResultMUS(r muss.Reader) (result OkResult, n int, err error) {
	b, n, err := ord.UnmarshalBool(r)
	result = OkResult(b)
	return
}

func SizeOkResultMUS(result OkResult) (size int) {
	return ord.SizeBool(bool(result))
}
