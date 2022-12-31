// Generated code for main.ExampleInterface
package generated

import (
  "melato.org/jsoncall"
)

func NewExampleClient(h *jsoncall.HttpClient) *ExampleInterface {
  return &ExampleInterface{h}
}

// ExampleInterface - Generated client for main.ExampleInterface
type ExampleInterface struct {
  Client   jsoncall.Client
}

type rA struct {
  P1 string `json:"result"`
  P2 *jsoncall.Error `json:"error"`
}

func (t *ExampleInterface) A(p1 string, p2 int) (string, error) {
  var out rA
  err := t.Client.Call(&out, "A", p1, p2)
  if err != nil {
    return out.P1, err
  }
  return out.P1, out.P2.ToError()
}

type rB struct {
  P1 string `json:"result"`
}

func (t *ExampleInterface) B() string{
  var out rB
  t.Client.Call(&out, "B")
  return out.P1
}
