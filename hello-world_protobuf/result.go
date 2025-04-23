package hwp

import muss "github.com/mus-format/mus-stream-go"

func NewGreeting(str string) Greeting {
	return Greeting{
		GreetingData: &GreetingData{Str: str},
	}
}

type Greeting struct {
	*GreetingData
}

func (r Greeting) String() string {
	return r.GreetingData.Str
}

func (r Greeting) LastOne() bool {
	return true
}

func (c Greeting) MarshalTypedProtobuf(w muss.Writer) (n int, err error) {
	return GreetingDTS.Marshal(c, w)
}

func (c Greeting) SizeTypedProtobuf() (size int) {
	return GreetingDTS.Size(c)
}
