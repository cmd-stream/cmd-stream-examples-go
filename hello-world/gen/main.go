// gen/main.go

package main

import (
	"os"
	"reflect"

	hw "github.com/cmd-stream/cmd-stream-examples-go/hello-world"

	"github.com/cmd-stream/base-go"

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

// main function generates the mus-format.gen.go file, containing MUS
// serialization code for hw.SayHelloCmd, hw.SayFancyHelloCmd, hw.Greeting,
// base.Cmd and base.Result interfaces.
func main() {
	// Create a generator.
	g, err := musgen.NewFileGenerator(
		genops.WithPkgPath("github.com/cmd-stream/cmd-stream-examples-go/hello-world"),
		genops.WithPackage("hw"),
		genops.WithStream(), // We're going to generate streaming code.
	)
	assert.EqualError(err, nil)

	// ValidateLength function will be used to validate the first Command field.
	// It protects the server from excessively large payloads - if
	// deserialization fails with an validation error, the corresponding client
	// connection will be closed.
	ops := structops.WithField(typeops.WithLenValidator("ValidateLength"))

	// hw.SayHelloCmd
	sayHelloCmdType := reflect.TypeFor[hw.SayHelloCmd]()
	err = g.AddStruct(sayHelloCmdType, ops)
	assert.EqualError(err, nil)

	// With this call the generator will produce SayHelloCmdDTS variable,
	// which helps to serialize 'DTM + SayHelloCmd'. DTS stands for Data Type
	// metadata Support.
	err = g.AddDTS(sayHelloCmdType)
	assert.EqualError(err, nil)

	// hw.SayFancyHelloCmd
	sayFancyHelloCmdType := reflect.TypeFor[hw.SayFancyHelloCmd]()
	err = g.AddStruct(sayFancyHelloCmdType, ops)
	assert.EqualError(err, nil)

	err = g.AddDTS(sayFancyHelloCmdType)
	assert.EqualError(err, nil)

	// This call instructs the generator to produce serializer for the
	// base.Cmd interface.
	err = g.AddInterface(reflect.TypeFor[base.Cmd[hw.Greeter]](),
		introps.WithImpl(sayHelloCmdType),
		introps.WithImpl(sayFancyHelloCmdType),
		introps.WithMarshaller(), /// SayHelloCmd and SayFancyHelloCmd should
		// also implement the MarshallerTypedMUS interface. More on this later.
	)
	assert.EqualError(err, nil)

	// hw.Greeting
	greetingType := reflect.TypeFor[hw.Greeting]()
	err = g.AddDefinedType(greetingType)
	assert.EqualError(err, nil)

	err = g.AddDTS(greetingType)
	assert.EqualError(err, nil)

	// base.Result interface
	err = g.AddInterface(reflect.TypeFor[base.Result](),
		introps.WithImpl(greetingType),
		introps.WithMarshaller(),
	)
	assert.EqualError(err, nil)

	// Generate
	bs, err := g.Generate()
	assert.EqualError(err, nil)
	err = os.WriteFile("./mus-format.gen.go", bs, 0755)
	assert.EqualError(err, nil)
}
