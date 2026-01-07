package main

import "fmt"

type Mediator interface {
Notify(sender string, event string)
}

type ConcreteMediator struct {
component1 *Component1
component2 *Component2
}

func (m *ConcreteMediator) Notify(sender string, event string) {
fmt.Printf("Mediator reacting to %s with event %s\n", sender, event)
}

type Component1 struct {
mediator Mediator
}

func (c *Component1) DoA() {
fmt.Println("Component1 does A")
c.mediator.Notify("Component1", "A")
}

type Component2 struct {
mediator Mediator
}

func (c *Component2) DoB() {
fmt.Println("Component2 does B")
c.mediator.Notify("Component2", "B")
}

func main() {
mediator := &ConcreteMediator{}
c1 := &Component1{mediator: mediator}
c2 := &Component2{mediator: mediator}
mediator.component1 = c1
mediator.component2 = c2

c1.DoA()
c2.DoB()
}
