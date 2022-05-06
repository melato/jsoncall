package main

import (
	"fmt"
	"reflect"

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

func (t *ClientOps) Add(a, b int32) error {
	x, err := t.demo.Add(a, b)
	if err != nil {
		return err
	}
	fmt.Println(x)
	return nil
}

func (t *ClientOps) Wait(seconds int) error {
	return t.demo.Wait(seconds)
}

func showMethod(rType reflect.Type, name string) {
	m, exists := rType.MethodByName(name)
	if exists {
		fmt.Printf("%v(%v).%s in: %d\n", rType, rType.Kind(), name, m.Type.NumIn())
	} else {
		fmt.Printf("no such method %v.%s\n", rType, name)
	}
}

func (t *ClientOps) Reflect(name string) {
	var v *demo.Demo
	showMethod(reflect.TypeOf(v).Elem(), "Add")
	showMethod(reflect.TypeOf(t.demo), "Add")
}

func Command() *command.SimpleCommand {
	cmd := &command.SimpleCommand{}
	var ops ClientOps
	cmd.Flags(&ops)
	cmd.Command("ping").RunFunc(ops.Ping)
	cmd.Command("hello").RunFunc(ops.Hello)
	cmd.Command("add").RunFunc(ops.Add)
	cmd.Command("wait").RunFunc(ops.Wait)
	cmd.Command("reflect").RunFunc(ops.Reflect)
	return cmd
}

func main() {
	cmd := Command()
	command.Main(cmd)
}
