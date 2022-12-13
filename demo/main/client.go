package main

import (
	"encoding/json"
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
	math  jsoncall.Client
}

func (t *ClientOps) Init() error {
	t.Url = "http://localhost:8080/"
	return nil
}

func (t *ClientOps) newMathClient() (jsoncall.Client, error) {
	var math *demo.Math
	caller, err := jsoncall.NewCaller(math, nil)
	if err != nil {
		return nil, err
	}
	caller.Prefix = "math"
	return &jsoncall.HttpClient{Caller: caller, Url: t.Url}, nil
}

func (t *ClientOps) Configured() error {
	if t.Trace {
		jsoncall.TraceInit = true
		jsoncall.TraceCalls = true
		jsoncall.TraceData = true
	}
	var err error
	t.demo, err = client.NewDemoClient(t.Url)
	if err != nil {
		return err
	}
	t.math, err = t.newMathClient()
	if err != nil {
		return err
	}
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

func (t *ClientOps) callMath(method string, args ...any) error {
	var result map[string]any
	err := t.math.Call(&result, method, args...)
	if err != nil {
		return err
	}
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}
	fmt.Printf("%s\nj", string(data))
	return nil
}

func (t *ClientOps) Add(a, b int32) error {
	return t.callMath("Add", a, b)
}

func (t *ClientOps) Div(a, b int32) error {
	return t.callMath("Div", a, b)
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
	cmd.Command("div").RunFunc(ops.Div)
	cmd.Command("wait").RunFunc(ops.Wait)
	cmd.Command("reflect").RunFunc(ops.Reflect)
	return cmd
}

func main() {
	cmd := Command()
	command.Main(cmd)
}
