package jsoncall

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

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

	// Url is prepended to the name of each method.  Should include trailing "/"
	Url string

	// InitRequest initializes the request (optional).  Can be used to set authorization headers.
	InitRequest func(r *http.Request) error
}

// url should include trailing "/"
func (caller *Caller) NewHttpClient(url string) *HttpClient {
	t := NewHttpClient(url)
	t.Caller = caller
	return t
}

func NewHttpClient(url string) *HttpClient {
	var t HttpClient
	t.Url = url
	return &t
}

func (t *HttpClient) responseError(response *http.Response, data []byte) error {
	var e Error
	err := json.Unmarshal(data, &e)
	if err != nil {
		return fmt.Errorf(response.Status)
	}
	return &e
}

type statusError struct {
	Status int
}

func (t *statusError) Error() string {
	return fmt.Sprintf("http status %d", t.Status)
}

func (t *HttpClient) callData(name string, data []byte) ([]byte, error) {
	url := t.Url + name
	if TraceCalls {
		fmt.Printf("%s %s\n", http.MethodPost, url)
		fmt.Printf("%s\n", string(data))
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
		fmt.Printf("status: %d\n", response.StatusCode)
		fmt.Printf("%s\n", string(data))
	}
	if err != nil {
		return nil, err
	}
	if 200 <= response.StatusCode && response.StatusCode < 300 {
		return data, nil
	}
	return data, &statusError{response.StatusCode}
}

func IsErrStatus(err error) bool {
	var statusError *statusError
	return errors.As(err, &statusError)
}

// CallJson makes a POST HTTP call.
// It converts input to JSON and sends it as the body of the request.
// If there are no errors, and the response status is 2xx, it unmarshals the response body to output
// If the response status is not 2xx and errorOutput is not nil, and there is no other error, it unmarshals the response body to errorOutput
func (t *HttpClient) CallJson(name string, input interface{}, output interface{}, errorOutput interface{}) error {
	if TraceData {
		fmt.Printf("CallJson %s input:%T output:%T errorOutput:%T\n", name, input, output, errorOutput)
	}
	data, err := json.Marshal(input)
	if err != nil {
		return err
	}
	data, err = t.callData(name, data)
	if err == nil {
		if output != nil {
			return json.Unmarshal(data, output)
		}
		return nil
	}
	if errorOutput == nil || !IsErrStatus(err) {
		return err
	}
	_ = json.Unmarshal(data, errorOutput)
	return err
}

func (t *HttpClient) Call(output interface{}, name string, args ...interface{}) error {
	method := t.Caller.methods[name]
	if method == nil {
		return Errorf("no such method: %s", name)
	}
	inputs, err := method.Desc.InputMap(args)
	if err != nil {
		return fmt.Errorf("marshal %s: %w", method.Desc.Method, err)
	}
	return t.CallJson(method.Desc.Path, inputs, output, nil)
}
