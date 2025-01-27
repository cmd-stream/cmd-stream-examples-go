package hwp

import (
	"io"

	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/mus-stream-dts-go"
	muss "github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/varint"
	"google.golang.org/protobuf/proto"

	hw "cmd-stream-examples-go/hello-world"
)

// SayHelloCmd

func MarshalSayHelloCmdProtobuf(c SayHelloCmd, w muss.Writer) (n int,
	err error) {
	return marshalCmd(c, w)
}

func UnmarshalSayHelloCmdProtobuf(r muss.Reader) (c SayHelloCmd, n int, err error) {
	data := &SayHelloData{}
	n, err = unmarshalCmd[*SayHelloData](data,
		com.ValidatorFn[int](hw.ValidateLength), r)
	if err != nil {
		return
	}
	c.SayHelloData = data
	return
}

func SizeSayHelloCmdProtobuf(c SayHelloCmd) (size int) {
	panic("not implemented")
}

func SkipSayHelloCmdProtobuf(r muss.Reader) (n int, err error) {
	panic("not implemented")
}

// SayFancyHelloCmd

func MarshalSayFancyHelloCmdProtobuf(c SayFancyHelloCmd, w muss.Writer) (
	n int, err error) {
	return marshalCmd(c, w)
}

func UnmarshalSayFancyHelloCmdProtobuf(r muss.Reader) (c SayFancyHelloCmd,
	n int, err error) {
	data := &SayFancyHelloData{}
	n, err = unmarshalCmd[*SayFancyHelloData](data,
		com.ValidatorFn[int](hw.ValidateLength), r)
	if err != nil {
		return
	}
	c.SayFancyHelloData = data
	return
}

func SizeSayFancyHelloCmdProtobuf(c SayFancyHelloCmd) (size int) {
	panic("not implemented")
}

func SkipSayFancyHelloCmdProtobuf(r muss.Reader) (n int, err error) {
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
	data := &ResultData{}
	n, err = unmarshalCmd[*ResultData](data, nil, r)
	if err != nil {
		return
	}
	result.ResultData = data
	return
}

func SizeResultProtobuf(result Result) (size int) {
	panic("not implemented")
}

func SkipResultProtobuf(r muss.Reader) (n int, err error) {
	panic("not implemented")
}

var (
	SayHelloCmdDTS = dts.New[SayHelloCmd](hw.SayHelloCmdDTM,
		muss.MarshallerFn[SayHelloCmd](MarshalSayHelloCmdProtobuf),
		muss.UnmarshallerFn[SayHelloCmd](UnmarshalSayHelloCmdProtobuf),
		muss.SizerFn[SayHelloCmd](SizeSayHelloCmdProtobuf),
		muss.SkipperFn(SkipSayHelloCmdProtobuf),
	)
	SayFancyHelloCmdDTS = dts.New[SayFancyHelloCmd](hw.SayFancyHelloCmdDTM,
		muss.MarshallerFn[SayFancyHelloCmd](MarshalSayFancyHelloCmdProtobuf),
		muss.UnmarshallerFn[SayFancyHelloCmd](UnmarshalSayFancyHelloCmdProtobuf),
		muss.SizerFn[SayFancyHelloCmd](SizeSayFancyHelloCmdProtobuf),
		muss.SkipperFn(SkipSayFancyHelloCmdProtobuf),
	)
	ResultDTS = dts.New[Result](hw.ResultDTM,
		muss.MarshallerFn[Result](MarshalResultProtobuf),
		muss.UnmarshallerFn[Result](UnmarshalResultProtobuf),
		muss.SizerFn[Result](SizeResultProtobuf),
		muss.SkipperFn(SkipResultProtobuf),
	)
)

func marshalCmd[T proto.Message](c T, w muss.Writer) (n int, err error) {
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

func unmarshalCmd[T proto.Message](d T, v com.Validator[int],
	r muss.Reader) (n int, err error) {
	l, n, err := varint.UnmarshalPositiveInt(r)
	if err != nil {
		return
	}
	if v != nil {
		if err = v.Validate(l); err != nil {
			return
		}
	}
	bs := make([]byte, l)
	n1, err := io.ReadFull(r, bs)
	n += n1
	if err != nil {
		return
	}
	err = proto.Unmarshal(bs, d)
	return
}
