package main

import (
"sync"
"testing"
)

func TestBarrier(t *testing.T) {
n := 3
barrier := NewBarrier(n)
var wg sync.WaitGroup

for i := 0; i < n; i++ {
wg.Add(1)
go func() {
defer wg.Done()
barrier.Wait()
}()
}

wg.Wait()
}
