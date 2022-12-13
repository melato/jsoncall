package jsoncall

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Method struct {
	Desc *MethodDescriptor
	// InType is a Struct type that containers one field for each method input
	// It is used by the server to unmarshal the input arguments
	InType    reflect.Type
	OutErrors []bool
}

func newMethod(method reflect.Method, hasReceiver bool, desc *MethodDescriptor) (*Method, error) {
	var m Method
	numIn := method.Type.NumIn()
	numOut := method.Type.NumOut()
	var offset int
	if hasReceiver {
		offset = 1
		numIn--
	}
	if numOut > 0 {
		var errp *error
		errorType := reflect.TypeOf(errp).Elem()
		m.OutErrors = make([]bool, numOut)
		for i := 0; i < numOut; i++ {
			if method.Type.Out(i) == errorType {
				m.OutErrors[i] = true
			}
		}
	}
	if desc == nil {
		desc = DefaultMethodDescriptor(method.Name, numIn, m.OutErrors)
	}
	if numIn != len(desc.In) {
		return nil, fmt.Errorf("method %s has %d inputs, but its descriptor has %d", method.Name, numIn, len(desc.In))
	}
	if numOut != len(desc.Out) {
		return nil, fmt.Errorf("method %s has %d outputs, but its descriptor has %d", method.Name, numOut, len(desc.Out))
	}
	m.Desc = desc
	if TraceInit {
		fmt.Printf("method: %s in: %d out: %d\n", method.Name, numIn, numOut)
	}
	if numIn > 0 {
		fields := make([]reflect.StructField, numIn)
		for i := 0; i < numIn; i++ {
			var field reflect.StructField
			field.Name = fmt.Sprintf("P%d", i+1)
			field.Tag = reflect.StructTag(fmt.Sprintf(`json:"%s"`, desc.In[i]))
			field.Type = method.Type.In(offset + i)
			fields[i] = field
			if TraceInit {
				fmt.Printf("%s field[%d]: %s (%v)\n", method.Name, i, field.Name, field.Type)
			}
		}
		m.InType = reflect.StructOf(fields)
	}
	return &m, nil
}

func (t *Method) MarshalInputs(args ...interface{}) ([]byte, error) {
	return t.Desc.MarshalInputs(args...)
}

func (m *Method) unmarshalInputs(receiver interface{}, data []byte) ([]reflect.Value, error) {
	numIn := len(m.Desc.In)
	in := make([]reflect.Value, 1+numIn)
	in[0] = reflect.ValueOf(receiver)
	if numIn > 0 {
		a := reflect.New(m.InType)
		v := a.Interface()
		err := json.Unmarshal(data, &v)
		if err != nil {
			return nil, err
		}
		a = a.Elem()
		for i := 0; i < numIn; i++ {
			in[1+i] = a.Field(i)
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
func (m *Method) marshalOutputs(out []reflect.Value) ([]byte, ErrorCode, error) {
	outMap := make(map[string]interface{})
	for i, x := range out {
		v := x.Interface()
		if m.OutErrors[i] && !x.IsNil() {
			e := v.(error)
			return nil, ErrUser, e
		}
		outMap[m.Desc.Out[i]] = v
	}
	data, err := json.Marshal(outMap)
	if err != nil {
		return nil, ErrMarshal, err
	}
	if TraceCalls {
		fmt.Printf("result: %s\n", string(data))
	}
	return data, ErrNone, nil
}

// Call - call a receiver method.
// Unmarshal the method parameters from JSON and marshal the outputs to JSON
// arguments and return values are marshalled as a map
// If there is an error in unmarshalling/marshalling, return nil, *Error
// If any output is of type error and is not nil, Call() returns nil, ErrUser, and the error
func (m *Method) Call(receiver interface{}, jsonIn []byte) ([]byte, ErrorCode, error) {
	rType := reflect.TypeOf(receiver)
	name := m.Desc.Method
	if TraceCalls {
		fmt.Printf("%v.%s(%s)\n", rType, name, string(jsonIn))
	}
	method, exists := rType.MethodByName(name)
	if !exists {
		return nil, ErrNoSuchMethod, Errorf("unknown receiver method: %v.%s", rType, name)
	}

	in, err := m.unmarshalInputs(receiver, jsonIn)
	if err != nil {
		return nil, ErrMarshal, err
	}
	out := method.Func.Call(in)
	return m.marshalOutputs(out)
}
