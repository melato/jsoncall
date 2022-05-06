package demo

import (
	"time"
)

type DemoImpl struct {
}

func (t *DemoImpl) Ping() error {
	return nil
}
func (t *DemoImpl) Hello() (string, error) {
	return "hello", nil
}

func (t *DemoImpl) Add(a, b int32) (int32, error) {
	return a + b, nil
}

func (t *DemoImpl) Concat(a, b string) (string, error) {
	return a + b, nil
}

func (t *DemoImpl) Wait(seconds int) error {
	time.Sleep(time.Duration(seconds) * time.Second)
	return nil
}
