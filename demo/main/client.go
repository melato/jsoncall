package main

import (
	"fmt"

	"melato.org/command"
	"melato.org/jsoncall"
	"melato.org/jsoncall/demo"
	"melato.org/jsoncall/demo/client"
)

type ClientOps struct {
	Url   string
	Trace bool
	demo  demo.Demo
}

func (t *ClientOps) Init() error {
	t.Url = "http://localhost:8080/"
	return nil
}

func (t *ClientOps) Configured() error {
	if t.Trace {
		jsoncall.TraceInit = true
		jsoncall.TraceCalls = true
		jsoncall.TraceData = true
	}
	var err error
	t.demo, err = client.NewClient(t.Url)
	return err
}

func (t *ClientOps) Ping() error {
	return t.demo.Ping()
}

func (t *ClientOps) Hello() error {
	s, err := t.demo.Hello()
	if err != nil {
		return err
	}
	fmt.Println(s)
	return nil
}

func (t *ClientOps) Wait(seconds int) error {
	return t.demo.Wait(seconds)
}

func Command() *command.SimpleCommand {
	cmd := &command.SimpleCommand{}
	var ops ClientOps
	cmd.Flags(&ops)
	cmd.Command("ping").RunFunc(ops.Ping)
	cmd.Command("hello").RunFunc(ops.Hello)
	cmd.Command("wait").RunFunc(ops.Wait)
	return cmd
}

func main() {
	cmd := Command()
	command.Main(cmd)
}
