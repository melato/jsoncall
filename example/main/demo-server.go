// Copyright 2024 Alex Athanasopoulos.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"net/http"

	"melato.org/command"
	"melato.org/jsoncall"
	"melato.org/jsoncall/example"
	"melato.org/jsoncall/example/server"
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
		jsoncall.TraceData = true
	}
	return nil
}

func (t *Server) DemoReceiver(w http.ResponseWriter, r *http.Request) interface{} {
	return &server.DemoImpl{}
}

func (t *Server) ListenAndServe(handler http.Handler) error {
	addr := fmt.Sprintf(":%d", t.Port)
	fmt.Printf("starting server at %s\n", addr)
	return http.ListenAndServe(addr, handler)
}

// RunMux demostrates how to use http.ServeMux to implement a server that handles multiple interfaces
func (t *Server) Run() error {
	demoCaller, err := example.NewDemoCaller()
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.Handle("/demo/", demoCaller.NewHttpHandler(t.DemoReceiver))

	return t.ListenAndServe(mux)
}

func main() {
	var cmd command.SimpleCommand
	var serverOps Server
	cmd.Flags(&serverOps)
	cmd.RunFunc(serverOps.Run)
	command.Main(&cmd)
}
