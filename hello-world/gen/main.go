package main

import (
	"os"
	"reflect"

	hw "cmd-stream-examples-go/hello-world"

	"github.com/mus-format/musgen-go/basegen"
	musgen "github.com/mus-format/musgen-go/mus"
)

// main function generates the mus-format.gen.go file with MUS serialization
// code for SayHelloCmd, SayFancyHelloCmd and Result.
func main() {
	// Create a generator.
	g, err := musgen.NewFileGenerator(basegen.Conf{
		Package: "hw",
		Stream:  true, // We're going to generate streaming code.
	})
	if err != nil {
		panic(err)
	}
	// These options specify a validator for the first string field of the structure.
	opts := basegen.StructOptions{
		basegen.StringFieldOptions{
			StringOptions: basegen.StringOptions{
				LenValidator: "ValidateLength", // Where ValidateLength is the function name.
			},
		},
	}
	// With this call the generator will produce SayHelloCmdDTS variable, which
	// helps to serialize/deserialize 'DTM + SayHelloCmd'. DTS stands for Data
	// Type Metadata Support.
	err = g.AddStructDTSWith(reflect.TypeFor[hw.SayHelloCmd](), "", opts)
	if err != nil {
		panic(err)
	}
	err = g.AddStructDTSWith(reflect.TypeFor[hw.SayFancyHelloCmd](), "", opts)
	if err != nil {
		panic(err)
	}
	err = g.AddStructDTS(reflect.TypeFor[hw.Result]())
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
