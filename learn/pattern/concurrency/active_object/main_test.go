package main

import (
	"sync/atomic"
	"testing"
)

func TestActiveObject(t *testing.T) {
	ao := NewActiveObject()

	var counter int32
	for i := 0; i < 10; i++ {
		ao.DoWork(func() {
			atomic.AddInt32(&counter, 1)
		})
	}

	ao.Shutdown()

	if counter != 10 {
		t.Errorf("Expected 10, got %d", counter)
	}
}
