package demo

import (
	_ "embed"
	"encoding/json"

	"melato.org/jsoncall"
)

//go:embed demo.json
var demoNames []byte

func NewCaller() (*jsoncall.Caller, error) {
	var api *Demo
	var c jsoncall.Caller
	err := json.Unmarshal(demoNames, &c.Names)
	if err != nil {
		return nil, err
	}
	err = c.SetTypePointer(api)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
