package main

import (
	"fmt"
	"net/http"

	"melato.org/command"
	"melato.org/jsoncall"
	"melato.org/jsoncall/generate"
)

type ExampleInterface interface {
	A(s string, d int) (string, error)
}

type Example struct {
}

func (t *Example) A(s string, d int) (string, error) {
	return fmt.Sprintf("%s:%d", s, d), nil
}

func (t *Example) B() (string, error) {
	return "b", nil
}

func ExampleServer() error {
	var handler http.Handler
	handler, err := jsoncall.NewHttpHandler(&Example{}, nil)
	if err != nil {
		return err
	}
	return http.ListenAndServe(":8080", handler)
}

func ExampleServerWithInterface() error {
	var handler http.Handler
	var api *ExampleInterface
	handler, err := jsoncall.NewHttpHandler(&Example{}, api)
	if err != nil {
		return err
	}
	return http.ListenAndServe(":8080", handler)
}

func ExampleClient() error {
	var example *ExampleInterface
	caller, err := jsoncall.NewCaller(example, nil)
	if err != nil {
		return err
	}
	client := caller.NewHttpClient("http://localhost:8080/")

	var response map[string]any
	err = client.Call(&response, "A", "hello", 2)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", response["result"])
	return nil
}

func GenerateStub() error {
	g := generate.NewGenerator()
	g.Package = "client"
	g.Type = "ExampleClient"
	g.OutputFile = "../client/generated_example.go"

	var example *ExampleInterface
	caller, err := jsoncall.NewCaller(example, nil)
	if err != nil {
		return err
	}
	return g.Output(g.GenerateClient(caller))
}

func main() {
	/*
		jsoncall.TraceInit = true
		jsoncall.TraceCalls = true
		jsoncall.TraceData = true
	*/
	cmd := &command.SimpleCommand{}
	cmd.Command("server").RunFunc(ExampleServer)
	cmd.Command("server-interface").RunFunc(ExampleServerWithInterface)
	cmd.Command("client").RunFunc(ExampleClient)
	cmd.Command("generate").RunFunc(GenerateStub)
	command.Main(cmd)
}
