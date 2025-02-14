// Code generated by musgen-go. DO NOT EDIT.

package server_streaming

import (
	dts "github.com/mus-format/mus-stream-dts-go"
	muss "github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/ord"
)

func MarshalSayFancyHelloMultiCmdMUS(v SayFancyHelloMultiCmd, w muss.Writer) (n int, err error) {
	return ord.MarshalString(v.str, nil, w)
}

func UnmarshalSayFancyHelloMultiCmdMUS(r muss.Reader) (v SayFancyHelloMultiCmd, n int, err error) {
	v.str, n, err = ord.UnmarshalString(nil, r)
	return
}

func SizeSayFancyHelloMultiCmdMUS(v SayFancyHelloMultiCmd) (size int) {
	return ord.SizeString(v.str, nil)
}

func SkipSayFancyHelloMultiCmdMUS(r muss.Reader) (n int, err error) {
	return ord.SkipString(nil, r)
}

var SayFancyHelloMultiCmdDTS = dts.New[SayFancyHelloMultiCmd](SayFancyHelloMultiCmdDTM,
	muss.MarshallerFn[SayFancyHelloMultiCmd](MarshalSayFancyHelloMultiCmdMUS),
	muss.UnmarshallerFn[SayFancyHelloMultiCmd](UnmarshalSayFancyHelloMultiCmdMUS),
	muss.SizerFn[SayFancyHelloMultiCmd](SizeSayFancyHelloMultiCmdMUS),
	muss.SkipperFn(SkipSayFancyHelloMultiCmdMUS))

func MarshalResultMUS(v Result, w muss.Writer) (n int, err error) {
	n, err = ord.MarshalString(v.str, nil, w)
	if err != nil {
		return
	}
	var n1 int
	n1, err = ord.MarshalBool(v.lastOne, w)
	n += n1
	return
}

func UnmarshalResultMUS(r muss.Reader) (v Result, n int, err error) {
	v.str, n, err = ord.UnmarshalString(nil, r)
	if err != nil {
		return
	}
	var n1 int
	v.lastOne, n1, err = ord.UnmarshalBool(r)
	n += n1
	return
}

func SizeResultMUS(v Result) (size int) {
	size = ord.SizeString(v.str, nil)
	return size + ord.SizeBool(v.lastOne)
}

func SkipResultMUS(r muss.Reader) (n int, err error) {
	n, err = ord.SkipString(nil, r)
	if err != nil {
		return
	}
	var n1 int
	n1, err = ord.SkipBool(r)
	n += n1
	return
}

var ResultDTS = dts.New[Result](ResultDTM,
	muss.MarshallerFn[Result](MarshalResultMUS),
	muss.UnmarshallerFn[Result](UnmarshalResultMUS),
	muss.SizerFn[Result](SizeResultMUS),
	muss.SkipperFn(SkipResultMUS))
