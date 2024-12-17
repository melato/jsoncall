// Copyright 2024 Alex Athanasopoulos.
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"testing"

	"melato.org/jsoncall/example"
)

func TestDemoInterface(t *testing.T) {
	var d example.Demo
	var impl *DemoImpl
	d = impl
	_ = d
}
