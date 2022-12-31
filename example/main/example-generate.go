package main

import (
	"encoding/json"
	"fmt"
	"os"

	"melato.org/jsoncall"
	"melato.org/jsoncall/generate"
)

type ExampleInterface interface {
	A(s string, d int) (string, error)
	B() string
}

func GenerateStub() error {
	g := generate.NewGenerator()
	g.Package = "generated"
	g.Func = "NewExampleClient"
	g.OutputFile = "../generated/example.go"

	var example *ExampleInterface
	caller, err := jsoncall.NewCaller(example, nil)
	if err != nil {
		return err
	}
	return g.Output(g.GenerateClient(caller))
}

func GenerateDescriptor() error {
	var example *ExampleInterface
	desc := generate.GenerateDescriptor(example)
	data, err := json.Marshal(desc)
	if err != nil {
		return err
	}
	return os.WriteFile("example.json", data, 0666)
}

func main() {
	var err error
	if err == nil {
		err = GenerateDescriptor()
	}
	if err == nil {
		err = GenerateStub()
	}
	if err != nil {
		fmt.Printf("%v\n", err)
	}

}