package main

import (
	"os"
	"reflect"

	hw "cmd-stream-examples-go/hello-world"

	musgen "github.com/mus-format/musgen-go/mus"
	genops "github.com/mus-format/musgen-go/options/generate"
	structops "github.com/mus-format/musgen-go/options/struct"
	typeops "github.com/mus-format/musgen-go/options/type"
)

// The main function generates the mus-format.gen.go file containing
// MUS serialization code for SayHelloCmd, SayFancyHelloCmd, and Result.
func main() {
	// Create a generator.
	g := musgen.NewFileGenerator(
		genops.WithPackage("hw"),
		genops.WithStream(),
	)

	// SayHelloCmd.
	t := reflect.TypeFor[hw.SayHelloCmd]()
	err := g.AddStruct(t,
		// Specifies options for the first field.
		structops.WithField(
			// Specifies the length validator for the first field.
			typeops.WithLenValidator("ValidateLength"), // Where ValidateLength is the function name.
		),
	)
	if err != nil {
		panic(err)
	}
	// This call instructs the generator to produce the SayHelloCmdDTS variable,
	// which facilitates the serialization and deserialization of 'DTM + SayHelloCmd'.
	// DTS stands for Data Type Metadata Support.
	err = g.AddDTS(t)
	if err != nil {
		panic(err)
	}

	// SayFancyHelloCmd.
	t = reflect.TypeFor[hw.SayFancyHelloCmd]()
	err = g.AddStruct(t, structops.WithField(
		typeops.WithLenValidator("ValidateLength"),
	))
	if err != nil {
		panic(err)
	}
	err = g.AddDTS(t)
	if err != nil {
		panic(err)
	}

	// Result.
	t = reflect.TypeFor[hw.Result]()
	err = g.AddStruct(t)
	if err != nil {
		panic(err)
	}
	err = g.AddDTS(t)
	if err != nil {
		panic(err)
	}

	// Generate.
	bs, err := g.Generate()
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("./mus-format.gen.go", bs, 0755)
	if err != nil {
		panic(err)
	}
}
