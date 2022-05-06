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

func (t *GeneratedClient) Add(p0 int32, p1 int32) (int32, error) {
  var out rAdd
  err := t.Client.CallV(&out, "Add", p0, p1)
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
