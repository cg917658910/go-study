package main

import (
	"sync"
	"testing"
)

func TestRWLock(t *testing.T) {
	cache := NewCache()
	var wg sync.WaitGroup

	cache.Set("key", "value")

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = cache.Get("key")
		}()
	}

	wg.Wait()
}
