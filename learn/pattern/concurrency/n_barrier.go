package concurrency

// N-Barrier Pattern N屏障模式 并发设计模式
// 允许多个协程在达到某个点之前等待彼此
// 适用于需要多个协程在继续执行之前同步的场景

import (
	"sync"
)

type NBarrier struct {
	n          int
	count      int
	mutex      sync.Mutex
	cond       *sync.Cond
	generation int
}

func NewNBarrier(n int) *NBarrier {
	nb := &NBarrier{n: n}
	nb.cond = sync.NewCond(&nb.mutex)
	return nb
}

func (nb *NBarrier) Wait() {
	nb.mutex.Lock()
	defer nb.mutex.Unlock()

	gen := nb.generation
	nb.count++
	if nb.count == nb.n {
		nb.generation++
		nb.count = 0
		nb.cond.Broadcast()
	} else {
		for gen == nb.generation {
			nb.cond.Wait()
		}
	}
}
