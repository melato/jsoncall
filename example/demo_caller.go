// Copyright 2024 Alex Athanasopoulos.
// SPDX-License-Identifier: Apache-2.0

package example

import (
	_ "embed"

	"melato.org/jsoncall"
)

//go:embed demo.json
var demoNames []byte

// NewDemoCaller combines the Demo interface with the .json naming file
func NewDemoCaller() (*jsoncall.Caller, error) {
	var api *Demo
	return jsoncall.NewCaller(api, demoNames)
}
