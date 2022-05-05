package jsoncall

import (
	"encoding/json"
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

type Method struct {
	Method       reflect.Method
	InType       reflect.Type // a Struct type that contains the method's inputs
	OutType      reflect.Type // a Struct type that contains the method's outputs
	OutNames     []string
	OutErrors    []bool
	LastOutError int
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

func newMethod(method reflect.Method) *Method {
	if TraceInit {
		fmt.Printf("method: %s\n", method.Name)
	}
	var m Method
	m.Method = method
	numIn := method.Type.NumIn()
	if numIn > 1 {
		fields := make([]reflect.StructField, numIn-1)
		for i := 1; i < numIn; i++ {
			var field reflect.StructField
			field.Name = fmt.Sprintf("P%d", i)
			field.Type = method.Type.In(i)
			fields[i-1] = field
		}
		m.InType = reflect.StructOf(fields)
	}
	numOut := method.Type.NumOut()
	m.LastOutError = -1
	if numOut > 0 {
		var errp *error
		var Errp *Error
		errorType := reflect.TypeOf(errp).Elem()
		ErrorType := reflect.TypeOf(Errp)
		fields := make([]reflect.StructField, numOut)
		m.OutNames = make([]string, numOut)
		m.OutErrors = make([]bool, numOut)
		for i := 0; i < numOut; i++ {
			var field reflect.StructField
			field.Name = fmt.Sprintf("P%d", i+1)
			field.Type = method.Type.Out(i)
			if field.Type == errorType {
				field.Type = ErrorType
				m.OutErrors[i] = true
				m.LastOutError = i
			}
			fields[i] = field
			m.OutNames[i] = field.Name
		}
		m.OutType = reflect.StructOf(fields)
	}
	return &m
}

func Marshal(args ...interface{}) ([]byte, error) {
	m := make(map[string]interface{})
	for i, arg := range args {
		if arg != nil {
			m[fmt.Sprintf("P%d", i+1)] = arg
		}
	}
	return json.Marshal(m)
}

func (m *Method) unmarshalInputs(receiver interface{}, data []byte) ([]reflect.Value, error) {
	numIn := m.Method.Type.NumIn()
	in := make([]reflect.Value, numIn)
	in[0] = reflect.ValueOf(receiver)
	if numIn > 1 {
		a := reflect.New(m.InType)
		v := a.Interface()
		err := json.Unmarshal(data, &v)
		if err != nil {
			return nil, err
		}
		//a = reflect.ValueOf(v).Elem()
		a = a.Elem()
		for i := 1; i < numIn; i++ {
			in[i] = a.Field(i - 1)
		}
	}
	return in, nil
}

// marshalOutputs - convert outputs as returned from a method call to JSON data
// The outputs are marshalled to JSON as follows:
// If the last output is of type error, and it is not nil, return nil, *Error
// If the last output is of type error and it is nil, return the remaining outputs as follows:
// If there are no outputs, return nil
// If there is 1 output, return the JSON representation of that output
// If there are more than 1 outputs, return a JSON representation of a map where the keys are "P1", "P2", ...
// and the values are the outputs
// If there is an error in marshalling, return nil, *Error
func (m *Method) marshalOutputs(out []reflect.Value) ([]byte, error) {
	outMap := make(map[string]interface{})
	for i, x := range out {
		v := x.Interface()
		if m.OutErrors[i] && !x.IsNil() {
			e := v.(error)
			v = ToError(e)
		}
		outMap[m.OutNames[i]] = v
	}

	return json.Marshal(outMap)
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
