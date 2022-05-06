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
	pType := reflect.TypeOf(proto)
	if pType.Kind() == reflect.Pointer {
		eType := pType.Elem()
		if eType.Kind() == reflect.Interface {
			return NewJsonCaller(eType)
		}
	}
	return NewJsonCaller(pType)
}

func NewJsonCaller(api reflect.Type) (*Caller, error) {
	var c Caller
	if api == nil {
		return nil, fmt.Errorf("nil api type")
	}
	var hasReceiver bool
	switch api.Kind() {
	case reflect.Pointer:
		hasReceiver = true
	case reflect.Interface:
		hasReceiver = false
	default:
		return nil, fmt.Errorf("unsupported api (%v) kind: %v", api, api.Kind())
	}
	c.Api = api
	n := api.NumMethod()
	if TraceInit {
		fmt.Printf("api type: %v methods: %d\n", api, n)
	}
	c.Methods = make(map[string]*Method, n)
	for i := 0; i < n; i++ {
		method := api.Method(i)
		m := newMethod(method, hasReceiver)
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
	rType := reflect.TypeOf(receiver)
	if TraceCalls {
		fmt.Printf("%v.%s(%s)\n", rType, name, string(jsonIn))
	}
	m := t.Methods[name]
	if m == nil {
		return nil, ErrNoSuchMethod, Errorf("unknown api method: %v.%s", t.Api, name)
	}
	method, exists := rType.MethodByName(name)
	if !exists {
		return nil, ErrNoSuchMethod, Errorf("unknown receiver method: %v.%s", rType, name)
	}

	in, err := m.unmarshalInputs(receiver, jsonIn)
	if err != nil {
		return nil, ErrMarshal, &Error{err.Error()}
	}
	out := method.Func.Call(in)
	outputData, err := m.marshalOutputs(out)
	if err != nil {
		return nil, ErrMarshal, err
	}
	if TraceCalls {
		fmt.Printf("result: %s\n", string(outputData))
	}
	return outputData, ErrNone, nil
}
