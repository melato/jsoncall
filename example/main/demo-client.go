// Copyright 2024 Alex Athanasopoulos.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	_ "embed"
	"fmt"

	"melato.org/command/usage"

	"melato.org/command"
	"melato.org/jsoncall"
	"melato.org/jsoncall/example"
	"melato.org/jsoncall/example/generated"
)

//go:embed demo-client.yaml
var usageData []byte

type ClientOps struct {
	Url   string
	Trace bool
	Json  bool `name:"json" usage:"print json response for Math calls"`
	demo  example.Demo
}

func (t *ClientOps) Init() error {
	t.Url = "http://localhost:8080/"
	return nil
}

func (t *ClientOps) NewDemoClient() (example.Demo, error) {
	caller, err := example.NewDemoCaller()
	if err != nil {
		return nil, err
	}
	c := caller.NewHttpClient(t.Url + "demo/")
	return generated.NewDemoClient(c), nil
}

func (t *ClientOps) Configured() error {
	if t.Trace {
		//jsoncall.TraceInit = true
		jsoncall.TraceCalls = true
		jsoncall.TraceData = true
	}
	var err error
	t.demo, err = t.NewDemoClient()
	if err != nil {
		return err
	}
	return err
}

func (t *ClientOps) Wait(seconds int) error {
	return t.demo.Wait(seconds)
}

func (t *ClientOps) Time() {
	hour, minute, second := t.demo.Time()
	fmt.Printf("%02d:%02d:%02d\n", hour, minute, second)
}

func (t *ClientOps) TimeStruct() error {
	r, err := t.demo.TimeStruct()
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", r)
	return nil
}

func (t *ClientOps) TimePointer() error {
	r := t.demo.TimePointer()
	fmt.Printf("%v\n", r)
	return nil
}

func (t *ClientOps) Repeat(s string, count int) error {
	list, err := t.demo.Repeat(s, count)
	if err != nil {
		return err
	}
	for _, s := range list {
		fmt.Printf("%s\n", s)
	}
	return nil
}

func (t *ClientOps) Ping() error {
	return t.demo.Ping()
}

func (t *ClientOps) Error() error {
	s, err := t.demo.Error()
	fmt.Printf("%s\n", s)
	return err
}

func Command() *command.SimpleCommand {
	cmd := &command.SimpleCommand{}
	var ops ClientOps
	cmd.Flags(&ops)
	cmd.Command("ping").RunFunc(ops.Ping)
	cmd.Command("error").RunFunc(ops.Error)
	cmd.Command("wait").RunFunc(ops.Wait)
	cmd.Command("time").RunFunc(ops.Time)
	cmd.Command("time-struct").RunFunc(ops.TimeStruct)
	cmd.Command("time-pointer").RunFunc(ops.TimePointer)
	cmd.Command("repeat").RunFunc(ops.Repeat)

	usage.Apply(cmd, usageData)
	return cmd
}

func main() {
	cmd := Command()
	command.Main(cmd)
}
