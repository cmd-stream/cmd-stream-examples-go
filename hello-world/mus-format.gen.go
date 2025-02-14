// Code generated by musgen-go. DO NOT EDIT.

package hw

import (
	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/mus-stream-dts-go"
	muss "github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/ord"
)

func MarshalSayHelloCmdMUS(v SayHelloCmd, w muss.Writer) (n int, err error) {
	return ord.MarshalString(v.str, nil, w)
}

func UnmarshalSayHelloCmdMUS(r muss.Reader) (v SayHelloCmd, n int, err error) {
	v.str, n, err = ord.UnmarshalValidString(nil,
		com.ValidatorFn[int](ValidateLength),
		false,
		r)
	return
}

func SizeSayHelloCmdMUS(v SayHelloCmd) (size int) {
	return ord.SizeString(v.str, nil)
}

func SkipSayHelloCmdMUS(r muss.Reader) (n int, err error) {
	return ord.SkipString(nil, r)
}

var SayHelloCmdDTS = dts.New[SayHelloCmd](SayHelloCmdDTM,
	muss.MarshallerFn[SayHelloCmd](MarshalSayHelloCmdMUS),
	muss.UnmarshallerFn[SayHelloCmd](UnmarshalSayHelloCmdMUS),
	muss.SizerFn[SayHelloCmd](SizeSayHelloCmdMUS),
	muss.SkipperFn(SkipSayHelloCmdMUS))

func MarshalSayFancyHelloCmdMUS(v SayFancyHelloCmd, w muss.Writer) (n int, err error) {
	return ord.MarshalString(v.str, nil, w)
}

func UnmarshalSayFancyHelloCmdMUS(r muss.Reader) (v SayFancyHelloCmd, n int, err error) {
	v.str, n, err = ord.UnmarshalValidString(nil,
		com.ValidatorFn[int](ValidateLength),
		false,
		r)
	return
}

func SizeSayFancyHelloCmdMUS(v SayFancyHelloCmd) (size int) {
	return ord.SizeString(v.str, nil)
}

func SkipSayFancyHelloCmdMUS(r muss.Reader) (n int, err error) {
	return ord.SkipString(nil, r)
}

var SayFancyHelloCmdDTS = dts.New[SayFancyHelloCmd](SayFancyHelloCmdDTM,
	muss.MarshallerFn[SayFancyHelloCmd](MarshalSayFancyHelloCmdMUS),
	muss.UnmarshallerFn[SayFancyHelloCmd](UnmarshalSayFancyHelloCmdMUS),
	muss.SizerFn[SayFancyHelloCmd](SizeSayFancyHelloCmdMUS),
	muss.SkipperFn(SkipSayFancyHelloCmdMUS))

func MarshalResultMUS(v Result, w muss.Writer) (n int, err error) {
	return ord.MarshalString(v.str, nil, w)
}

func UnmarshalResultMUS(r muss.Reader) (v Result, n int, err error) {
	v.str, n, err = ord.UnmarshalString(nil, r)
	return
}

func SizeResultMUS(v Result) (size int) {
	return ord.SizeString(v.str, nil)
}

func SkipResultMUS(r muss.Reader) (n int, err error) {
	return ord.SkipString(nil, r)
}

var ResultDTS = dts.New[Result](ResultDTM,
	muss.MarshallerFn[Result](MarshalResultMUS),
	muss.UnmarshallerFn[Result](UnmarshalResultMUS),
	muss.SizerFn[Result](SizeResultMUS),
	muss.SkipperFn(SkipResultMUS))
