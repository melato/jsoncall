package demo

type Demo interface {
	Ping() error
	Hello() (string, error)
	Wait(seconds int) error
}
