package client

import (
  "melato.org/jsoncall"
)

// GeneratedClient - Generated client for demo.Demo
type GeneratedClient struct {
  Client   *jsoncall.Client
}

type rHello struct {
  P1 string
}

func (t *GeneratedClient) Hello() (string, error) {
  var out rHello
  err := t.Client.CallV(&out, "Hello")
  return  out.P1, err
}

type rPing struct {
}

func (t *GeneratedClient) Ping() error{
  var out rPing
  err := t.Client.CallV(&out, "Ping")
  return  err
}

type rWait struct {
}

func (t *GeneratedClient) Wait(p0 int) error{
  var out rWait
  err := t.Client.CallV(&out, "Wait", p0)
  return  err
}
