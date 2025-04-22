package hwp

import (
	"io"

	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/dts-stream-go"
	muss "github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/varint"
	"google.golang.org/protobuf/proto"

	hw "github.com/cmd-stream/cmd-stream-examples-go/hello-world"
)

// SayHelloCmd

var SayHelloCmdProtobuf = sayHelloCmdProtobuf{}

type sayHelloCmdProtobuf struct{}

func (s sayHelloCmdProtobuf) Marshal(c SayHelloCmd, w muss.Writer) (n int,
	err error) {
	return marshalCmd(c, w)
}

func (s sayHelloCmdProtobuf) Unmarshal(r muss.Reader) (c SayHelloCmd, n int,
	err error) {
	data := &SayHelloData{}
	n, err = unmarshalCmd[*SayHelloData](data,
		com.ValidatorFn[int](hw.ValidateLength), r)
	if err != nil {
		return
	}
	c.SayHelloData = data
	return
}

func (s sayHelloCmdProtobuf) Size(c SayHelloCmd) (size int) {
	panic("not implemented")
}

func (s sayHelloCmdProtobuf) Skip(r muss.Reader) (n int,
	err error) {
	panic("not implemented")
}

// SayFancyHelloCmd

var SayFancyHelloCmdProtobuf = sayFancyHelloCmdProtobuf{}

type sayFancyHelloCmdProtobuf struct{}

func (s sayFancyHelloCmdProtobuf) Marshal(c SayFancyHelloCmd, w muss.Writer) (
	n int, err error) {
	return marshalCmd(c, w)
}

func (s sayFancyHelloCmdProtobuf) Unmarshal(r muss.Reader) (c SayFancyHelloCmd,
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

func (s sayFancyHelloCmdProtobuf) Size(c SayFancyHelloCmd) (size int) {
	panic("not implemented")
}

func (s sayFancyHelloCmdProtobuf) Skip(r muss.Reader) (n int, err error) {
	panic("not implemented")
}

// Result

var ResultProtobuf = resultProtobuf{}

type resultProtobuf struct{}

func (s resultProtobuf) Marshal(result Result, w muss.Writer) (n int,
	err error) {
	bs, err := proto.Marshal(result.ResultData)
	if err != nil {
		return
	}
	l := len(bs)
	n, err = varint.PositiveInt.Marshal(l, w)
	if err != nil {
		return
	}
	n1, err := w.Write(bs)
	n += n1
	return
}

func (s resultProtobuf) Unmarshal(r muss.Reader) (result Result, n int,
	err error) {
	data := &ResultData{}
	n, err = unmarshalCmd[*ResultData](data, nil, r)
	if err != nil {
		return
	}
	result.ResultData = data
	return
}

func (s resultProtobuf) Size(result Result) (size int) {
	panic("not implemented")
}

func (s resultProtobuf) Skip(r muss.Reader) (n int, err error) {
	panic("not implemented")
}

var (
	SayHelloCmdDTS = dts.New[SayHelloCmd](hw.SayHelloCmdDTM,
		SayHelloCmdProtobuf)
	SayFancyHelloCmdDTS = dts.New[SayFancyHelloCmd](hw.SayFancyHelloCmdDTM,
		SayFancyHelloCmdProtobuf)
	ResultDTS = dts.New[Result](hw.ResultDTM, ResultProtobuf)
)

func marshalCmd[T proto.Message](c T, w muss.Writer) (n int, err error) {
	bs, err := proto.Marshal(c)
	if err != nil {
		return
	}
	l := len(bs)
	n, err = varint.PositiveInt.Marshal(l, w)
	if err != nil {
		return
	}
	_, err = w.Write(bs)
	n += l
	return
}

func unmarshalCmd[T proto.Message](d T, v com.Validator[int],
	r muss.Reader) (n int, err error) {
	l, n, err := varint.PositiveInt.Unmarshal(r)
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
