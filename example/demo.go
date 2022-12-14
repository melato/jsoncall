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
