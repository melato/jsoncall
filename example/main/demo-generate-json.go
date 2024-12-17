// Copyright 2024 Alex Athanasopoulos.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"

	"melato.org/jsoncall/example"
	"melato.org/jsoncall/generate"
)

func GenerateJson() error {
	var api *example.Demo
	return generate.UpdateDescriptor(api, "../demo.json")

}

func main() {
	err := GenerateJson()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
