package main

import (
"fmt"
"sync"
)

// Flyweight 享元接口
type Flyweight interface {
Operation(extrinsicState string)
}

// ConcreteFlyweight 具体享元
type ConcreteFlyweight struct {
intrinsicState string
}

func (c *ConcreteFlyweight) Operation(extrinsicState string) {
fmt.Printf("IntrinsicState: %s, ExtrinsicState: %s\n", c.intrinsicState, extrinsicState)
}

// FlyweightFactory 享元工厂
type FlyweightFactory struct {
flyweights map[string]Flyweight
mu         sync.RWMutex
}

func NewFlyweightFactory() *FlyweightFactory {
return &FlyweightFactory{
flyweights: make(map[string]Flyweight),
}
}

func (f *FlyweightFactory) GetFlyweight(key string) Flyweight {
f.mu.RLock()
flyweight, exists := f.flyweights[key]
f.mu.RUnlock()

if !exists {
f.mu.Lock()
flyweight = &ConcreteFlyweight{intrinsicState: key}
f.flyweights[key] = flyweight
f.mu.Unlock()
fmt.Printf("Creating new flyweight for key: %s\n", key)
}

return flyweight
}

func (f *FlyweightFactory) GetFlyweightCount() int {
f.mu.RLock()
defer f.mu.RUnlock()
return len(f.flyweights)
}

func main() {
factory := NewFlyweightFactory()

fw1 := factory.GetFlyweight("A")
fw2 := factory.GetFlyweight("B")
fw3 := factory.GetFlyweight("A")

fw1.Operation("First")
fw2.Operation("Second")
fw3.Operation("Third")

fmt.Printf("Total flyweights created: %d\n", factory.GetFlyweightCount())
}
