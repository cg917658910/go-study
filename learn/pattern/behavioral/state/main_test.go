package main

import "testing"

func TestState(t *testing.T) {
	context := &Context{state: &ConcreteStateA{}}
	context.Request()
}
