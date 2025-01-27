// Code generated by musgen-go. DO NOT EDIT.

package versioning

import (
	dts "github.com/mus-format/mus-stream-dts-go"
	muss "github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/ord"
)

func MarshalOldSayHelloCmdMUS(v OldSayHelloCmd, w muss.Writer) (n int, err error) {
	return ord.MarshalString(v.str, nil, w)
}

func UnmarshalOldSayHelloCmdMUS(r muss.Reader) (v OldSayHelloCmd, n int, err error) {
	v.str, n, err = ord.UnmarshalString(nil, r)
	return
}

func SizeOldSayHelloCmdMUS(v OldSayHelloCmd) (size int) {
	return ord.SizeString(v.str, nil)
}

func SkipOldSayHelloCmdMUS(r muss.Reader) (n int, err error) {
	return ord.SkipString(nil, r)
}

var OldSayHelloCmdDTS = dts.New[OldSayHelloCmd](OldSayHelloCmdDTM,
	muss.MarshallerFn[OldSayHelloCmd](MarshalOldSayHelloCmdMUS),
	muss.UnmarshallerFn[OldSayHelloCmd](UnmarshalOldSayHelloCmdMUS),
	muss.SizerFn[OldSayHelloCmd](SizeOldSayHelloCmdMUS),
	muss.SkipperFn(SkipOldSayHelloCmdMUS))
