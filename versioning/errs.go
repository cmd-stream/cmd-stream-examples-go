package main

import "errors"

var ErrUnsupportedCmdType = errors.New("unsupported command type")
var ErrUnsupportedResultType = errors.New("unsupported result type")
