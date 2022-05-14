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
		jsoncall.TraceDebug = true
	}
	return nil
}

func (t *Server) Receiver(c jsoncall.ReceiverContext) (interface{}, error) {
	return &demo.DemoImpl{}, nil
}

func (t *Server) Run() error {
	caller, err := demo.NewCaller()
	if err != nil {
		return err
	}
	handler := jsoncall.NewHttpHandler(caller, t.Receiver)
	addr := fmt.Sprintf(":%d", t.Port)
	fmt.Printf("starting server at %s\n", addr)
	return http.ListenAndServe(addr, handler)
}

func main() {
	var cmd command.SimpleCommand
	var serverOps Server
	cmd.Flags(&serverOps).RunFunc(serverOps.Run)
	command.Main(&cmd)
}
