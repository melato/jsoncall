// Package client - provides an implementation of demo.Demo that calls a web service.
// This is a separate package, so the server can compile before generating generated_client.go
package client

import (
	"melato.org/jsoncall"
	"melato.org/jsoncall/demo"
)

func NewClient(url string) (demo.Demo, error) {
	c := &GeneratedClient{}
	var api *demo.Demo
	var err error
	c.Client, err = jsoncall.NewClientP(api)
	if err != nil {
		return nil, err
	}
	c.Client.Url = url
	return c, nil
}
