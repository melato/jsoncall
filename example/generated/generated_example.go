// Generated code for main.ExampleInterface
package generated

import (
  "melato.org/jsoncall"
)

func NewExampleInterface(h *jsoncall.HttpClient) *ExampleInterface {
  return &ExampleInterface{h}
}

// ExampleInterface - Generated client for main.ExampleInterface
type ExampleInterface struct {
  Client   jsoncall.Client
}

type rA struct {
  P1 string `json:"result"`
}

func (t *ExampleInterface) A(p1 string, p2 int) (string, error) {
  var out rA
  err := t.Client.Call(&out, "A", p1, p2)
  return  out.P1, err
}

type rB struct {
  P1 string `json:"result"`
}

func (t *ExampleInterface) B() string{
  var out rB
  t.Client.Call(&out, "B")
  return  out.P1
}
