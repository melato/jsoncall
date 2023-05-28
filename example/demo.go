package example

import (
	"fmt"
	"time"
)

type Demo interface {
	// Repeat returns a []string consisting repeated elements
	// It demonstrates  inputs and outputs of various types.
	// It is also an easy way to generate large responses.
	Repeat(s string, count int) ([]string, error)

	// Time takes no arguments, returns multiple values, and does not return an error.
	Time() (hours, minutes, seconds int)

	// TimeStruct returns a struct
	TimeStruct() (Time, error)

	// TimePointer returns a struct pointer
	TimePointer() *Time

	// Wait waits the specified number of seconds
	// It can be used to test a long-running response.
	Wait(seconds int) error

	// Ping should return no errors
	// can be used to verify that communication is good
	Ping() error

	// Error returns a string and an error.
	// Tests error handling.
	Error() (string, error)
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

func (t *DemoImpl) Repeat(s string, count int) ([]string, error) {
	if count < 0 {
		return nil, fmt.Errorf("negative count: %d", count)
	}
	list := make([]string, count)
	for i := 0; i < count; i++ {
		list[i] = s
	}
	return list, nil
}

func (t *DemoImpl) Ping() error {
	return nil
}

func (t *DemoImpl) Error() (string, error) {
	return "test", fmt.Errorf("err")
}
