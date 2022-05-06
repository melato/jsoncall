package jsoncall

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var TraceAuth bool
var TraceData bool

type Client struct {
	Url         string
	InitRequest func(r *http.Request) error
	Caller      *Caller
}

func (t *Client) responseError(response *http.Response, data []byte) error {
	var e Error
	err := json.Unmarshal(data, &e)
	fmt.Printf("unmarshal error: %v\n", err)
	if err != nil {
		return fmt.Errorf(response.Status)
	}
	return &e
}

func (t *Client) callData(m *Method, args []interface{}) ([]byte, error) {
	data, err := m.MarshalInputs(args...)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, t.Url+m.Names.Path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	if t.InitRequest != nil {
		err := t.InitRequest(request)
		if err != nil {
			return nil, err
		}
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "*/*")
	var h http.Client
	response, err := h.Do(request)
	if err != nil {
		return nil, err
	}
	data, err = ioutil.ReadAll(response.Body)
	if TraceData {
		fmt.Printf("status: %d body: %s\n", response.StatusCode, string(data))
	}
	if err != nil {
		return nil, err
	}
	if 200 <= response.StatusCode && response.StatusCode < 300 {
		return data, nil
	}
	return nil, t.responseError(response, data)
}

func (t *Client) CallVM(result interface{}, m *Method, args ...interface{}) error {
	data, err := t.callData(m, args)
	if err != nil {
		return err
	}
	if result != nil {
		return json.Unmarshal(data, result)
	}
	return nil
}

func (t *Client) CallV(result interface{}, name string, args ...interface{}) error {
	m := t.Caller.Methods[name]
	if m == nil {
		return Errorf("no such method: %s", name)
	}
	return t.CallVM(result, m, args...)
}

func (t *Client) ToError(v interface{}) error {
	if v == nil {
		return nil
	}
	return v.(error)
}
