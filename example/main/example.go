package main

import (
	"fmt"
	"net/http"

	"melato.org/command"
	"melato.org/jsoncall"
	"melato.org/jsoncall/example/generated"
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

func NewExampleClient() (ExampleInterface, error) {
	var example *ExampleInterface
	caller, err := jsoncall.NewCaller(example, nil)
	if err != nil {
		return nil, err
	}
	c := caller.NewHttpClient("http://localhost:8080/")
	return generated.NewExampleInterface(c), nil
}

func ExampleClientWithGeneratedCode() error {
	client, err := NewExampleClient()
	if err != nil {
		return err
	}
	s, err := client.A("hello", 7)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", s)
	return nil
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
	cmd.Command("generated-client").RunFunc(ExampleClientWithGeneratedCode)
	command.Main(cmd)
}
