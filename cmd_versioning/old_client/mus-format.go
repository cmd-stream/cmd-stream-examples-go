package old_client

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

// PrintCmd

func MarshalPrintCmdV1MUS(c PrintCmd, w muss.Writer) (n int, err error) {
	return ord.MarshalString(c.Text, nil, w)
}

func UnmarshalPrintCmdV1MUS(r muss.Reader) (c PrintCmd, n int, err error) {
	c.Text, n, err = ord.UnmarshalString(nil, r)
	return
}

func SizePrintCmdV1MUS(c PrintCmd) (size int) {
	return ord.SizeString(c.Text, nil)
}

func SkipPtrintCmdV1MUS(r muss.Reader) (n int, err error) {
	return ord.SkipString(nil, r)
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
	PrintCmdDTS = dts.New[PrintCmd](PrintCmdV1DTM,
		muss.MarshallerFn[PrintCmd](MarshalPrintCmdV1MUS),
		muss.UnmarshallerFn[PrintCmd](UnmarshalPrintCmdV1MUS),
		muss.SizerFn[PrintCmd](SizePrintCmdV1MUS),
		muss.SkipperFn(SkipPtrintCmdV1MUS),
	)
	OkResultDTS = dts.New[OkResult](OkResultDTM,
		muss.MarshallerFn[OkResult](MarshalOkResultMUS),
		muss.UnmarshallerFn[OkResult](UnmarshalOkResultMUS),
		muss.SizerFn[OkResult](SizeOkResultMUS),
		muss.SkipperFn(SkipOkResultMUS),
	)
)
