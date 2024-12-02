package main

type Calculator struct{}

func (c Calculator) Add(n1, n2 int64) int64 {
	return n1 + n2
}

func (c Calculator) Sub(n1, n2 int64) int64 {
	return n1 - n2
}
