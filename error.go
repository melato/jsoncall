package jsoncall

import (
	"fmt"
)

type Error string

func (t *Error) Error() string {
	return string(*t)
}

func (t *Error) ToError() error {
	if t == nil {
		return nil
	}
	return t
}

func newError(s string) *Error {
	e := Error(s)
	return &e
}

/*
type Error struct {
	Message string `json:"error"`
}

func (t *Error) Error() string {
	return t.Message
}
*/

func Errorf(format string, args ...interface{}) *Error {
	return newError(fmt.Sprintf(format, args...))
}

func ToError(err error) *Error {
	e, isError := err.(*Error)
	if isError {
		return e
	}
	se := Error(err.Error())
	return &se
}

type ErrorCode int

const (
	ErrNone         = iota
	ErrNoSuchMethod = iota
	ErrMarshal      = iota
	ErrUser         = iota
)
