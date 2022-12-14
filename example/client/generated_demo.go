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

type rSubstring struct {
  P1 string `json:"result"`
}

func (t *DemoClient) Substring(p1 string, p2 int, p3 int) (string, error) {
  var out rSubstring
  err := t.Client.Call(&out, "Substring", p1, p2, p3)
  return  out.P1, err
}

type rTime struct {
  P1 int `json:"hour"`
  P2 int `json:"minute"`
  P3 int `json:"second"`
}

func (t *DemoClient) Time() (int, int, int, error) {
  var out rTime
  err := t.Client.Call(&out, "Time")
  return  out.P1, out.P2, out.P3, err
}

func (t *DemoClient) Wait(p1 int) error{
  err := t.Client.Call(nil, "Wait", p1)
  return  err
}
