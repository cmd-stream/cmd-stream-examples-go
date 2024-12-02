package main

import (
	"io"

	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/mus-stream-dts-go"
	muss "github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/varint"
	"google.golang.org/protobuf/proto"
)

const (
	Eq1DTM com.DTM = iota
	Eq2DTM
	ResultDTM
)

// Eq1CmdMUS

func MarshalEq1CmdProtobuf(c Eq1Cmd, w muss.Writer) (n int, err error) {
	bs, err := proto.Marshal(c)
	if err != nil {
		return
	}
	l := len(bs)
	n, err = varint.MarshalPositiveInt(l, w)
	if err != nil {
		return
	}
	_, err = w.Write(bs)
	n += l
	return
}

func UnmarshalEq1CmdProtobuf(r muss.Reader) (c Eq1Cmd, n int, err error) {
	l, n, err := varint.UnmarshalPositiveInt(r)
	if err != nil {
		return
	}
	bs := make([]byte, l)
	n1, err := io.ReadFull(r, bs)
	n += n1
	if err != nil {
		return
	}
	c.Eq1Data = &Eq1Data{}
	err = proto.Unmarshal(bs, c)
	return
}

func SizeEq1CmdProtobuf(c Eq1Cmd) (size int) {
	panic("not implemented")
}

func SkipEq1CmdProtobuf(r muss.Reader) (n int, err error) {
	panic("not implemented")
}

// Eq2CmdMUS

func MarshalEq2CmdProtobuf(c Eq2Cmd, w muss.Writer) (n int, err error) {
	bs, err := proto.Marshal(c)
	if err != nil {
		return
	}
	l := len(bs)
	n, err = varint.MarshalPositiveInt(l, w)
	if err != nil {
		return
	}
	_, err = w.Write(bs)
	n += l
	return
}

func UnmarshalEq2CmdProtobuf(r muss.Reader) (c Eq2Cmd, n int, err error) {
	l, n, err := varint.UnmarshalPositiveInt(r)
	if err != nil {
		return
	}
	bs := make([]byte, l)
	n1, err := io.ReadFull(r, bs)
	n += n1
	if err != nil {
		return
	}
	c.Eq2Data = &Eq2Data{}
	err = proto.Unmarshal(bs, c)
	return
}

func SizeEq2CmdProtobuf(c Eq2Cmd) (size int) {
	panic("not implemented")
}

func SkipEq2CmdProtobuf(r muss.Reader) (n int, err error) {
	panic("not implemented")
}

// Result

func MarshalResultProtobuf(result Result, w muss.Writer) (n int, err error) {
	bs, err := proto.Marshal(result.ResultData)
	if err != nil {
		return
	}
	l := len(bs)
	n, err = varint.MarshalPositiveInt(l, w)
	if err != nil {
		return
	}
	n1, err := w.Write(bs)
	n += n1
	return
}

func UnmarshalResultProtobuf(r muss.Reader) (result Result, n int, err error) {
	l, n, err := varint.UnmarshalPositiveInt(r)
	if err != nil {
		return
	}
	bs := make([]byte, l)
	n1, err := io.ReadFull(r, bs)
	n += n1
	if err != nil {
		return
	}
	result.ResultData = &ResultData{}
	err = proto.Unmarshal(bs, result)
	return
}

func SizeResultProtobuf(result Result) (size int) {
	panic("not implemented")
}

func SkipResultProtobuf(r muss.Reader) (n int, err error) {
	panic("not implemented")
}

var (
	Eq1DTS = dts.New[Eq1Cmd](Eq1DTM,
		muss.MarshallerFn[Eq1Cmd](MarshalEq1CmdProtobuf),
		muss.UnmarshallerFn[Eq1Cmd](UnmarshalEq1CmdProtobuf),
		muss.SizerFn[Eq1Cmd](SizeEq1CmdProtobuf),
		muss.SkipperFn(SkipEq1CmdProtobuf),
	)
	Eq2DTS = dts.New[Eq2Cmd](Eq2DTM,
		muss.MarshallerFn[Eq2Cmd](MarshalEq2CmdProtobuf),
		muss.UnmarshallerFn[Eq2Cmd](UnmarshalEq2CmdProtobuf),
		muss.SizerFn[Eq2Cmd](SizeEq2CmdProtobuf),
		muss.SkipperFn(SkipEq2CmdProtobuf),
	)
	ResultDTS = dts.New[Result](ResultDTM,
		muss.MarshallerFn[Result](MarshalResultProtobuf),
		muss.UnmarshallerFn[Result](UnmarshalResultProtobuf),
		muss.SizerFn[Result](SizeResultProtobuf),
		muss.SkipperFn(SkipResultProtobuf),
	)
)
