package main

import (
	_ "embed"
	"encoding/json"
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
	math  jsoncall.Client
}

func (t *ClientOps) Init() error {
	t.Url = "http://localhost:8080/"
	return nil
}

func (t *ClientOps) newMathClient() (jsoncall.Client, error) {
	var math *example.Math
	caller, err := jsoncall.NewCaller(math, nil)
	if err != nil {
		return nil, err
	}
	return caller.NewHttpClient(t.Url + "math/"), nil
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
	t.math, err = t.newMathClient()
	if err != nil {
		return err
	}
	return err
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

func (t *ClientOps) Div(a, b int32) error {
	return t.Math("Div", a, b)
}

func (t *ClientOps) Math(method string, a, b int32) error {
	return t.callMath(method, a, b)
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
	cmd.Command("div").RunFunc(ops.Div)
	cmd.Command("ping").RunFunc(ops.Ping)
	cmd.Command("error").RunFunc(ops.Error)
	cmd.Command("math").RunFunc(ops.Math)
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
