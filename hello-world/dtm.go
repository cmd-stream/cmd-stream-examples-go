//go:generate go run gen/main.go
package hw

import (
	com "github.com/mus-format/common-go"
)

const (
	SayHelloCmdDTM com.DTM = iota
	SayFancyHelloCmdDTM
)

const (
	ResultDTM com.DTM = iota
)
