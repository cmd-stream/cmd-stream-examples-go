//go:generate go run gen/main.go
package server_streaming

import (
	com "github.com/mus-format/common-go"
)

const (
	SayFancyHelloMultiCmdDTM com.DTM = iota + 10
	ResultDTM
)
