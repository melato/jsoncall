package example

import (
	_ "embed"

	"melato.org/jsoncall"
)

//go:embed demo.json
var demoNames []byte

func NewDemoCaller() (*jsoncall.Caller, error) {
	var api *Demo
	return jsoncall.NewCaller(api, demoNames)
}
