package main

import "fmt"

type State interface {
	Handle(context *Context)
}

type Context struct {
	state State
}

func (c *Context) SetState(state State) {
	c.state = state
}

func (c *Context) Request() {
	c.state.Handle(c)
}

type ConcreteStateA struct{}

func (s *ConcreteStateA) Handle(context *Context) {
	fmt.Println("State A handling")
	context.SetState(&ConcreteStateB{})
}

type ConcreteStateB struct{}

func (s *ConcreteStateB) Handle(context *Context) {
	fmt.Println("State B handling")
	context.SetState(&ConcreteStateA{})
}

func main() {
	context := &Context{state: &ConcreteStateA{}}
	context.Request()
	context.Request()
}
