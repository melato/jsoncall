package jsoncall

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
)

type ReceiverProvider func(ReceiverContext) (interface{}, error)

type ReceiverContext interface {
	Request() *http.Request
	WriteHeader(statusCode int)
	MethodName() string
}

// HttpHandler - net/http.Handler that maps POST requests to method calls on a receiver
type HttpHandler struct {
	receiver    ReceiverProvider
	Caller      *Caller
	methodPaths map[string]*Method
}

func (caller *Caller) NewHttpHandler(receiver ReceiverProvider) *HttpHandler {
	var t HttpHandler
	t.Caller = caller
	t.receiver = receiver
	t.methodPaths = make(map[string]*Method)
	for _, m := range t.Caller.methods {
		t.methodPaths[m.Desc.Path] = m
	}
	return &t
}

func NewHttpHandler(receiver interface{}) (http.Handler, error) {
	caller, err := NewCaller(receiver, nil)
	if err != nil {
		return nil, err
	}
	var f ReceiverProvider = func(context ReceiverContext) (interface{}, error) { return receiver, nil }
	return caller.NewHttpHandler(f), nil
}

// SetReceiverProvider - provides a method receiver for each call
// If it returns nil, processing of the request stops.
/*
func (t *HttpHandler) SetReceiverProvider(r ReceiverProvider) {
	t.receiver = r
}
*/

func (t *HttpHandler) getBytes(r *http.Request) ([]byte, error) {
	var buf bytes.Buffer
	defer r.Body.Close()
	_, err := io.Copy(&buf, r.Body)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *HttpHandler) writeResponse(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

func (t *HttpHandler) writeError(w http.ResponseWriter, status int, e error) {
	data, _ := json.Marshal(ToError(e))
	t.writeResponse(w, status, data)
}

func (t *HttpHandler) DoService(receiver interface{}, m *Method, w http.ResponseWriter, r *http.Request) {
	inputData, err := t.getBytes(r)

	var outputData []byte
	status := http.StatusInternalServerError
	var errCode ErrorCode
	if err == nil {
		outputData, errCode, err = m.Call(receiver, inputData)
		if TraceData {
			fmt.Printf("call result: %s %v\n", string(outputData), err)
		}
		switch errCode {
		case ErrNone:
			status = http.StatusOK
		case ErrMarshal:
			status = http.StatusBadRequest
		case ErrNoSuchMethod:
			status = http.StatusNotFound
		case ErrUser:
			status = http.StatusInternalServerError
		default:
			status = http.StatusInternalServerError
		}
	}
	if err != nil {
		t.writeError(w, status, err)
	} else {
		t.writeResponse(w, status, outputData)
	}
}

type methodContext struct {
	request *http.Request
	writer  http.ResponseWriter
	method  *Method
}

func (c *methodContext) Request() *http.Request {
	return c.request
}

func (c *methodContext) WriteHeader(statusCode int) {
	c.writer.WriteHeader(statusCode)
}

func (c *methodContext) MethodName() string {
	return c.method.Desc.Method
}

func (t *HttpHandler) ServeMethod(m *Method, w http.ResponseWriter, r *http.Request) {
	var receiver interface{}
	var err error
	if t.receiver != nil {
		c := methodContext{request: r, writer: w, method: m}
		receiver, err = t.receiver(&c)
	}
	if err == nil {
		t.DoService(receiver, m, w, r)
	} else {
		outputData, err2 := json.Marshal(ToError(err))
		if err2 == nil {
			w.Write(outputData)
		}
	}
	return
}

func (t *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := path.Base(r.URL.Path)
	if TraceCalls {
		fmt.Printf("path: %s, method=%s\n", r.URL.Path, method)
	}
	m, found := t.methodPaths[method]
	if found {
		if TraceCalls {
			fmt.Printf("method: %s\n", m.Desc.Method)
		}
		t.ServeMethod(m, w, r)
	} else {
		t.writeError(w, http.StatusNotFound, fmt.Errorf("unknown method: %v.%s", t.Caller.rType, method))
	}
}
