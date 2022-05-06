package demo

type Demo interface {
	Ping() error
	Hello() (string, error)
	Add(a, b int32) (int32, error)
	Wait(seconds int) error
}
