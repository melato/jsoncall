package demo

import (
	_ "embed"
	"encoding/json"

	"melato.org/jsoncall"
)

//go:embed demo.json
var demoNames []byte

func NewCaller() (*jsoncall.Caller, error) {
	var c jsoncall.Caller
	err := json.Unmarshal(demoNames, &c.Names)
	if err != nil {
		return nil, err
	}
	var api *Demo
	err = c.SetTypePointer(api)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
