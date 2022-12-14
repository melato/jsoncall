package example

import (
	_ "embed"

	"melato.org/jsoncall"
)

//go:embed demo.json
var demoNames []byte

func NewCaller() (*jsoncall.Caller, error) {
	var api *Demo
	return jsoncall.NewCaller(api, demoNames)
}
