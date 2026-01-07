package main

import (
"testing"
)

func TestBridge(t *testing.T) {
implA := &ConcreteImplementorA{}
implB := &ConcreteImplementorB{}

abstraction1 := NewAbstraction(implA)
if abstraction1.Operation() != "ConcreteImplementorA" {
t.Errorf("Expected ConcreteImplementorA, got %s", abstraction1.Operation())
}

abstraction2 := NewAbstraction(implB)
if abstraction2.Operation() != "ConcreteImplementorB" {
t.Errorf("Expected ConcreteImplementorB, got %s", abstraction2.Operation())
}
}

func TestRefinedAbstraction(t *testing.T) {
implA := &ConcreteImplementorA{}
refined := NewRefinedAbstraction(implA)
expected := "RefinedAbstraction: ConcreteImplementorA"
if refined.Operation() != expected {
t.Errorf("Expected %s, got %s", expected, refined.Operation())
}
}
