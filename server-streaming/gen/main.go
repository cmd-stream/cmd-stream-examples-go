package main

import (
	"os"
	"reflect"

	server_streaming "github.com/cmd-stream/cmd-stream-examples-go/server-streaming"

	musgen "github.com/mus-format/musgen-go/mus"
	genops "github.com/mus-format/musgen-go/options/generate"
)

// The main function generates the mus-format.gen.go file containing MUS
// serialization code for SayHelloCmd, SayFancyHelloCmd, and Result.
func main() {
	// Create a generator.
	g, err := musgen.NewFileGenerator(
		genops.WithPkgPath("github.com/cmd-stream/cmd-stream-examples-go/server-streaming"),
		genops.WithPackage("streaming"),
		genops.WithStream())
	if err != nil {
		panic(err)
	}

	// SayFancyHelloMultiCmd.
	t := reflect.TypeFor[server_streaming.SayFancyHelloMultiCmd]()
	err = g.AddStruct(t)
	if err != nil {
		panic(err)
	}
	err = g.AddDTS(t)
	if err != nil {
		panic(err)
	}

	// Result.
	t = reflect.TypeFor[server_streaming.Result]()
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
