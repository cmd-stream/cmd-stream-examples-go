// receiver.go

package hw

import "strings"

// NewGreeter creates a new Greeter.
func NewGreeter(interjection, adjective, sep string) Greeter {
	return Greeter{
		interjection: interjection,
		adjective:    adjective,
		sep:          " ",
	}
}

// Greeter represents a Receiver and provides the functionality for
// creating greetings.
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
