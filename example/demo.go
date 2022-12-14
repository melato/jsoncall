package example

import (
	"strings"
	"time"
)

type Demo interface {
	// Repeat calls strings.Repeat.
	// It demonstrates  passing arguments of different types.
	// It is also an easy way to generate large responses.
	Repeat(s string, count int) (string, error)

	// Time takes no arguments, returns multiple values, and does not return an error.
	Time() (hours, minutes, seconds int)

	// TimeStruct returns a struct
	TimeStruct() (Time, error)

	// TimePointer returns a struct pointer
	TimePointer() *Time

	// Wait waits the specified number of seconds
	// It can be used to test a long-running response.
	Wait(seconds int) error
}

type Time struct {
	Hour, Minute, Second int
}

type DemoImpl struct {
}

// Extra a method that is not in the interface
func (t *DemoImpl) Extra() {
}

func (t *DemoImpl) Wait(seconds int) error {
	time.Sleep(time.Duration(seconds) * time.Second)
	return nil
}

func (t *DemoImpl) Time() (hour, minute, second int) {
	now := time.Now()
	return now.Hour(), now.Minute(), now.Second()
}

func (t *DemoImpl) TimeStruct() (Time, error) {
	var m Time
	m.Hour, m.Minute, m.Second = t.Time()
	return m, nil
}

func (t *DemoImpl) TimePointer() *Time {
	var m Time
	m.Hour, m.Minute, m.Second = t.Time()
	return &m
}

func (t *DemoImpl) Repeat(s string, count int) (string, error) {
	return strings.Repeat(s, count), nil
}
