package jsoncall

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// MethodNames provides external names to use when marshalling/unmarshalling a method, its inputs, and its outputs
type MethodNames struct {
	// Method is the name of the Go Method
	Method string `json:"method"`
	// Path is the name of the method used in transit, typically part of a URL
	Path string `json:"path"`
	// In - the names of the inputs
	In []string `json:"in"`
	// Out - the names of the outputs
	Out []string `json:"out"`
}

// Names provides external names for methods, their inputs and their outputs
type Names []*MethodNames

func (t *MethodNames) MarshalInputs(args ...interface{}) ([]byte, error) {
	if len(args) != len(t.In) {
		return nil, fmt.Errorf("wrong # of arguments for %s: %d/%d", t.Method, len(args), len(t.In))
	}
	m := make(map[string]interface{})
	for i, arg := range args {
		if arg != nil {
			m[t.In[i]] = arg
		}
	}
	return json.Marshal(m)
}

func DefaultMethodNames(method reflect.Method, hasReceiver bool) *MethodNames {
	var m MethodNames
	m.Method = method.Name
	m.Path = method.Name
	numIn := method.Type.NumIn()
	if hasReceiver {
		numIn--
	}
	if numIn > 0 {
		m.In = make([]string, numIn)
		for i := 0; i < numIn; i++ {
			m.In[i] = fmt.Sprintf("P%d", i+1)
		}
	}
	numOut := method.Type.NumOut()
	if numOut > 0 {
		m.Out = make([]string, numOut)
		for i := 0; i < numOut; i++ {
			m.Out[i] = fmt.Sprintf("P%d", i+1)
		}
	}
	return &m
}

// Merge overrides a with b
// Typically a would be the default method names, and b would be user-defined
// If a method exists in both a and b, it uses the specification from b
// If the two specifications of a method have different number of inputs or outputs, it returns an error
// If a method exists in b but not in a, it is ignored
func (a Names) Merge(b Names) error {
	bmap := make(map[string]*MethodNames)
	for _, names := range b {
		bmap[names.Method] = names
	}
	for _, x := range a {
		y, exists := bmap[x.Method]
		if exists {
			if len(x.In) != len(y.In) {
				return fmt.Errorf("the inputs of method %s have changed.  Please update them.", x.Method)
			}
			if len(x.Out) != len(y.Out) {
				return fmt.Errorf("the outputs of method %s have changed.  Please update them.", x.Method)
			}
			x.In = y.In
			x.Out = y.Out
			x.Path = y.Path
		}
	}
	return nil
}

func ParseNames(data []byte) (Names, error) {
	var names Names
	err := json.Unmarshal(data, &names)
	if err != nil {
		return nil, err
	}
	return names, nil
}
