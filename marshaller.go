package exmpls

import "github.com/cmd-stream/transport-go"

// Marshaller defines a Marshal method.
//
// This interface is used in codecs for encoding, so all commands and results
// should implement it.
type Marshaller interface {
	Marshal(w transport.Writer) error
}
