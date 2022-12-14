// Package client - provides an implementation of demo.Demo that calls a web service.
// This is a separate package, so that other packages can compile without the generated files.
package client

import (
	"melato.org/jsoncall"
	"melato.org/jsoncall/example"
)

func NewDemoClient(url string) (example.Demo, error) {
	caller, err := example.NewCaller()
	if err != nil {
		return nil, err
	}
	caller.Prefix = "demo"
	c := &jsoncall.HttpClient{Caller: caller, Url: url}
	return &DemoClient{c}, nil
}
