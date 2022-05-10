package main

import (
	"melato.org/command"
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
	g.Type = "GeneratedClient"
	g.OutputFile = "../client/generated_client.go"
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

func main() {
	var cmd command.SimpleCommand
	var generateOp GenerateOp
	cmd.Command("client").Flags(&generateOp).RunFunc(generateOp.Generate)
	var namesOp NamesOp
	cmd.Command("names").Flags(&namesOp).RunFunc(namesOp.UpdateNames)
	command.Main(&cmd)
}
