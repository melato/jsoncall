package main

import (
	"fmt"

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

func main() {
	err := GenerateStub()
	if err != nil {
		fmt.Printf("%v\n", err)
	}

}
