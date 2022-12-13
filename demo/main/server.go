package main

import (
	"fmt"
	"net/http"

	"melato.org/command"
	"melato.org/jsoncall"
	"melato.org/jsoncall/demo"
)

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
	return &demo.DemoImpl{}, nil
}

func (t *Server) MathReceiver(c jsoncall.ReceiverContext) (interface{}, error) {
	return &demo.MathImpl{}, nil
}

// RunMux demostrates how to use http.ServeMux to implement a server that handles multiple interfaces
func (t *Server) Run() error {
	mux := http.NewServeMux()
	demoCaller, err := demo.NewCaller()
	if err != nil {
		return err
	}
	mux.Handle("/demo/", jsoncall.NewHttpHandler(demoCaller, t.DemoReceiver))

	var math *demo.Math
	mathCaller, err := jsoncall.NewCaller(math, nil)
	if err != nil {
		return err
	}
	mux.Handle("/math/", jsoncall.NewHttpHandler(mathCaller, t.MathReceiver))

	addr := fmt.Sprintf(":%d", t.Port)
	fmt.Printf("starting server at %s\n", addr)
	return http.ListenAndServe(addr, mux)
}

func main() {
	var cmd command.SimpleCommand
	var serverOps Server
	cmd.Flags(&serverOps)
	cmd.Command("run").RunFunc(serverOps.Run)
	command.Main(&cmd)
}
