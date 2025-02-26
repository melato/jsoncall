// Copyright 2024 Alex Athanasopoulos.
// SPDX-License-Identifier: Apache-2.0

package jsoncall

import (
	"encoding/json"
	"fmt"
)

// MethodDescriptor provides external names to use when marshalling/unmarshalling a method, its inputs, and its outputs
// Providing an API descriptor ensures that the Go method names and input/output order can change without affecting any clients.
type MethodDescriptor struct {
	// Method is the name of the Go Method
	Method string `json:"method"`
	// Path is the name of the method used in transit, typically part of a URL
	Path string `json:"path"`
	// In - the names of the inputs
	In []string `json:"in,omitempty"`
	// Out - the names of the outputs
	Out []string `json:"out,omitempty"`
}

// ApiDescriptor has method descriptors
type ApiDescriptor []*MethodDescriptor

func (t *MethodDescriptor) InputMap(args []interface{}) (map[string]interface{}, error) {
	if len(args) != len(t.In) {
		return nil, fmt.Errorf("wrong # of arguments for %s: %d/%d", t.Method, len(args), len(t.In))
	}
	m := make(map[string]interface{})
	for i, arg := range args {
		if arg != nil {
			m[t.In[i]] = arg
		}
	}
	return m, nil
}

func (t *MethodDescriptor) MarshalInputs(args ...interface{}) ([]byte, error) {
	m, err := t.InputMap(args)
	if err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

// Merge overrides a with b
// Typically a would contain the default method descriptors, and b would contain user-defined method descriptors
// If a method exists in both a and b, it uses the specification from b
// If the two specifications of a method have different number of inputs or outputs, it returns an error
// If a method exists in b but not in a, it is ignored
func (a ApiDescriptor) Merge(b ApiDescriptor) error {
	bmap := make(map[string]*MethodDescriptor)
	for _, m := range b {
		bmap[m.Method] = m
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

func ParseApiDescriptor(data []byte) (ApiDescriptor, error) {
	var desc ApiDescriptor
	err := json.Unmarshal(data, &desc)
	if err != nil {
		return nil, err
	}
	return desc, nil
}
