package main

import (
	"net/http"

	"melato.org/command"
	"melato.org/jsoncall"
	"melato.org/jsoncall/demo"
)

type Server struct {
	Port       int32
	ConfigFile string `name:"c" usage:"database config file"`
	Trace      bool
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

func (t *Server) Receiver(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return &demo.DemoImpl{}, nil
}

func (t *Server) Run() error {
	caller, err := demo.NewCaller()
	if err != nil {
		return err
	}
	server := jsoncall.HttpServer{Caller: caller,
		ReceiverFunc: t.Receiver,
		Port:         t.Port,
	}
	return server.Run()
}

func main() {
	var cmd command.SimpleCommand
	var serverOps Server
	cmd.Flags(&serverOps).RunFunc(serverOps.Run)
	command.Main(&cmd)
}
