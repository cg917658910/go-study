package main

import (
	"fmt"
	"sync"
)

type Barrier struct {
	n  int
	wg sync.WaitGroup
}

func NewBarrier(n int) *Barrier {
	b := &Barrier{n: n}
	b.wg.Add(n)
	return b
}

func (b *Barrier) Wait() {
	b.wg.Done()
	b.wg.Wait()
}

func main() {
	n := 3
	barrier := NewBarrier(n)

	for i := 1; i <= n; i++ {
		go func(id int) {
			fmt.Printf("Goroutine %d waiting\n", id)
			barrier.Wait()
			fmt.Printf("Goroutine %d passed\n", id)
		}(i)
	}

	// Wait for all to complete
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Give time for goroutines to print
		for i := 0; i < 1000000; i++ {
		}
	}()
	wg.Wait()
}
