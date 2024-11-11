package syncf

import (
	"fmt"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	q := make(Queue)
	go q.Evaluate()
	var d int
	done := make(chan struct{})
	q.Put(func() {
		fmt.Printf("running")
		time.Sleep(time.Duration(10) * time.Millisecond)
		d = 7
		done <- struct{}{}
	})
	if d != 0 {
		t.Fatalf("expected 0: %d", d)
	}
	_ = <-done
	if d != 7 {
		t.Fatalf("expected 7: %d", d)
	}
}
