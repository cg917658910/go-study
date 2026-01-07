package concurrency

import (
	"fmt"
	"sync"
	"testing"
)

func TestNBarrier(t *testing.T) {
	n := 3 // 需要同步的协程数量
	barrier := NewNBarrier(n)
	var wg sync.WaitGroup

	// 模拟多个协程
	for i := 1; i <= n; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Goroutine %d is waiting at the barrier\n", id)
			barrier.Wait() // 等待其他协程
			fmt.Printf("Goroutine %d passed the barrier\n", id)
		}(i)
	}

	wg.Wait()
}

func TestNBarrier_MultipleCycles(t *testing.T) {
	n := 2 // 每次屏障需要同步的协程数量
	barrier := NewNBarrier(n)
	var wg sync.WaitGroup

	// 模拟多个协程，测试多次屏障同步
	for i := 1; i <= n; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for cycle := 1; cycle <= 3; cycle++ {
				fmt.Printf("Goroutine %d is waiting at the barrier (Cycle %d)\n", id, cycle)
				barrier.Wait() // 等待其他协程
				fmt.Printf("Goroutine %d passed the barrier (Cycle %d)\n", id, cycle)
			}
		}(i)
	}

	wg.Wait()
}
