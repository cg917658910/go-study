package main

import (
"fmt"
"sync"
)

type Task func()

type ThreadPool struct {
tasks   chan Task
workers int
wg      sync.WaitGroup
}

func NewThreadPool(workers int) *ThreadPool {
return &ThreadPool{
tasks:   make(chan Task, 100),
workers: workers,
}
}

func (p *ThreadPool) Start() {
for i := 0; i < p.workers; i++ {
p.wg.Add(1)
go p.worker()
}
}

func (p *ThreadPool) worker() {
defer p.wg.Done()
for task := range p.tasks {
task()
}
}

func (p *ThreadPool) Submit(task Task) {
p.tasks <- task
}

func (p *ThreadPool) Shutdown() {
close(p.tasks)
p.wg.Wait()
}

func main() {
pool := NewThreadPool(3)
pool.Start()

for i := 0; i < 10; i++ {
id := i
pool.Submit(func() {
fmt.Printf("Task %d executed\n", id)
})
}

pool.Shutdown()
}
