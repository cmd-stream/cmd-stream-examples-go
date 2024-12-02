package main

import muss "github.com/mus-format/mus-stream-go"

type MarshallerProtobuf interface {
	MarshalProtobuf(w muss.Writer) (n int, err error)
}
