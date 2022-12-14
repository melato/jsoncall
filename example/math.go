package example

import "fmt"

type Math interface {
	Add(a, b int32) (int32, error)
	Div(a, b int32) (int32, error) // can cause a division by zero
}

type MathImpl struct {
}

func (t *MathImpl) Add(a, b int32) (int32, error) {
	return a + b, nil
}

func (t *MathImpl) Div(a, b int32) (int32, error) {
	if b == 0 {
		return 0, fmt.Errorf("Division by zero")
	}
	return a / b, nil
}
