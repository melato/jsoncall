package jsoncall

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var TraceData bool

// Client is the interface used by generated clients.  It can also be used by itself.
type Client interface {
	// Call a remote method.  method, args will be used to make the request.
	// response will be passed to Json.Unmarshal with the response JSON.
	// A response type that should always work is *map[string]interface{}
	// Generated clients use  more specific response types by knowing the expected
	// return types of each method.
	Call(response interface{}, method string, args ...interface{}) error
}

type HttpClient struct {
	// Caller specifies the client/server API
	Caller *Caller

	// Url is prepended to the path for each method.  Required.  Should include trailing "/"
	Url string

	// InitRequest initializes the request (optional).  Can be used to set authorization headers.
	InitRequest func(r *http.Request) error
}

func (t *HttpClient) responseError(response *http.Response, data []byte) error {
	var e Error
	err := json.Unmarshal(data, &e)
	if err != nil {
		return fmt.Errorf(response.Status)
	}
	return &e
}

func (t *HttpClient) callData(m *Method, args []interface{}) ([]byte, error) {
	data, err := m.MarshalInputs(args...)
	if err != nil {
		return nil, fmt.Errorf("marshal %s: %v", m.Desc.Method, err)
	}
	var url string
	if t.Caller.Prefix != "" {
		url = t.Url + t.Caller.Prefix + "/" + m.Desc.Path
	} else {
		url = t.Url + m.Desc.Path
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
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

func (t *HttpClient) callMethod(result interface{}, m *Method, args ...interface{}) error {
	data, err := t.callData(m, args)
	if err != nil {
		return err
	}
	if result != nil {
		return json.Unmarshal(data, result)
	}
	return nil
}

func (t *HttpClient) Call(result interface{}, name string, args ...interface{}) error {
	m := t.Caller.methods[name]
	if m == nil {
		return Errorf("no such method: %s", name)
	}
	return t.callMethod(result, m, args...)
}
