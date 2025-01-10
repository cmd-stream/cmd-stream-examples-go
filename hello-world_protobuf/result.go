package main

import "github.com/cmd-stream/transport-go"

func NewResult(str string, lastOne bool) Result {
	return Result{
		ResultData: &ResultData{Str: str, LastOne: lastOne},
	}
}

type Result struct {
	*ResultData
}

func (r Result) LastOne() bool {
	return r.ResultData.LastOne
}

func (c Result) Marshal(w transport.Writer) (err error) {
	_, err = ResultDTS.Marshal(c, w)
	return
}
