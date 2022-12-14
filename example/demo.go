package example

import (
	"time"
)

type Demo interface {
	Nop()
	Ping() error
	Hello() (string, error)
	Wait(seconds int) error
	Seconds(hours, minutes, seconds int) (int, error)
	Time() (hours, minutes, seconds int, err error)
	Substring(s string, start int, length int) (string, error)
}

type DemoImpl struct {
}

func (t *DemoImpl) Ping() error {
	return nil
}

func (t *DemoImpl) Nop() {
}

func (t *DemoImpl) Hello() (string, error) {
	return "hello", nil
}

func (t *DemoImpl) Seconds(hours, minutes, seconds int) (int, error) {
	return hours*3600 + minutes*60 + seconds, nil
}

func (t *DemoImpl) Wait(seconds int) error {
	time.Sleep(time.Duration(seconds) * time.Second)
	return nil
}

func (t *DemoImpl) Time() (hour, minute, second int, err error) {
	now := time.Now()
	return now.Hour(), now.Minute(), now.Second(), nil
}

func (t *DemoImpl) Substring(s string, start int, length int) (string, error) {
	return s[start : start+length], nil
}
