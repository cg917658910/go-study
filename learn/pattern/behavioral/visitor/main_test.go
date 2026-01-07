package main

import "testing"

func TestVisitor(t *testing.T) {
	elements := []Element{
		&ConcreteElementA{name: "A"},
		&ConcreteElementB{value: 42},
	}

	visitor := &ConcreteVisitor{}
	for _, element := range elements {
		element.Accept(visitor)
	}
}
