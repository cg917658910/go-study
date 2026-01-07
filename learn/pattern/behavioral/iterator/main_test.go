package main

import "testing"

func TestIterator(t *testing.T) {
collection := &ConcreteCollection{
items: []interface{}{1, 2, 3},
}

iterator := collection.CreateIterator()
count := 0
for iterator.HasNext() {
iterator.Next()
count++
}

if count != 3 {
t.Errorf("Expected 3 items, got %d", count)
}
}
