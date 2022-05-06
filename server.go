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
	Port int32
	// ReceiverFunc - provides a method receiver for each call
	// If it returns nil, processing of the request stops.
	ReceiverFunc func(w http.ResponseWriter, r *http.Request) (interface{}, error)
	Caller       *Caller
	methodPaths  map[string]*Method
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

func (t *Server) writeResponse(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

func (t *Server) writeError(w http.ResponseWriter, status int, e error) {
	data, _ := json.Marshal(ToError(e))
	t.writeResponse(w, status, data)
}

func (t *Server) DoService(receiver interface{}, m *Method, w http.ResponseWriter, r *http.Request) {
	inputData, err := t.getBytes(r)

	var outputData []byte
	status := http.StatusInternalServerError
	var errCode ErrorCode
	if err == nil {
		outputData, errCode, err = t.Caller.Call(m, receiver, inputData)
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

func (t *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	path = strings.TrimPrefix(path, "/")
	m, found := t.methodPaths[path]
	if TraceCalls {
		fmt.Printf("path: %s method: %s\n", path, m.Names.Method)
	}
	if found {
		t.ServeMethod(m, w, r)
	} else {
		t.writeError(w, http.StatusNotFound, fmt.Errorf("unknown api path: %v/%s", t.Caller.Api, path))
	}
}

func (t *Server) ServeMethod(m *Method, w http.ResponseWriter, r *http.Request) {
	receiver, err := t.ReceiverFunc(w, r)
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

func (t *Server) Run() error {
	t.methodPaths = make(map[string]*Method)
	for _, m := range t.Caller.Methods {
		t.methodPaths[m.Names.Path] = m
	}
	server := &http.Server{Addr: fmt.Sprintf(":%d", t.Port), Handler: t}
	log.Printf("starting server at %s\n", server.Addr)
	err := server.ListenAndServe()
	return err
}
