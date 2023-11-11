package jsoncall

import (
	"fmt"
)

func traceFunc(v *bool) func(bool) { return func(b bool) { *v = b } }

func TraceVariables() map[string]*bool {
	fmt.Printf("TraceVariables is deprecated.  Use TraceFuncs")
	return map[string]*bool{
		"data":  &TraceData,
		"calls": &TraceCalls,
		"init":  &TraceInit,
	}
}

func traceFunc(v *bool) func(bool) { return func(b bool) { *v = b } }

func TraceFuncs() map[string]func(bool) {
	funcs := make(map[string]func(bool))
	for name, v := range TraceVariables() {
		funcs[name] = traceFunc(v)
	}
	return funcs

}
