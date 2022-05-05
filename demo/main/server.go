package main

import (
	"net/http"

	"melato.org/command"
	"melato.org/jsoncall"
	"melato.org/jsoncall/demo"
	"melato.org/jsoncall/generate"
)

type GenerateOp struct {
	generate.Generator
}

func (t *GenerateOp) Init() error {
	g := &t.Generator
	g.Init()
	g.Package = "client"
	g.Type = "GeneratedClient"
	g.OutputFile = "../client/generated_client.go"
	return nil
}

func (t *GenerateOp) Generate() error {
	var v *demo.Demo
	return t.Output(t.GenerateP(v))

}

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
	var api *demo.DemoImpl
	server, err := jsoncall.NewServer(api)
	if err != nil {
		return err
	}
	server.ReceiverFunc = t.Receiver
	server.Port = t.Port
	return server.Run()
}

func main() {
	var cmd command.SimpleCommand
	var generateOp GenerateOp
	cmd.Command("generate").Flags(&generateOp).RunFunc(generateOp.Generate)

	var serverOps Server
	cmd.Command("listen").Flags(&serverOps).RunFunc(serverOps.Run)
	command.Main(&cmd)
}
