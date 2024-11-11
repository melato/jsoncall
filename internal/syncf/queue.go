package syncf

// Queue - a channel of functions that are called in one or more goroutines
type Queue chan func()

// Evaluate all the functions in the Queue, forever
//
// Use like this:
//
//	go q.Evaluate()
func (q Queue) Evaluate() {
	for fn := range q {
		fn()
	}
}

// Put - another way of appending a function to the channel.
func (q Queue) Put(fn func()) {
	q <- fn
}
