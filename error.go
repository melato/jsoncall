package jsoncall

import (
	"fmt"
)

type Error struct {
	Message string `json:"error"`
}

func (t *Error) Error() string {
	return t.Message
}

func Errorf(format string, args ...interface{}) *Error {
	return &Error{fmt.Sprintf(format, args...)}
}

func ToError(err error) *Error {
	e, isError := err.(*Error)
	if isError {
		return e
	}
	return &Error{err.Error()}
}
