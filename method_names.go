package jsoncall

import (
	"encoding/json"
	"fmt"
	"reflect"
)

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
