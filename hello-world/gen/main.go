package main

import (
	"os"
	"reflect"

	"github.com/cmd-stream/base-go"
	hw "github.com/cmd-stream/cmd-stream-examples-go/hello-world"

	musgen "github.com/mus-format/musgen-go/mus"
	genops "github.com/mus-format/musgen-go/options/generate"
	introps "github.com/mus-format/musgen-go/options/interface"
	structops "github.com/mus-format/musgen-go/options/struct"
	typeops "github.com/mus-format/musgen-go/options/type"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

// The main function generates the mus-format.gen.go file containing
// MUS serialization code for SayHelloCmd, SayFancyHelloCmd, and Greeting.
func main() {
	// Create a generator.
	g, err := musgen.NewFileGenerator(
		genops.WithPkgPath("github.com/cmd-stream/cmd-stream-examples-go/hello-world"),
		genops.WithPackage("hw"),
		genops.WithStream(),
	)
	assert.EqualError(err, nil)

	// SayHelloCmd.
	sayHelloCmdType := reflect.TypeFor[hw.SayHelloCmd]()
	err = g.AddStruct(sayHelloCmdType,
		// Specifies options for the first field.
		structops.WithField(
			// Specifies the length validator for the first field.
			typeops.WithLenValidator("ValidateLength"), // Where ValidateLength is the function name.
		),
	)
	assert.EqualError(err, nil)

	// This call instructs the generator to produce the SayHelloCmdDTS variable,
	// which facilitates the serialization and deserialization of 'DTM + SayHelloCmd'.
	// DTS stands for Data Type Metadata Support.
	err = g.AddDTS(sayHelloCmdType)
	assert.EqualError(err, nil)

	// SayFancyHelloCmd.
	sayFancyHelloCmdType := reflect.TypeFor[hw.SayFancyHelloCmd]()
	err = g.AddStruct(sayFancyHelloCmdType, structops.WithField(
		typeops.WithLenValidator("ValidateLength"),
	))
	assert.EqualError(err, nil)

	err = g.AddDTS(sayFancyHelloCmdType)
	assert.EqualError(err, nil)

	// base.Cmd[hw.Greeter]
	err = g.AddInterface(reflect.TypeFor[base.Cmd[hw.Greeter]](),
		introps.WithImpl(sayHelloCmdType),
		introps.WithImpl(sayFancyHelloCmdType),
		introps.WithMarshaller(),
	)
	assert.EqualError(err, nil)

	// Greeting.
	greetingType := reflect.TypeFor[hw.Greeting]()
	err = g.AddDefinedType(greetingType)
	assert.EqualError(err, nil)

	err = g.AddDTS(greetingType)
	assert.EqualError(err, nil)

	// base.Result
	err = g.AddInterface(reflect.TypeFor[base.Result](),
		introps.WithImpl(greetingType),
	)
	assert.EqualError(err, nil)

	// Generate.
	bs, err := g.Generate()
	assert.EqualError(err, nil)
	err = os.WriteFile("./mus-format.gen.go", bs, 0755)
	assert.EqualError(err, nil)
}
