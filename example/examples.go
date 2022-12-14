package example

import (
	"net/http"

	"melato.org/jsoncall"
)

type ExampleApi interface {
	A(s string, d int) (string, error)
}

type Example struct {
}

func (t *Example) A(s string, d int) (string, error) {
	return "a", nil
}

func (t *Example) B() (string, error) {
	return "b", nil
}

func ExampleServer() error {
	var handler http.Handler
	handler, err := jsoncall.NewHttpHandler(&Example{}, nil)
	if err != nil {
		return err
	}
	return http.ListenAndServe(":8080", handler)
}

func ExampleServerWithInterface() error {
	var handler http.Handler
	var api *ExampleApi
	handler, err := jsoncall.NewHttpHandler(&Example{}, api)
	if err != nil {
		return err
	}
	return http.ListenAndServe(":8080", handler)
}
