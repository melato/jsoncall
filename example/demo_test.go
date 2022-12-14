package example

import (
	"testing"
)

func TestDemoInterface(t *testing.T) {
	var d Demo
	var impl *DemoImpl
	d = impl
	_ = d
}
