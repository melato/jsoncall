// Package client - provides an implementation of demo.Demo that calls a web service.
// This is a separate package, so that other packages can compile without generated_client.go
package client

import (
	"melato.org/jsoncall"
	"melato.org/jsoncall/demo"
)

func NewClient(url string) (demo.Demo, error) {
	caller, err := demo.NewCaller()
	if err != nil {
		return nil, err
	}
	c := &jsoncall.HttpClient{Caller: caller, Url: url}
	return &GeneratedClient{c}, nil
}
