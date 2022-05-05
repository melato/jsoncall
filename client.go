package jsoncall

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
)

var TraceAuth bool
var TraceData bool

type Client struct {
	Url         string
	InitRequest func(r *http.Request) error
	caller      *Caller
}

func NewClient(api reflect.Type) (*Client, error) {
	var c Client
	var err error
	c.caller, err = NewJsonCaller(api)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func NewClientP(api interface{}) (*Client, error) {
	apiType := reflect.TypeOf(api)
	switch apiType.Kind() {
	case reflect.Pointer, reflect.Slice:
		return NewClient(apiType.Elem())
	default:
		return nil, fmt.Errorf("provided interface{} must be pointer or slice")
	}
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
	data, err := Marshal(args)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, t.Url+m.Method.Name, bytes.NewReader(data))
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

func (t *Client) CallE(m *Method, args ...interface{}) ([]interface{}, error) {
	data, err := t.callData(m, args)
	if err != nil {
		return nil, ToError(err)
	}
	return m.unmarshalOutputs(data)
}

func (t *Client) CallVM(result interface{}, m *Method, args ...interface{}) error {
	data, err := t.callData(m, args)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, result)
}

func (t *Client) CallV(result interface{}, name string, args ...interface{}) error {
	m := t.caller.Methods[name]
	if m == nil {
		return Errorf("no such method: %s", name)
	}
	return t.CallVM(result, m, args...)
}

func (t *Client) Call(name string, args ...interface{}) []interface{} {
	m := t.caller.Methods[name]
	var out []interface{}
	var err error
	if m != nil {
		out, err = t.CallE(m, args...)
		if err == nil {
			return out
		}
	} else {
		fmt.Printf("methods: %d\n", len(t.caller.Methods))
		panic(fmt.Sprintf("no such method: %s", name))
	}
	err = ToError(err)
	mType := m.Method.Type
	numOut := mType.NumOut()
	out = make([]interface{}, mType.NumOut())
	for i := 0; i < numOut; i++ {
		var v interface{}
		if m.OutErrors[i] {
			v = err
		} else {
			zero := reflect.Zero(mType.Out(i))
			v = zero.Interface()
		}
		out[i] = v
	}
	return out
}

func (t *Client) ToError(v interface{}) error {
	if v == nil {
		return nil
	}
	return v.(error)
}
