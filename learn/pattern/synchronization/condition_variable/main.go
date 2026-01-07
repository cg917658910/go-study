package main

import (
	"fmt"
	"sync"
)

type Queue struct {
	mu    sync.Mutex
	cond  *sync.Cond
	items []int
}

func NewQueue() *Queue {
	q := &Queue{
		items: make([]int, 0),
	}
	q.cond = sync.NewCond(&q.mu)
	return q
}

func (q *Queue) Enqueue(item int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, item)
	q.cond.Signal()
}

func (q *Queue) Dequeue() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	for len(q.items) == 0 {
		q.cond.Wait()
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func main() {
	queue := NewQueue()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		item := queue.Dequeue()
		fmt.Printf("Dequeued: %d\n", item)
	}()

	queue.Enqueue(42)

	wg.Wait()
}
