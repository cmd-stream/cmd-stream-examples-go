package hw

import (
	com "github.com/mus-format/common-go"
	muss "github.com/mus-format/mus-stream-go"
)

const GreetingDTM com.DTM = iota

type Greeting string

func (g Greeting) LastOne() bool {
	return true
}

func (g Greeting) String() string {
	return string(g)
}

func (g Greeting) MarshalTypedMUS(w muss.Writer) (n int, err error) {
	return GreetingDTS.Marshal(g, w)
}

func (g Greeting) SizeTypedMUS() (size int) {
	return GreetingDTS.Size(g)
}
