package main

import (
	"fmt"

	"melato.org/jsoncall/example/server"
)

func main() {
	err := server.Generate()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
