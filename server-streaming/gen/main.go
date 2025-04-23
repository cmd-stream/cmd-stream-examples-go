package main

import (
	"os"
	"reflect"

	"github.com/cmd-stream/base-go"
	hw "github.com/cmd-stream/cmd-stream-examples-go/hello-world"
	server_streaming "github.com/cmd-stream/cmd-stream-examples-go/server-streaming"
	assert "github.com/ymz-ncnk/assert/panic"

	musgen "github.com/mus-format/musgen-go/mus"
	genops "github.com/mus-format/musgen-go/options/generate"
	introps "github.com/mus-format/musgen-go/options/interface"
)

func init() {
	assert.On = true
}

// The main function generates the mus-format.gen.go file containing MUS
// serialization code for SayHelloCmd, SayFancyHelloCmd, and Result.
func main() {
	// Create a generator.
	g, err := musgen.NewFileGenerator(
		genops.WithPkgPath("github.com/cmd-stream/cmd-stream-examples-go/server-streaming"),
		genops.WithPackage("streaming"),
		genops.WithImportAlias("github.com/cmd-stream/cmd-stream-examples-go/hello-world",
			"hw"),
		genops.WithStream())
	assert.EqualError(err, nil)

	// SayFancyHelloMultiCmd.
	sayFancyHelloMultiCmdType := reflect.TypeFor[server_streaming.SayFancyHelloMultiCmd]()
	err = g.AddStruct(sayFancyHelloMultiCmdType)
	assert.EqualError(err, nil)

	err = g.AddDTS(sayFancyHelloMultiCmdType)
	assert.EqualError(err, nil)

	err = g.AddInterface(reflect.TypeFor[base.Cmd[hw.Greeter]](),
		introps.WithImpl(sayFancyHelloMultiCmdType))
	assert.EqualError(err, nil)

	// Greeting.
	greetingType := reflect.TypeFor[server_streaming.Greeting]()
	err = g.AddStruct(greetingType)
	assert.EqualError(err, nil)

	err = g.AddDTS(greetingType)
	assert.EqualError(err, nil)

	err = g.AddInterface(reflect.TypeFor[base.Result](),
		introps.WithImpl(greetingType))
	assert.EqualError(err, nil)

	// Generate.
	bs, err := g.Generate()
	// fmt.Println(string(err.(*musgen.TmplEngineError).ByteSlice()))
	// fmt.Println(string(err.(*musgen.TmplEngineError).Error()))
	assert.EqualError(err, nil)
	err = os.WriteFile("./mus-format.gen.go", bs, 0755)
	assert.EqualError(err, nil)
}
