package main

import (
"fmt"
"sync"
)

type Counter struct {
mu    sync.Mutex
value int
}

func (c *Counter) Inc() {
c.mu.Lock()
defer c.mu.Unlock()
c.value++
}

func (c *Counter) Value() int {
c.mu.Lock()
defer c.mu.Unlock()
return c.value
}

func main() {
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
fmt.Printf("Counter: %d\n", counter.Value())
}
