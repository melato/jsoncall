package jsoncall

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// Server - provides a web service that maps requests to method calls on a receiver
type Server struct {
	Port  int32
	Trace bool
	// ReceiverFunc - provides a method receiver for each call
	// If it returns nil, processing of the request stops.
	ReceiverFunc func(w http.ResponseWriter, r *http.Request) (interface{}, error)
	caller       *Caller
}

func NewServer(api interface{}) (*Server, error) {
	var s Server
	var err error
	s.caller, err = NewJsonCallerP(api)
	if err != nil {
		return nil, err
	}
	s.ReceiverFunc = func(w http.ResponseWriter, r *http.Request) (interface{}, error) { return api, nil }
	return &s, nil
}

func (t *Server) getBytes(r *http.Request) ([]byte, error) {
	var buf bytes.Buffer
	defer r.Body.Close()
	_, err := io.Copy(&buf, r.Body)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *Server) DoService(receiver interface{}, method string, w http.ResponseWriter, r *http.Request) {
	inputData, err := t.getBytes(r)

	var outputData []byte
	var status int
	if err == nil {
		outputData, err = t.caller.Call(method, receiver, inputData)
		if TraceData {
			fmt.Printf("call result: %s %v\n", string(outputData), err)
		}
		if err == nil {
			status = http.StatusOK
		}
	}
	if err != nil {
		status = http.StatusInternalServerError
		outputData, err = json.Marshal(ToError(err))
	}
	w.WriteHeader(status)
	w.Write(outputData)
}

func (t *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	path = strings.TrimPrefix(path, "/")
	if TraceCalls {
		fmt.Printf("method: %s\n", path)
	}
	t.ServeMethod(path, w, r)
}

func (t *Server) ServeMethod(method string, w http.ResponseWriter, r *http.Request) {
	receiver, err := t.ReceiverFunc(w, r)
	if err == nil {
		t.DoService(receiver, method, w, r)
	} else {
		outputData, err2 := json.Marshal(ToError(err))
		if err2 == nil {
			w.Write(outputData)
		}
	}
	return
}

func (t *Server) Run() error {
	server := &http.Server{Addr: fmt.Sprintf(":%d", t.Port), Handler: t}
	log.Printf("starting server at %s\n", server.Addr)
	err := server.ListenAndServe()
	return err
}
