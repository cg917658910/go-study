package main

import (
"testing"
)

func TestFlyweightFactory(t *testing.T) {
factory := NewFlyweightFactory()

fw1 := factory.GetFlyweight("A")
fw2 := factory.GetFlyweight("A")

if fw1 != fw2 {
t.Error("Same key should return same flyweight instance")
}

if factory.GetFlyweightCount() != 1 {
t.Errorf("Expected 1 flyweight, got %d", factory.GetFlyweightCount())
}
}

func TestFlyweightMultipleKeys(t *testing.T) {
factory := NewFlyweightFactory()

factory.GetFlyweight("A")
factory.GetFlyweight("B")
factory.GetFlyweight("C")

if factory.GetFlyweightCount() != 3 {
t.Errorf("Expected 3 flyweights, got %d", factory.GetFlyweightCount())
}
}
