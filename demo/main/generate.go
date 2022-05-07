package main

import (
	"melato.org/command"
	"melato.org/jsoncall/demo"
	"melato.org/jsoncall/generate"
)

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

func (t *GenerateOp) WriteNames(file string) error {
	var api *demo.Demo
	names := generate.GenerateMethodNames(api)
	return generate.WriteMethodNamesJSON(names, file)
}

func main() {
	var cmd command.SimpleCommand
	var generateOp GenerateOp
	cmd.Command("client").Flags(&generateOp).RunFunc(generateOp.Generate)
	cmd.Command("names").RunFunc(generateOp.WriteNames)
	command.Main(&cmd)
}
