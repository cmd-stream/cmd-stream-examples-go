package hwp

import (
	"fmt"
	"io"
	reflect "reflect"

	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/dts-stream-go"
	exts "github.com/mus-format/ext-protobuf-stream-go"
	muss "github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/varint"
	"google.golang.org/protobuf/proto"

	"github.com/cmd-stream/base-go"
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

// Greeting

var GreetingProtobuf = greetingProtobuf{}

type greetingProtobuf struct{}

func (s greetingProtobuf) Marshal(result Greeting, w muss.Writer) (n int,
	err error) {
	bs, err := proto.Marshal(result.GreetingData)
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

func (s greetingProtobuf) Unmarshal(r muss.Reader) (result Greeting, n int,
	err error) {
	data := &GreetingData{}
	n, err = unmarshalCmd[*GreetingData](data, nil, r)
	if err != nil {
		return
	}
	result.GreetingData = data
	return
}

func (s greetingProtobuf) Size(result Greeting) (size int) {
	panic("not implemented")
}

func (s greetingProtobuf) Skip(r muss.Reader) (n int, err error) {
	panic("not implemented")
}

var (
	SayHelloCmdDTS = dts.New[SayHelloCmd](hw.SayHelloCmdDTM,
		SayHelloCmdProtobuf)
	SayFancyHelloCmdDTS = dts.New[SayFancyHelloCmd](hw.SayFancyHelloCmdDTM,
		SayFancyHelloCmdProtobuf)
	GreetingDTS = dts.New[Greeting](hw.GreetingDTM, GreetingProtobuf)
)

// base.Cmd --------------------------------------------------------------------

var CmdProtobuf = cmdProtobuf{}

type cmdProtobuf struct{}

func (s cmdProtobuf) Marshal(v base.Cmd[hw.Greeter], w muss.Writer) (n int,
	err error) {
	if m, ok := v.(exts.MarshallerTypedProtobuf); ok {
		return m.MarshalTypedProtobuf(w)
	}
	panic(fmt.Sprintf("%v doesn't implement the exts.MarshallerTypedProtobuf interface",
		reflect.TypeOf(v)))
}

func (s cmdProtobuf) Unmarshal(r muss.Reader) (v base.Cmd[hw.Greeter], n int, err error) {
	dtm, n, err := dts.DTMSer.Unmarshal(r)
	if err != nil {
		return
	}
	var n1 int
	switch dtm {
	case hw.SayHelloCmdDTM:
		v, n1, err = SayHelloCmdDTS.UnmarshalData(r)
	case hw.SayFancyHelloCmdDTM:
		v, n1, err = SayFancyHelloCmdDTS.UnmarshalData(r)
	default:
		err = fmt.Errorf("unexpected %v DTM", dtm)
		return
	}
	n += n1
	return
}

func (s cmdProtobuf) Size(v base.Cmd[hw.Greeter]) (size int) {
	if m, ok := v.(exts.MarshallerTypedProtobuf); ok {
		return m.SizeTypedProtobuf()
	}
	panic(fmt.Sprintf("%v doesn't implement the exts.MarshallerTypedProtobuf interface",
		reflect.TypeOf(v)))
}

func (s cmdProtobuf) Skip(r muss.Reader) (n int, err error) {
	dtm, n, err := dts.DTMSer.Unmarshal(r)
	if err != nil {
		return
	}
	var n1 int
	switch dtm {
	case hw.SayHelloCmdDTM:
		n1, err = SayHelloCmdDTS.SkipData(r)
	case hw.SayFancyHelloCmdDTM:
		n1, err = SayFancyHelloCmdDTS.SkipData(r)
	default:
		err = fmt.Errorf("unexpected %v DTM", dtm)
		return
	}
	n += n1
	return
}

// base.Result -----------------------------------------------------------------

var ResultProtobuf = resultProtobuf{}

type resultProtobuf struct{}

func (s resultProtobuf) Marshal(v base.Result, w muss.Writer) (n int, err error) {
	if m, ok := v.(exts.MarshallerTypedProtobuf); ok {
		return m.MarshalTypedProtobuf(w)
	}
	panic(fmt.Sprintf("%v doesn't implement the exts.MarshallerTypedProtobuf interface",
		reflect.TypeOf(v)))
}

func (s resultProtobuf) Unmarshal(r muss.Reader) (v base.Result, n int, err error) {
	dtm, n, err := dts.DTMSer.Unmarshal(r)
	if err != nil {
		return
	}
	var n1 int
	switch dtm {
	case hw.GreetingDTM:
		v, n1, err = GreetingDTS.UnmarshalData(r)
	default:
		err = fmt.Errorf("unexpected %v DTM", dtm)
		return
	}
	n += n1
	return
}

func (s resultProtobuf) Size(v base.Result) (size int) {
	if m, ok := v.(exts.MarshallerTypedProtobuf); ok {
		return m.SizeTypedProtobuf()
	}
	panic(fmt.Sprintf("%v doesn't implement the exts.MarshallerTypedProtobuf interface",
		reflect.TypeOf(v)))
}

func (s resultProtobuf) Skip(r muss.Reader) (n int, err error) {
	dtm, n, err := dts.DTMSer.Unmarshal(r)
	if err != nil {
		return
	}
	var n1 int
	switch dtm {
	case hw.GreetingDTM:
		n1, err = GreetingDTS.SkipData(r)
	default:
		err = fmt.Errorf("unexpected %v DTM", dtm)
		return
	}
	n += n1
	return
}

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
