package main

import (
	"os"
	"reflect"

	versioning "cmd-stream-examples-go/versioning"

	"github.com/mus-format/musgen-go/basegen"
	musgen "github.com/mus-format/musgen-go/mus"
)

// main function will generate the mus-format.gen.go file with MUS serialization
// code for SayHelloCmd, SayFancyHelloCmd and Result.
func main() {
	// Create a generator.
	g, err := musgen.NewFileGenerator(basegen.Conf{
		Package: "versioning",
		Stream:  true, // We're going to generate streaming code.
	})
	if err != nil {
		panic(err)
	}
	// With this call the generator will produce OldSayHelloCmdDTS variable, which
	// helps to serialize/deserialize 'DTM + OldSayHelloCmd'. DTS stands for Data
	// Type Metadata Support.
	err = g.AddStructDTS(reflect.TypeFor[versioning.OldSayHelloCmd]())
	if err != nil {
		panic(err)
	}
	bs, err := g.Generate()
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("./mus-format.gen.go", bs, 0755)
	if err != nil {
		panic(err)
	}
}
