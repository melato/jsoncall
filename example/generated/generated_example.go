package generated

import (
  "melato.org/jsoncall"
)

// ExampleClient - Generated client for main.ExampleInterface
type ExampleClient struct {
  Client   jsoncall.Client
}

type rA struct {
  P1 string `json:"result"`
}

func (t *ExampleClient) A(p1 string, p2 int) (string, error) {
  var out rA
  err := t.Client.Call(&out, "A", p1, p2)
  return  out.P1, err
}
