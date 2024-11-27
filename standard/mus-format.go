package main

import (
	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/mus-stream-dts-go"
	muss "github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/varint"
)

const (
	Eq1DTM com.DTM = iota
	Eq2DTM
	ResultDTM
)

// Eq1CmdMUS

func MarshalEq1CmdMUS(c Eq1Cmd, w muss.Writer) (n int, err error) {
	n, err = varint.MarshalInt(c.a, w)
	if err != nil {
		return
	}
	var n1 int
	n1, err = varint.MarshalInt(c.b, w)
	n += n1
	if err != nil {
		return
	}
	n1, err = varint.MarshalInt(c.c, w)
	n += n1
	return
}

func UnmarshalEq1CmdMUS(r muss.Reader) (c Eq1Cmd, n int, err error) {
	c.a, n, err = varint.UnmarshalInt(r)
	if err != nil {
		return
	}
	var n1 int
	c.b, n, err = varint.UnmarshalInt(r)
	n += n1
	if err != nil {
		return
	}
	c.c, n, err = varint.UnmarshalInt(r)
	n += n1
	return
}

func SizeEq1CmdMUS(c Eq1Cmd) (size int) {
	size = varint.SizeInt(c.a)
	size += varint.SizeInt(c.b)
	return size + varint.SizeInt(c.c)
}

func SkipEq1CmdMUS(r muss.Reader) (n int, err error) {
	n, err = varint.SkipInt(r)
	if err != nil {
		return
	}
	var n1 int
	n, err = varint.SkipInt(r)
	n += n1
	if err != nil {
		return
	}
	n, err = varint.SkipInt(r)
	n += n1
	return
}

// Eq2CmdMUS

func MarshalEq2CmdMUS(c Eq2Cmd, w muss.Writer) (n int, err error) {
	n, err = varint.MarshalInt(c.a, w)
	if err != nil {
		return
	}
	var n1 int
	n1, err = varint.MarshalInt(c.b, w)
	n += n1
	if err != nil {
		return
	}
	n1, err = varint.MarshalInt(c.c, w)
	n += n1
	return
}

func UnmarshalEq2CmdMUS(r muss.Reader) (c Eq2Cmd, n int, err error) {
	c.a, n, err = varint.UnmarshalInt(r)
	if err != nil {
		return
	}
	var n1 int
	c.b, n, err = varint.UnmarshalInt(r)
	n += n1
	if err != nil {
		return
	}
	c.c, n, err = varint.UnmarshalInt(r)
	n += n1
	return
}

func SizeEq2CmdMUS(c Eq2Cmd) (size int) {
	size = varint.SizeInt(c.a)
	size += varint.SizeInt(c.b)
	return size + varint.SizeInt(c.c)
}

func SkipEq2CmdMUS(r muss.Reader) (n int, err error) {
	n, err = varint.SkipInt(r)
	if err != nil {
		return
	}
	var n1 int
	n, err = varint.SkipInt(r)
	n += n1
	if err != nil {
		return
	}
	n, err = varint.SkipInt(r)
	n += n1
	return
}

// Result

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

func SkipResultMUS(r muss.Reader) (n int, err error) {
	return varint.SkipInt(r)
}

var (
	Eq1DTS = dts.New[Eq1Cmd](Eq1DTM,
		muss.MarshallerFn[Eq1Cmd](MarshalEq1CmdMUS),
		muss.UnmarshallerFn[Eq1Cmd](UnmarshalEq1CmdMUS),
		muss.SizerFn[Eq1Cmd](SizeEq1CmdMUS),
		muss.SkipperFn(SkipEq1CmdMUS),
	)
	Eq2DTS = dts.New[Eq2Cmd](Eq2DTM,
		muss.MarshallerFn[Eq2Cmd](MarshalEq2CmdMUS),
		muss.UnmarshallerFn[Eq2Cmd](UnmarshalEq2CmdMUS),
		muss.SizerFn[Eq2Cmd](SizeEq2CmdMUS),
		muss.SkipperFn(SkipEq2CmdMUS),
	)
	ResultDTS = dts.New[Result](ResultDTM,
		muss.MarshallerFn[Result](MarshalResultMUS),
		muss.UnmarshallerFn[Result](UnmarshalResultMUS),
		muss.SizerFn[Result](SizeResultMUS),
		muss.SkipperFn(SkipResultMUS),
	)
)
