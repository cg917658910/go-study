package main

import (
"sync/atomic"
"testing"
)

func TestThreadPool(t *testing.T) {
pool := NewThreadPool(3)
pool.Start()

var counter int32
for i := 0; i < 10; i++ {
pool.Submit(func() {
atomic.AddInt32(&counter, 1)
})
}

pool.Shutdown()

if counter != 10 {
t.Errorf("Expected 10 tasks, got %d", counter)
}
}
