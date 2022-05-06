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
	Names   []*MethodNames
}

func HasReceiver(api reflect.Type) bool {
	if api.Kind() == reflect.Interface {
		return false
	}
	return true
}

func (c *Caller) SetTypePointer(proto interface{}) error {
	pType := reflect.TypeOf(proto)
	switch pType.Kind() {
	case reflect.Interface:
		return fmt.Errorf("cannot use interface{} prototype.  must be pointer or slice")
	case reflect.Pointer:
		eType := pType.Elem()
		if eType.Kind() == reflect.Interface {
			return c.SetType(eType)
		}
	}
	return c.SetType(pType)
}

func (c *Caller) SetType(api reflect.Type) error {
	if api == nil {
		return fmt.Errorf("nil api type")
	}
	hasReceiver := HasReceiver(api)
	switch api.Kind() {
	case reflect.Pointer:
	case reflect.Interface:
	default:
		return fmt.Errorf("unsupported api (%v) kind: %v", api, api.Kind())
	}
	c.Api = api
	n := api.NumMethod()
	if TraceInit {
		fmt.Printf("api type: %v methods: %d\n", api, n)
	}
	namesMap := make(map[string]*MethodNames)
	for _, m := range c.Names {
		namesMap[m.Method] = m
	}
	c.Methods = make(map[string]*Method, n)
	for i := 0; i < n; i++ {
		method := api.Method(i)
		m := newMethod(method, hasReceiver, namesMap[method.Name])
		c.Methods[method.Name] = m
	}
	return nil
}

// Call - call a receiver method.
// Unmarshal the method parameters from JSON and marshal the outputs to JSON
// arguments and return values are marshalled as a map, with keys "P1", "P2", ...
// If the last output is of type error, it is unmashalled as *Error
// If there is an error in unmarshalling/marshalling, return nil, *Error
func (t *Caller) Call(m *Method, receiver interface{}, jsonIn []byte) ([]byte, ErrorCode, error) {
	rType := reflect.TypeOf(receiver)
	name := m.Names.Method
	if TraceCalls {
		fmt.Printf("%v.%s(%s)\n", rType, name, string(jsonIn))
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
