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

func (t *Server) Run(handler http.Handler) error {
	addr := fmt.Sprintf(":%d", t.Port)
	fmt.Printf("starting server at %s\n", addr)
	return http.ListenAndServe(addr, handler)
}

func (t *Server) RunDemo() error {
	caller, err := demo.NewCaller()
	if err != nil {
		return err
	}
	handler := jsoncall.NewHttpHandler(caller, t.DemoReceiver)
	return t.Run(handler)
}

func (t *Server) RunMath() error {
	var math *demo.Math
	caller, err := jsoncall.NewCaller(math, nil)
	if err != nil {
		return err
	}
	handler := jsoncall.NewHttpHandler(caller, t.MathReceiver)
	return t.Run(handler)
}

// RunMux demostrates how to use http.ServeMux to implement a server that handles multiple interfaces
func (t *Server) RunMux() error {
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

	// we could add more handlers.
	return t.Run(mux)
}

func main() {
	var cmd command.SimpleCommand
	var serverOps Server
	cmd.Flags(&serverOps)
	cmd.Command("demo").RunFunc(serverOps.RunDemo)
	cmd.Command("math").RunFunc(serverOps.RunMath)
	cmd.Command("mux").RunFunc(serverOps.RunMux)
	command.Main(&cmd)
}
