package example

import (
	"net/http"
	"testing"

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

func TestExample(t *testing.T) {
	var handler http.Handler
	var err error
	handler, err = jsoncall.NewHttpHandler(&Example{})
	if err != nil || handler == nil {
		t.Fail()
	}
}
