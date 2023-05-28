package main

import (
	_ "embed"
	"fmt"
	"net/http"

	"melato.org/command"
	"melato.org/command/usage"
	"melato.org/jsoncall"
	"melato.org/jsoncall/example"
)

//go:embed demo-server.yaml
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
		jsoncall.TraceData = true
	}
	return nil
}

func (t *Server) DemoReceiver(w http.ResponseWriter, r *http.Request) interface{} {
	return &example.DemoImpl{}
}

func (t *Server) MathReceiver(w http.ResponseWriter, r *http.Request) interface{} {
	return &example.MathImpl{}
}

// RunMux demostrates how to use http.ServeMux to implement a server that handles multiple interfaces
func (t *Server) RunFull() error {
	demoCaller, err := example.NewDemoCaller()
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.Handle("/demo/", demoCaller.NewHttpHandler(t.DemoReceiver))

	var math *example.Math
	mathHandler, err := jsoncall.NewHttpHandler(math)
	if err != nil {
		return err
	}
	mathHandler.SetReceiver((&example.MathImpl{}))
	mux.Handle("/math/", mathHandler)

	addr := fmt.Sprintf(":%d", t.Port)
	fmt.Printf("starting server at %s\n", addr)
	return http.ListenAndServe(addr, mux)
}

func (t *Server) RunDirect() error {
	handler, err := jsoncall.NewHttpHandler(&example.MathImpl{})
	if err != nil {
		return err
	}
	handler.SetReceiver((&example.MathImpl{}))
	return http.ListenAndServe(":8080", handler)
}

func (t *Server) RunInterface() error {
	var api *example.Math
	handler, err := jsoncall.NewHttpHandler(api)
	if err != nil {
		return err
	}
	handler.SetReceiver((&example.MathImpl{}))
	return http.ListenAndServe(":8080", handler)
}

func main() {
	var cmd command.SimpleCommand
	var serverOps Server
	cmd.Flags(&serverOps)
	cmd.Command("full").RunFunc(serverOps.RunFull)
	cmd.Command("receiver").RunFunc(serverOps.RunDirect)
	cmd.Command("interface").RunFunc(serverOps.RunInterface)
	usage.Apply(&cmd, usageData)
	command.Main(&cmd)
}
