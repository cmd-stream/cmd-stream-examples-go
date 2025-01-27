package hwp

import "github.com/cmd-stream/transport-go"

func NewResult(str string) Result {
	return Result{
		ResultData: &ResultData{Str: str},
	}
}

type Result struct {
	*ResultData
}

func (r Result) Greeting() string {
	return r.ResultData.Str
}

func (r Result) LastOne() bool {
	return true
}

func (c Result) Marshal(w transport.Writer) (err error) {
	_, err = ResultDTS.Marshal(c, w)
	return
}
