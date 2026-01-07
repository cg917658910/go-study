package main

import "fmt"

// Implementor 实现接口
type Implementor interface {
OperationImpl() string
}

// ConcreteImplementorA 具体实现A
type ConcreteImplementorA struct{}

func (c *ConcreteImplementorA) OperationImpl() string {
return "ConcreteImplementorA"
}

// ConcreteImplementorB 具体实现B
type ConcreteImplementorB struct{}

func (c *ConcreteImplementorB) OperationImpl() string {
return "ConcreteImplementorB"
}

// Abstraction 抽象类
type Abstraction struct {
implementor Implementor
}

func NewAbstraction(impl Implementor) *Abstraction {
return &Abstraction{implementor: impl}
}

func (a *Abstraction) Operation() string {
return a.implementor.OperationImpl()
}

// RefinedAbstraction 扩展抽象类
type RefinedAbstraction struct {
*Abstraction
}

func NewRefinedAbstraction(impl Implementor) *RefinedAbstraction {
return &RefinedAbstraction{
Abstraction: NewAbstraction(impl),
}
}

func (r *RefinedAbstraction) Operation() string {
return "RefinedAbstraction: " + r.Abstraction.Operation()
}

func main() {
implA := &ConcreteImplementorA{}
implB := &ConcreteImplementorB{}

abstraction1 := NewAbstraction(implA)
fmt.Println(abstraction1.Operation())

abstraction2 := NewAbstraction(implB)
fmt.Println(abstraction2.Operation())

refined1 := NewRefinedAbstraction(implA)
fmt.Println(refined1.Operation())

refined2 := NewRefinedAbstraction(implB)
fmt.Println(refined2.Operation())
}
