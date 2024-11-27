package main

import (
	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/mus-stream-dts-go"
	muss "github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/ord"
)

// -----------------------------------------------------------------------------
// DTM
// -----------------------------------------------------------------------------

const (
	PrintCmdV1DTM com.DTM = iota
	PrintCmdV2DTM
	OkResultDTM
)

// -----------------------------------------------------------------------------
// Marshal/Unmarshal/Size functions
// -----------------------------------------------------------------------------

// PrintCmdV1

func MarshalPrintCmdV1MUS(c PrintCmdV1, w muss.Writer) (n int, err error) {
	return ord.MarshalString(c.text, nil, w)
}

func UnmarshalPrintCmdV1MUS(r muss.Reader) (c PrintCmdV1, n int, err error) {
	c.text, n, err = ord.UnmarshalString(nil, r)
	return
}

func SizePrintCmdV1MUS(c PrintCmdV1) (size int) {
	return ord.SizeString(c.text, nil)
}

func SkipPtrintCmdV1MUS(r muss.Reader) (n int, err error) {
	return ord.SkipString(nil, r)
}

// PrintCmdV2

func MarshalPrintCmdV2MUS(c PrintCmdV2, w muss.Writer) (n int, err error) {
	n, err = ord.MarshalString(c.from, nil, w)
	if err != nil {
		return
	}
	var n1 int
	n1, err = ord.MarshalString(c.text, nil, w)
	n += n1
	return
}

func UnmarshalPrintCmdV2MUS(r muss.Reader) (c PrintCmdV2, n int, err error) {
	c.from, n, err = ord.UnmarshalString(nil, r)
	if err != nil {
		return
	}
	var n1 int
	c.text, n1, err = ord.UnmarshalString(nil, r)
	n += n1
	return
}

func SizePrintCmdV2MUS(c PrintCmdV2) (size int) {
	size = ord.SizeString(c.from, nil)
	return size + ord.SizeString(c.text, nil)
}

func SkipPrintCmdV2MUS(r muss.Reader) (n int, err error) {
	n, err = ord.SkipString(nil, r)
	if err != nil {
		return
	}
	var n1 int
	n1, err = ord.SkipString(nil, r)
	n += n1
	return
}

// OkResult

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

func SkipOkResultMUS(r muss.Reader) (n int, err error) {
	return ord.SkipBool(r)
}

// -----------------------------------------------------------------------------
// DTS
// -----------------------------------------------------------------------------

var (
	PrintCmdV1DTS = dts.New[PrintCmdV1](PrintCmdV1DTM,
		muss.MarshallerFn[PrintCmdV1](MarshalPrintCmdV1MUS),
		muss.UnmarshallerFn[PrintCmdV1](UnmarshalPrintCmdV1MUS),
		muss.SizerFn[PrintCmdV1](SizePrintCmdV1MUS),
		muss.SkipperFn(SkipPtrintCmdV1MUS),
	)
	PrintCmdV2DTS = dts.New[PrintCmdV2](PrintCmdV2DTM,
		muss.MarshallerFn[PrintCmdV2](MarshalPrintCmdV2MUS),
		muss.UnmarshallerFn[PrintCmdV2](UnmarshalPrintCmdV2MUS),
		muss.SizerFn[PrintCmdV2](SizePrintCmdV2MUS),
		muss.SkipperFn(SkipPrintCmdV2MUS),
	)
	OkResultDTS = dts.New[OkResult](OkResultDTM,
		muss.MarshallerFn[OkResult](MarshalOkResultMUS),
		muss.UnmarshallerFn[OkResult](UnmarshalOkResultMUS),
		muss.SizerFn[OkResult](SizeOkResultMUS),
		muss.SkipperFn(SkipOkResultMUS),
	)
)
