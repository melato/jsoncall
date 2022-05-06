package jsoncall

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Method struct {
	Name      string
	InType    reflect.Type // a Struct type that contains the method's inputs
	InNames   []string
	OutNames  []string
	OutErrors []bool
}

func newMethod(method reflect.Method, hasReceiver bool) *Method {
	var m Method
	m.Name = method.Name
	numIn := method.Type.NumIn()
	var offset int
	if hasReceiver {
		offset = 1
		numIn--
	}
	if TraceInit {
		fmt.Printf("method: %s in: %d\n", method.Name, numIn)
	}
	if numIn > 0 {
		m.InNames = make([]string, numIn)
		fields := make([]reflect.StructField, numIn)
		for i := 0; i < numIn; i++ {
			var field reflect.StructField
			field.Name = fmt.Sprintf("P%d", i+1)
			field.Type = method.Type.In(offset + i)
			fields[i] = field
			m.InNames[i] = field.Name
			if TraceInit {
				fmt.Printf("%s field[%d]: %s (%v)\n", m.Name, i, field.Name, field.Type)
			}
		}
		m.InType = reflect.StructOf(fields)
	}
	numOut := method.Type.NumOut()
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
			}
			fields[i] = field
			m.OutNames[i] = field.Name
		}
	}
	return &m
}

func (t *Method) MarshalInputs(args ...interface{}) ([]byte, error) {
	if len(args) != len(t.InNames) {
		return nil, fmt.Errorf("wrong # of arguments for %s: %d/%d", t.Name, len(args), len(t.InNames))
	}
	m := make(map[string]interface{})
	for i, arg := range args {
		if arg != nil {
			m[t.InNames[i]] = arg
		}
	}
	return json.Marshal(m)
}

func (m *Method) unmarshalInputs(receiver interface{}, data []byte) ([]reflect.Value, error) {
	numIn := len(m.InNames)
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
