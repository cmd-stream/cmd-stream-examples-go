package main

// Result implements the base.Result interface. The client will wait for more
// command results if the LastOne method of the received result returns false.
type Result int

func (r Result) LastOne() bool {
	return true
}
