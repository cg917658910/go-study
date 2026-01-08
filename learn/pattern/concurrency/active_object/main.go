package main

import (
	"fmt"
	"sync"
)

type ActiveObject struct {
	requests chan func()
	wg       sync.WaitGroup
}

func NewActiveObject() *ActiveObject {
	ao := &ActiveObject{
		requests: make(chan func(), 100),
	}
	ao.wg.Add(1)
	go ao.run()
	return ao
}

func (ao *ActiveObject) run() {
	defer ao.wg.Done()
	for req := range ao.requests {
		req()
	}
}

func (ao *ActiveObject) DoWork(work func()) {
	ao.requests <- work
}

func (ao *ActiveObject) Shutdown() {
	close(ao.requests)
	ao.wg.Wait()
}

func main() {
	ao := NewActiveObject()

	for i := 0; i < 5; i++ {
		id := i
		ao.DoWork(func() {
			fmt.Printf("Task %d\n", id)
		})
	}

	ao.Shutdown()
}
