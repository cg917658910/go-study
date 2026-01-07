package main

import "fmt"

type Memento struct {
	state string
}

type Originator struct {
	state string
}

func (o *Originator) SetState(state string) {
	o.state = state
}

func (o *Originator) GetState() string {
	return o.state
}

func (o *Originator) SaveToMemento() *Memento {
	return &Memento{state: o.state}
}

func (o *Originator) RestoreFromMemento(memento *Memento) {
	o.state = memento.state
}

type Caretaker struct {
	mementos []*Memento
}

func (c *Caretaker) AddMemento(memento *Memento) {
	c.mementos = append(c.mementos, memento)
}

func (c *Caretaker) GetMemento(index int) *Memento {
	return c.mementos[index]
}

func main() {
	originator := &Originator{}
	caretaker := &Caretaker{}

	originator.SetState("State1")
	caretaker.AddMemento(originator.SaveToMemento())

	originator.SetState("State2")
	caretaker.AddMemento(originator.SaveToMemento())

	fmt.Println("Current:", originator.GetState())
	originator.RestoreFromMemento(caretaker.GetMemento(0))
	fmt.Println("Restored:", originator.GetState())
}
