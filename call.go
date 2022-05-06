package jsoncall

import (
	"fmt"
	"reflect"
)

var TraceCalls bool
var TraceInit bool
var TraceDebug bool

type Caller struct {
	Api     reflect.Type
	Methods map[string]*Method
}

func NewJsonCallerP(proto interface{}) (*Caller, error) {
	return NewJsonCaller(reflect.TypeOf(proto))
}

func NewJsonCaller(api reflect.Type) (*Caller, error) {
	var c Caller
	c.Api = api
	if c.Api == nil {
		return nil, fmt.Errorf("nil api type")
	}
	n := api.NumMethod()
	if TraceInit {
		fmt.Printf("api type: %v methods: %d\n", c.Api, n)
	}
	c.Methods = make(map[string]*Method, n)
	for i := 0; i < n; i++ {
		method := api.Method(i)
		m := newMethod(method)
		c.Methods[method.Name] = m
	}
	return &c, nil
}

// Call - call the receiver method with the given name.
// Unmarshal the method parameters from JSON and marshal the outputs to JSON
// arguments and return values are marshalled as a map, with keys "P1", "P2", ...
// If the last output is of type error, it is unmashalled as *Error
// If there is an error in unmarshalling/marshalling, return nil, *Error
func (t *Caller) Call(name string, receiver interface{}, jsonIn []byte) ([]byte, ErrorCode, error) {
	if TraceCalls {
		fmt.Printf("%v.%s(%s)\n", reflect.TypeOf(receiver), name, string(jsonIn))
	}
	m := t.Methods[name]
	if m == nil {
		return nil, ErrNoSuchMethod, Errorf("unknown method: %v.%s", t.Api, name)
	}

	in, err := m.unmarshalInputs(receiver, jsonIn)
	if err != nil {
		return nil, ErrMarshal, &Error{err.Error()}
	}
	out := m.Method.Func.Call(in)
	outputData, err := m.marshalOutputs(out)
	if err != nil {
		return nil, ErrMarshal, err
	}
	if TraceCalls {
		fmt.Printf("result: %s\n", string(outputData))
	}
	return outputData, ErrNone, nil
}
