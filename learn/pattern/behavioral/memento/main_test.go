package main

import "testing"

func TestMemento(t *testing.T) {
originator := &Originator{}
caretaker := &Caretaker{}

originator.SetState("State1")
caretaker.AddMemento(originator.SaveToMemento())

originator.SetState("State2")

originator.RestoreFromMemento(caretaker.GetMemento(0))

if originator.GetState() != "State1" {
t.Errorf("Expected State1, got %s", originator.GetState())
}
}
