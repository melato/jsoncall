package server

import (
	"fmt"
	"time"

	"melato.org/jsoncall/example"
)

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

func (t *DemoImpl) TimeStruct() (example.Time, error) {
	var m example.Time
	m.Hour, m.Minute, m.Second = t.Time()
	return m, nil
}

func (t *DemoImpl) TimePointer() *example.Time {
	var m example.Time
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
