// dtm.go

//go:generate go run gen/main.go
package hw

import (
	com "github.com/mus-format/common-go"
)

const (
	SayHelloCmdDTM com.DTM = iota + 1
	SayFancyHelloCmdDTM
)

const (
	GreetingDTM com.DTM = iota + 1
)
