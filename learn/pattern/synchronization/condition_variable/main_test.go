package main

import (
"sync"
"testing"
)

func TestConditionVariable(t *testing.T) {
queue := NewQueue()
var wg sync.WaitGroup

wg.Add(1)
go func() {
defer wg.Done()
item := queue.Dequeue()
if item != 42 {
t.Errorf("Expected 42, got %d", item)
}
}()

queue.Enqueue(42)

wg.Wait()
}
