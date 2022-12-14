package main

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"melato.org/command/usage"

	"melato.org/command"
	"melato.org/jsoncall"
	"melato.org/jsoncall/demo"
	"melato.org/jsoncall/demo/client"
)

//go:embed client.yaml
var usageData []byte

type ClientOps struct {
	Url   string
	Trace bool
	Json  bool `name:"json" usage:"print json response for Math calls"`
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

func (t *ClientOps) Nop() {
	t.demo.Nop()
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
	var response map[string]any
	err := t.math.Call(&response, method, args...)
	if err != nil {
		return err
	}
	if t.Json {
		data, err := json.Marshal(response)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", string(data))
	} else {
		fmt.Printf("%v\n", response["result"])
	}
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

func Command() *command.SimpleCommand {
	cmd := &command.SimpleCommand{}
	var ops ClientOps
	cmd.Flags(&ops)
	cmd.Command("ping").RunFunc(ops.Ping)
	cmd.Command("nop").RunFunc(ops.Nop)
	cmd.Command("hello").RunFunc(ops.Hello)
	cmd.Command("add").RunFunc(ops.Add)
	cmd.Command("div").RunFunc(ops.Div)
	cmd.Command("wait").RunFunc(ops.Wait)

	usage.Apply(cmd, usageData)
	return cmd
}

func main() {
	cmd := Command()
	command.Main(cmd)
}
