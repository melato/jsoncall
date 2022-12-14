package server

import (
	"melato.org/jsoncall/demo"
	"melato.org/jsoncall/generate"
)

// Program that generates the client stub.  Needed only at development time.
// Usage: go run generate.go client
type GenerateOp struct {
	generate.Generator
}

func (t *GenerateOp) Init() error {
	g := &t.Generator
	g.Init()
	g.Package = "client"
	g.Type = "DemoClient"
	g.OutputFile = "../client/generated_demo.go"
	return nil
}

func (t *GenerateOp) Generate() error {
	c, err := demo.NewCaller()
	if err != nil {
		return err
	}
	return t.Output(t.GenerateClient(c))

}

type NamesOp struct {
	File string `name:"f" usage:"names JSON file"`
}

func (t *NamesOp) Init() error {
	t.File = "../demo.json"
	return nil
}

func (t *NamesOp) UpdateNames() error {
	var api *demo.Demo
	return generate.UpdateDescriptor(api, t.File)
}
