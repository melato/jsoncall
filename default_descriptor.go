package jsoncall

import (
	"fmt"
	"reflect"
)

func DefaultMethodDescriptor(method reflect.Method, hasReceiver bool) *MethodDescriptor {
	numIn := method.Type.NumIn()
	if hasReceiver {
		numIn--
	}
	outErrors := getErrorOutputs(method)
	return defaultMethodDescriptor(method.Name, numIn, outErrors)
}

func defaultMethodDescriptor(name string, numIn int, outErrors []bool) *MethodDescriptor {
	var m MethodDescriptor
	m.Method = name
	m.Path = name
	if numIn > 0 {
		m.In = make([]string, numIn)
		for i := 0; i < numIn; i++ {
			m.In[i] = fmt.Sprintf("p%d", i+1)
		}
	}
	numOut := len(outErrors)
	if numOut > 0 {
		m.Out = make([]string, numOut)
		var numErrors int
		for _, isError := range outErrors {
			if isError {
				numErrors++
			}
		}
		for i, isError := range outErrors {
			var outName string
			if isError {
				if numErrors == 1 {
					outName = "error"
				} else {
					outName = fmt.Sprintf("e%d", i+1)
				}
			} else if numOut-numErrors == 1 {
				outName = "result"
			} else {
				outName = fmt.Sprintf("r%d", i+1)
			}
			m.Out[i] = outName
		}
	}
	return &m
}
