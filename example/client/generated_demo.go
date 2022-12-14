package client

import (
  "melato.org/jsoncall"
  "melato.org/jsoncall/example"
)

// DemoClient - Generated client for example.Demo
type DemoClient struct {
  Client   jsoncall.Client
}

type rRepeat struct {
  P1 string `json:"result"`
}

func (t *DemoClient) Repeat(p1 string, p2 int) (string, error) {
  var out rRepeat
  err := t.Client.Call(&out, "Repeat", p1, p2)
  return  out.P1, err
}

type rTime struct {
  P1 int `json:"hour"`
  P2 int `json:"minute"`
  P3 int `json:"second"`
}

func (t *DemoClient) Time() (int, int, int) {
  var out rTime
  t.Client.Call(&out, "Time")
  return  out.P1, out.P2, out.P3
}

type rTimePointer struct {
  P1 *example.Time `json:"result"`
}

func (t *DemoClient) TimePointer() *example.Time{
  var out rTimePointer
  t.Client.Call(&out, "TimePointer")
  return  out.P1
}

type rTimeStruct struct {
  P1 example.Time `json:"result"`
}

func (t *DemoClient) TimeStruct() (example.Time, error) {
  var out rTimeStruct
  err := t.Client.Call(&out, "TimeStruct")
  return  out.P1, err
}

func (t *DemoClient) Wait(p1 int) error{
  err := t.Client.Call(nil, "Wait", p1)
  return  err
}
