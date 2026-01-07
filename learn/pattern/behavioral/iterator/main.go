package main

import "fmt"

// Iterator 迭代器接口
type Iterator interface {
HasNext() bool
Next() interface{}
}

// Collection 集合接口
type Collection interface {
CreateIterator() Iterator
}

// ConcreteIterator 具体迭代器
type ConcreteIterator struct {
items []interface{}
index int
}

func (i *ConcreteIterator) HasNext() bool {
return i.index < len(i.items)
}

func (i *ConcreteIterator) Next() interface{} {
if i.HasNext() {
item := i.items[i.index]
i.index++
return item
}
return nil
}

// ConcreteCollection 具体集合
type ConcreteCollection struct {
items []interface{}
}

func (c *ConcreteCollection) CreateIterator() Iterator {
return &ConcreteIterator{
items: c.items,
index: 0,
}
}

func main() {
collection := &ConcreteCollection{
items: []interface{}{1, 2, 3, 4, 5},
}

iterator := collection.CreateIterator()
for iterator.HasNext() {
fmt.Println(iterator.Next())
}
}
