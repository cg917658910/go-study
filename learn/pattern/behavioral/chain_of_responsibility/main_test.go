package main

import "testing"

func TestChainOfResponsibility(t *testing.T) {
	handler1 := &ConcreteHandler1{}
	handler2 := &ConcreteHandler2{}

	handler1.SetNext(handler2)

	// 测试不会panic即可
	handler1.Handle("Handler1")
	handler1.Handle("Handler2")
}
