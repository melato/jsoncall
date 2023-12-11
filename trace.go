package jsoncall

import (
	_ "embed"
)

//go:embed trace.yaml
var TraceDescriptions []byte

func TraceFuncs() map[string]func(bool) {
	return map[string]func(bool){
		"data":  func(b bool) { TraceData = b },
		"calls": func(b bool) { TraceCalls = b },
		"init":  func(b bool) { TraceInit = b },
	}
}

type Trace struct {
}

func (t *Trace) Funcs() map[string]func(bool) {
	return TraceFuncs()
}
func (t *Trace) Descriptions() []byte {
	return TraceDescriptions
}
