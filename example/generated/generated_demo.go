// Generated code for example.Demo
package generated

import (
  "melato.org/jsoncall"
  "melato.org/jsoncall/example"
)

func NewDemoClient(h *jsoncall.HttpClient) *DemoClient {
  return &DemoClient{h}
}

// DemoClient - Generated client for example.Demo
type DemoClient struct {
  Client   jsoncall.Client
}

type rError struct {
  P1 string `json:"result"`
  P2 *jsoncall.Error `json:"error"`
}

func (t *DemoClient) Error() (string, error) {
  var out rError
  err := t.Client.Call(&out, "Error")
  if err != nil {
    return out.P1, err
  }
  return out.P1, out.P2.ToError()
}

type rPing struct {
  P1 *jsoncall.Error `json:"error"`
}

func (t *DemoClient) Ping() error{
  var out rPing
  err := t.Client.Call(&out, "Ping")
  if err != nil {
    return err
  }
  return out.P1.ToError()
}

type rRepeat struct {
  P1 []string `json:"result"`
  P2 *jsoncall.Error `json:"error"`
}

func (t *DemoClient) Repeat(p1 string, p2 int) ([]string, error) {
  var out rRepeat
  err := t.Client.Call(&out, "Repeat", p1, p2)
  if err != nil {
    return out.P1, err
  }
  return out.P1, out.P2.ToError()
}

type rTime struct {
  P1 int `json:"hour"`
  P2 int `json:"minute"`
  P3 int `json:"second"`
}

func (t *DemoClient) Time() (int, int, int) {
  var out rTime
  t.Client.Call(&out, "Time")
  return out.P1, out.P2, out.P3
}

type rTimePointer struct {
  P1 *example.Time `json:"result"`
}

func (t *DemoClient) TimePointer() *example.Time{
  var out rTimePointer
  t.Client.Call(&out, "TimePointer")
  return out.P1
}

type rTimeStruct struct {
  P1 example.Time `json:"result"`
  P2 *jsoncall.Error `json:"error"`
}

func (t *DemoClient) TimeStruct() (example.Time, error) {
  var out rTimeStruct
  err := t.Client.Call(&out, "TimeStruct")
  if err != nil {
    return out.P1, err
  }
  return out.P1, out.P2.ToError()
}

type rWait struct {
  P1 *jsoncall.Error `json:"e"`
}

func (t *DemoClient) Wait(p1 int) error{
  var out rWait
  err := t.Client.Call(&out, "Wait", p1)
  if err != nil {
    return err
  }
  return out.P1.ToError()
}
