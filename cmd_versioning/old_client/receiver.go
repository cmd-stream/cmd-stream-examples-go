package old_client

import "fmt"

// An old receiver.
type Printer struct{}

func (p Printer) Print(text string) {
	fmt.Printf("text: %v\n", text)
}
