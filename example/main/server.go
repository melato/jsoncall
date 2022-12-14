package main

import (
	_ "embed"
	"fmt"
	"net/http"

	"melato.org/command"
	"melato.org/command/usage"
	"melato.org/jsoncall"
	"melato.org/jsoncall/example"
	"melato.org/jsoncall/example/server"
)

//go:embed server.yaml
var usageData []byte

type Server struct {
	Port  int32
	Trace bool
}

func (t *Server) Init() error {
	t.Port = 8080
	return nil
}

func (t *Server) Configured() error {
	if t.Trace {
		jsoncall.TraceCalls = true
		jsoncall.TraceInit = true
	}
	return nil
}

func (t *Server) DemoReceiver(c jsoncall.ReceiverContext) (interface{}, error) {
	return &example.DemoImpl{}, nil
}

func (t *Server) MathReceiver(c jsoncall.ReceiverContext) (interface{}, error) {
	return &example.MathImpl{}, nil
}

// RunMux demostrates how to use http.ServeMux to implement a server that handles multiple interfaces
func (t *Server) Run() error {
	mux := http.NewServeMux()
	demoCaller, err := example.NewCaller()
	if err != nil {
		return err
	}
	mux.Handle("/demo/", demoCaller.NewHttpHandler(t.DemoReceiver))

	var math *example.Math
	mathHandler, err := jsoncall.NewHttpHandler(&example.MathImpl{}, math)
	if err != nil {
		return err
	}
	mux.Handle("/math/", mathHandler)

	addr := fmt.Sprintf(":%d", t.Port)
	fmt.Printf("starting server at %s\n", addr)
	return http.ListenAndServe(addr, mux)
}

func (t *Server) RunDirect() error {
	handler, err := jsoncall.NewHttpHandler(&example.MathImpl{}, nil)
	if err != nil {
		return err
	}
	return http.ListenAndServe(":8080", handler)
}

func (t *Server) RunInterface() error {
	var api *example.Math
	handler, err := jsoncall.NewHttpHandler(&example.MathImpl{}, api)
	if err != nil {
		return err
	}
	return http.ListenAndServe(":8080", handler)
}

func GenerateCommand() *command.SimpleCommand {
	var cmd command.SimpleCommand
	var generateOp server.GenerateOp
	cmd.Command("clients").Flags(&generateOp).RunFunc(generateOp.Generate)
	var namesOp server.NamesOp
	cmd.Command("descriptors").Flags(&namesOp).RunFunc(namesOp.UpdateNames)
	return &cmd
}

func main() {
	var cmd command.SimpleCommand
	var serverOps Server
	cmd.Flags(&serverOps)
	cmd.Command("run").RunFunc(serverOps.Run)
	cmd.Command("receiver").RunFunc(serverOps.RunDirect)
	cmd.Command("interface").RunFunc(serverOps.RunInterface)
	cmd.AddCommand("generate", GenerateCommand())
	usage.Apply(&cmd, usageData)
	command.Main(&cmd)
}
