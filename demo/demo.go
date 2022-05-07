package demo

type Demo interface {
	Ping() error
	Hello() (string, error)
	Add(a, b int32) (int32, error)
	Div(a, b int32) (int32, error) // can cause a division by zero
	Wait(seconds int) error
}
