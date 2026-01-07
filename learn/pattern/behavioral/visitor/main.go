package main

import "fmt"

type Visitor interface {
VisitConcreteElementA(element *ConcreteElementA)
VisitConcreteElementB(element *ConcreteElementB)
}

type Element interface {
Accept(visitor Visitor)
}

type ConcreteElementA struct {
name string
}

func (e *ConcreteElementA) Accept(visitor Visitor) {
visitor.VisitConcreteElementA(e)
}

type ConcreteElementB struct {
value int
}

func (e *ConcreteElementB) Accept(visitor Visitor) {
visitor.VisitConcreteElementB(e)
}

type ConcreteVisitor struct{}

func (v *ConcreteVisitor) VisitConcreteElementA(element *ConcreteElementA) {
fmt.Printf("Visiting ElementA: %s\n", element.name)
}

func (v *ConcreteVisitor) VisitConcreteElementB(element *ConcreteElementB) {
fmt.Printf("Visiting ElementB: %d\n", element.value)
}

func main() {
elements := []Element{
&ConcreteElementA{name: "A"},
&ConcreteElementB{value: 42},
}

visitor := &ConcreteVisitor{}
for _, element := range elements {
element.Accept(visitor)
}
}
