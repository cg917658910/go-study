package main

import "fmt"

// Handler 处理器接口
type Handler interface {
SetNext(handler Handler)
Handle(request string)
}

// BaseHandler 基础处理器
type BaseHandler struct {
next Handler
}

func (h *BaseHandler) SetNext(handler Handler) {
h.next = handler
}

func (h *BaseHandler) Handle(request string) {
if h.next != nil {
h.next.Handle(request)
}
}

// ConcreteHandler1 具体处理器1
type ConcreteHandler1 struct {
BaseHandler
}

func (h *ConcreteHandler1) Handle(request string) {
if request == "Handler1" {
fmt.Println("ConcreteHandler1 handled the request")
} else if h.next != nil {
h.next.Handle(request)
}
}

// ConcreteHandler2 具体处理器2
type ConcreteHandler2 struct {
BaseHandler
}

func (h *ConcreteHandler2) Handle(request string) {
if request == "Handler2" {
fmt.Println("ConcreteHandler2 handled the request")
} else if h.next != nil {
h.next.Handle(request)
}
}

func main() {
handler1 := &ConcreteHandler1{}
handler2 := &ConcreteHandler2{}

handler1.SetNext(handler2)

handler1.Handle("Handler1")
handler1.Handle("Handler2")
handler1.Handle("Handler3")
}
