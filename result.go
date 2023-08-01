package examples

type OneEchoResult EchoCmd

func (e OneEchoResult) LastOne() bool {
	return true
}
