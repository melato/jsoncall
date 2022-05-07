package client

import (
  "melato.org/jsoncall"
)

// GeneratedClient - Generated client for demo.Demo
type GeneratedClient struct {
  Client   *jsoncall.Client
}

type rAdd struct {
  P1 int32 `json:"sum"`
}

func (t *GeneratedClient) Add(p1 int32, p2 int32) (int32, error) {
  var out rAdd
  err := t.Client.CallV(&out, "Add", p1, p2)
  return  out.P1, err
}

type rHello struct {
  P1 string `json:"s"`
}

func (t *GeneratedClient) Hello() (string, error) {
  var out rHello
  err := t.Client.CallV(&out, "Hello")
  return  out.P1, err
}

func (t *GeneratedClient) Ping() error{
  err := t.Client.CallV(nil, "Ping")
  return  err
}

func (t *GeneratedClient) Wait(p1 int) error{
  err := t.Client.CallV(nil, "Wait", p1)
  return  err
}
