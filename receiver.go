package exmpls

import (
	"strings"
)

// NewGreeter creates a new Greeter.
func NewGreeter(interjection, adjective, sep string) Greeter {
	return Greeter{
		interjection: interjection,
		adjective:    adjective,
		sep:          " ",
	}
}

// Greeter represents a Receiver.
type Greeter struct {
	interjection string
	adjective    string
	sep          string
}

func (g Greeter) Interjection() string {
	return g.interjection
}

func (g Greeter) Adjective() string {
	return g.adjective
}

func (g Greeter) Join(strs ...string) string {
	return strings.Join(strs, g.sep)
}

// OldGreeter represents an old Receiver used by the server.
type OldGreeter struct{}

func (g OldGreeter) SayHello(str string) string {
	return "Hello " + str
}
