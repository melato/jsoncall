package client

import (
  "melato.org/jsoncall"
)

// GeneratedClient - Generated client for demo.Demo
type GeneratedClient struct {
  Client   *jsoncall.Client
}

func (t *GeneratedClient) Hello() (string, error) {
  result := t.Client.Call("Hello")
  var x0 string = result[0].(string)
  var x1 error
  if result[1] != nil {
	 x1 = result[1].(error)
  }
  return x0, x1
}

func (t *GeneratedClient) Ping() error{
  result := t.Client.Call("Ping")
  var x0 error
  if result[0] != nil {
	 x0 = result[0].(error)
  }
  return x0
}

func (t *GeneratedClient) Wait(p0 int) error{
  result := t.Client.Call("Wait")
  var x0 error
  if result[0] != nil {
	 x0 = result[0].(error)
  }
  return x0
}
