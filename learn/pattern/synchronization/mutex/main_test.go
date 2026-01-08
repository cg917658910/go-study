package main

import (
	"sync"
	"testing"
)

func TestMutex(t *testing.T) {
	counter := &Counter{}
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Inc()
		}()
	}

	wg.Wait()

	if counter.Value() != 100 {
		t.Errorf("Expected 100, got %d", counter.Value())
	}
}
