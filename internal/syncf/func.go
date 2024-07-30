package syncf

// FuncQueue maintains a queue that runs functions sequentially, in a single goroutine.
type FuncQueue chan func()

// NewFuncQueue creates a queue with a channel of 0 capacity
// Such a queue cannot have nested calls
func NewFuncQueue() FuncQueue {
	t := FuncQueue(make(chan func(), 0))
	go t.runner()
	return t
}

func (t FuncQueue) runner() {
	for fn := range t {
		fn()
	}
}

// RunAsync - put a function in the queue and return immediately
// RunAsync makes it possible to nest function calls.
// If a function in the queue needs to put another function in the same queue,
// It can do so with RunAsync().
// The queue capacity should be no less than the nesting depth.
// Otherwise a deadlock can occur.
func (t FuncQueue) RunAsync(fn func()) {
	t <- fn
}
