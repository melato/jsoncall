package cli

import (
	"melato.org/jsoncall"
)

func TraceMap() map[string]*bool {
	return map[string]*bool{
		"data":  &jsoncall.TraceData,
		"calls": &jsoncall.TraceCalls,
		"init":  &jsoncall.TraceInit,
	}
}
