package client

import (
  "melato.org/jsoncall"
)

// DemoClient - Generated client for example.Demo
type DemoClient struct {
  Client   jsoncall.Client
}

type rHello struct {
  P1 string `json:"s"`
}

func (t *DemoClient) Hello() (string, error) {
  var out rHello
  err := t.Client.Call(&out, "Hello")
  return  out.P1, err
}

func (t *DemoClient) Nop() {
  t.Client.Call(nil, "Nop")
}

func (t *DemoClient) Ping() error{
  err := t.Client.Call(nil, "Ping")
  return  err
}

type rSeconds struct {
  P1 int `json:"result"`
}

func (t *DemoClient) Seconds(p1 int, p2 int, p3 int) (int, error) {
  var out rSeconds
  err := t.Client.Call(&out, "Seconds", p1, p2, p3)
  return  out.P1, err
}

func (t *DemoClient) Wait(p1 int) error{
  err := t.Client.Call(nil, "Wait", p1)
  return  err
}
