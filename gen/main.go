package main

import (
	"os"
	"reflect"

	exmpls "github.com/cmd-stream/cmd-stream-examples-go"
	"github.com/mus-format/musgen-go/basegen"
	musgen "github.com/mus-format/musgen-go/mus"
)

func main() {
	g, err := musgen.NewFileGenerator(basegen.Conf{
		Package: "exmpls",
		Stream:  true,
	})
	if err != nil {
		panic(err)
	}
	err = g.AddStructDTS(reflect.TypeFor[exmpls.SayHelloCmd]())
	if err != nil {
		panic(err)
	}
	err = g.AddStructDTS(reflect.TypeFor[exmpls.SayFancyHelloCmd]())
	if err != nil {
		panic(err)
	}
	err = g.AddStructDTS(reflect.TypeFor[exmpls.UnsupportedCmd]())
	if err != nil {
		panic(err)
	}
	err = g.AddStructDTS(reflect.TypeFor[exmpls.SayFancyHelloMultiCmd]())
	if err != nil {
		panic(err)
	}
	err = g.AddStructDTS(reflect.TypeFor[exmpls.OldSayHelloCmd]())
	if err != nil {
		panic(err)
	}
	err = g.AddStructDTS(reflect.TypeFor[exmpls.Result]())
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
