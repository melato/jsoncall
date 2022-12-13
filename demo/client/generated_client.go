package client

import (
  "melato.org/jsoncall"
)

// GeneratedClient - Generated client for demo.Demo
type GeneratedClient struct {
  Client   jsoncall.Client
}

type rHello struct {
  P1 string `json:"s"`
}

func (t *GeneratedClient) Hello() (string, error) {
  var out rHello
  err := t.Client.Call(&out, "Hello")
  return  out.P1, err
}

func (t *GeneratedClient) Ping() error{
  err := t.Client.Call(nil, "Ping")
  return  err
}

type rSeconds struct {
  P1 int `json:"seconds"`
}

func (t *GeneratedClient) Seconds(p1 int, p2 int, p3 int) (int, error) {
  var out rSeconds
  err := t.Client.Call(&out, "Seconds", p1, p2, p3)
  return  out.P1, err
}

func (t *GeneratedClient) Wait(p1 int) error{
  err := t.Client.Call(nil, "Wait", p1)
  return  err
}
