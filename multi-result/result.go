package main

import (
	muss "github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/ord"
)

type MultiEchoResult struct {
	MultiEchoCmd
	lastOne bool
}

func (r MultiEchoResult) LastOne() bool {
	return r.lastOne
}

func MarshalMultiEchoResultMUS(result MultiEchoResult, w muss.Writer) (n int,
	err error) {
	n, err = MarshalMultiEchoCmdMUS(result.MultiEchoCmd, w)
	if err != nil {
		return
	}
	var n1 int
	n1, err = ord.MarshalBool(result.lastOne, w)
	n += n1
	return
}

func UnmarshalMultiEchoResultMUS(r muss.Reader) (result MultiEchoResult, n int,
	err error) {
	result.MultiEchoCmd, n, err = UnmarshalMultiEchoCmdMUS(r)
	if err != nil {
		return
	}
	var n1 int
	result.lastOne, n1, err = ord.UnmarshalBool(r)
	n += n1
	return
}

func SizeMultiEchoResultMUS(result MultiEchoResult) (size int) {
	size = SizeMultiEchoCmdMUS(result.MultiEchoCmd)
	return size + ord.SizeBool(result.lastOne)
}
