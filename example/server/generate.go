package server

import (
	"melato.org/jsoncall/example"
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
	g.Package = "generated"
	g.Type = "DemoClient"
	g.OutputFile = "../generated/generated_demo.go"
	g.Imports = []string{"melato.org/jsoncall/example"}
	return nil
}

func (t *GenerateOp) Generate() error {
	c, err := example.NewDemoCaller()
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
	var api *example.Demo
	return generate.UpdateDescriptor(api, t.File)
}

func Generate() error {
	var g GenerateOp
	_ = g.Init()
	err := g.Generate()
	if err != nil {
		return err
	}
	var names NamesOp
	_ = names.Init()
	err = names.UpdateNames()
	if err != nil {
		return err
	}
	return nil
}
