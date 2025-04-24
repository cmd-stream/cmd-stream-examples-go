// result.go

package hw

import muss "github.com/mus-format/mus-stream-go"

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
