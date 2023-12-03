package jsoncall

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
)

// ReceiverFunc is a function that is called at the beginning of each request
// to provide a receiver for the request method.
// If it returns nil, then the processing of the request stops, and no receiver methods are called.
// In that case, it should write an appropriate error response to the ResponseWriter.
//
// For example, it may look at the headers of the request, determine
// the user that is making the request, and incorporate this information
// in the receiver.
type ReceiverFunc func(w http.ResponseWriter, r *http.Request) interface{}

func NewReceiverFunc(receiver interface{}) ReceiverFunc {
	return func(w http.ResponseWriter, r *http.Request) interface{} { return receiver }
}

// HttpHandler - net/http.Handler that maps POST requests to method calls on a receiver
type HttpHandler struct {
	ReceiverFunc ReceiverFunc
	Caller       *Caller
	methodPaths  map[string]*Method
	synchronizer func(func())
}

func (caller *Caller) NewHttpHandler(receiver ReceiverFunc) *HttpHandler {
	var t HttpHandler
	t.Caller = caller
	t.ReceiverFunc = receiver
	t.methodPaths = make(map[string]*Method)
	for _, m := range t.Caller.methods {
		t.methodPaths[m.Desc.Path] = m
	}
	return &t
}

// NewHttpHandler creates an http.Handler, using the default API descriptor.
// The methods and method signatures of the API are the methods of the prototype.
// prototype is either a pointer to an interface or a pointer to a struct.
// If it is a pointer to a struct, it is also used as the receiver of the methods,
// unless another receiver is specified with SetReceiver() or SetReceiverFunc().
// If you want to specify an Api Descriptor, use a Caller to create the HttpHandler.
func NewHttpHandler(prototype interface{}) (*HttpHandler, error) {
	caller, err := NewCaller(prototype, nil)
	if err != nil {
		return nil, err
	}
	handler := caller.NewHttpHandler(nil)
	handler.SetReceiver(prototype)
	handler.synchronizer = func(f func()) { f() }
	return handler, nil
}

// SetSynchronizer installs a synchronizer.
// A synchronizer is a function that wraps another function call,
// in order to provide synchronization.
// The default synchronizer runs functions immediately.
// If synchronization is desired, set a synchronizer
// that runs functions sequentially in a single goroutine.
func (t *HttpHandler) SetSynchronizer(synchronizer func(func())) {
	t.synchronizer = synchronizer
}

// SetReceiverFunc specifies a function that is called for each request to produce
// the receiver for the requested method.
// If it returns nil, the requested method is not called.
func (t *HttpHandler) SetReceiverFunc(f ReceiverFunc) {
	t.ReceiverFunc = f
}

// SetReceiver specifies the receiver to be used for all requested methods.
// Calling SetReceiver is equivalent to calling SetReceiveFunc() with a function
// that returns receiver.
func (t *HttpHandler) SetReceiver(receiver interface{}) {
	t.SetReceiverFunc(func(w http.ResponseWriter, r *http.Request) interface{} { return receiver })
}

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
	if TraceCalls {
		fmt.Printf("writeRespose status=%d len(data)=%d\n", status, len(data))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

func (t *HttpHandler) writeError(w http.ResponseWriter, status int, e error) {
	data, _ := json.Marshal(ToError(e))
	t.writeResponse(w, status, data)
}

func (t *HttpHandler) ServeMethod(m *Method, w http.ResponseWriter, r *http.Request) {
	inputData, err := t.getBytes(r)
	if TraceData {
		fmt.Printf("getBytes: %d %v\n", len(inputData), err)
	}

	var outputData []byte
	status := http.StatusInternalServerError
	var errCode ErrorCode
	if err == nil {
		var receiver interface{}

		t.synchronizer(func() {
			if t.ReceiverFunc != nil {
				receiver = t.ReceiverFunc(w, r)
			}
			if receiver == nil {
				return
			}

			outputData, errCode, err = m.Call(receiver, inputData)
		})
		if TraceData {
			fmt.Printf("call result: %s %v\n", outputData, err)
		}
		switch errCode {
		case ErrNone:
			status = http.StatusOK
		case ErrMarshal:
			status = http.StatusBadRequest
		case ErrNoSuchMethod:
			status = http.StatusNotImplemented
		case ErrUser:
			status = http.StatusUnprocessableEntity
		default:
			status = http.StatusInternalServerError
		}
	}
	if err != nil && !m.HasErrors() {
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
