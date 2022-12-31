package jsoncall

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestApi(t *testing.T) {
	receiver := &TestImpl{}
	_, err := NewHttpHandler(receiver)
	if err != nil {
		t.Fail()
	}
	var test *TestInterface
	_, err = NewHttpHandler(test)
	if err != nil {
		t.Fail()
	}
}

func TestSet(t *testing.T) {
	w := newResponseWriter()
	var test TestImpl
	handler, _ := NewHttpHandler(&test)
	var request http.Request
	request.URL = &url.URL{}
	request.URL.Path = "/SetA"
	body := `{"p1":2}`
	request.Body = io.NopCloser(strings.NewReader(body))
	handler.ServeHTTP(w, &request)
	s := w.buf.String()
	fmt.Printf("%s\n", s)
	w.CheckStatus(t, http.StatusOK)
	if s != `{}` {
		t.Fail()
	}
	if test.A != 2 {
		t.Fail()
	}
}

func TestUnprocessableEntity(t *testing.T) {
	w := newResponseWriter()
	receiver := &TestImpl{}
	handler, _ := NewHttpHandler(receiver)
	var request http.Request
	request.URL = &url.URL{}
	request.URL.Path = "/Div"
	body := `{"p1":1,"p2":0}`
	request.Body = io.NopCloser(strings.NewReader(body))
	handler.ServeHTTP(w, &request)
	s := w.buf.String()
	fmt.Printf("%s\n", s)
	w.CheckStatus(t, http.StatusUnprocessableEntity)
	if s != `{"error":"division by 0","result":0}` {
		t.Fail()
	}
}

func TestNotFound(t *testing.T) {
	w := newResponseWriter()
	receiver := &TestImpl{}
	handler, _ := NewHttpHandler(receiver)
	var request http.Request
	request.URL = &url.URL{}
	request.URL.Path = "/ccc"
	handler.ServeHTTP(w, &request)
	s := w.buf.String()
	fmt.Printf("%s\n", s)
	w.CheckStatus(t, http.StatusNotFound)
}

func TestBadRequest(t *testing.T) {
	w := newResponseWriter()
	receiver := &TestImpl{}
	handler, _ := NewHttpHandler(receiver)
	var request http.Request
	request.URL = &url.URL{}
	request.URL.Path = "/Div"
	body := `{"p1":""}`
	request.Body = io.NopCloser(strings.NewReader(body))
	handler.ServeHTTP(w, &request)
	s := w.buf.String()
	fmt.Printf("%s\n", s)
	w.CheckStatus(t, http.StatusBadRequest)
}
