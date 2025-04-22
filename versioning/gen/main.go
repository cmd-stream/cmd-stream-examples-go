package main

import (
	"os"
	"reflect"

	versioning "github.com/cmd-stream/cmd-stream-examples-go/versioning"

	musgen "github.com/mus-format/musgen-go/mus"
	genops "github.com/mus-format/musgen-go/options/generate"
)

// The main function generates the mus-format.gen.go file containing
// MUS serialization code for SayHelloCmd, SayFancyHelloCmd, and Result.
func main() {
	// Create a generator.
	g, err := musgen.NewFileGenerator(
		genops.WithPkgPath("github.com/cmd-stream/cmd-stream-examples-go/versioning"),
		genops.WithPackage("versioning"),
		genops.WithStream(), // We're going to generate streaming code.
	)
	if err != nil {
		panic(err)
	}

	// OldSayHelloCmd.
	t := reflect.TypeFor[versioning.OldSayHelloCmd]()
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
