package main

import "testing"

func TestMediator(t *testing.T) {
	mediator := &ConcreteMediator{}
	c1 := &Component1{mediator: mediator}
	c1.DoA()
}
