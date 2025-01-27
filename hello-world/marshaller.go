package hw

import (
	"github.com/cmd-stream/transport-go"
)

// Marshaller defines a Marshal method.
type Marshaller interface {
	Marshal(w transport.Writer) error
}
