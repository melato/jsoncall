package cli

import (
	"melato.org/jsoncall"
)

func TraceVariables() map[string]*bool {
	return map[string]*bool{
		"data":  &jsoncall.TraceData,
		"calls": &jsoncall.TraceCalls,
		"init":  &jsoncall.TraceInit,
	}
}
