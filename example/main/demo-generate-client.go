// Copyright 2024 Alex Athanasopoulos.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"

	"melato.org/jsoncall/example"
	"melato.org/jsoncall/generate"
)

func GenerateClient() error {
	var g generate.Generator
	g.Init()
	g.Package = "generated"
	g.Type = "DemoClient"
	g.Func = "NewDemoClient"
	g.OutputFile = "../generated/generated_demo.go"
	g.Imports = []string{"melato.org/jsoncall/example"}
	caller, err := example.NewDemoCaller()
	if err != nil {
		return err
	}
	return g.Output(g.GenerateClient(caller))
}

func main() {
	err := GenerateClient()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
