package main

import "fmt"

type Printer struct{}

func (p Printer) Print(from string, text string) {
	fmt.Printf("from: %v, text: %v\n", from, text)
}
