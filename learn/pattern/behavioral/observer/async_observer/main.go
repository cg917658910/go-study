package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Subject 接口
type Subject interface {
	RegisterObserver(observer Observer)
	RemoveObserver(observer Observer)
	NotifyObservers()
}

// Observer 接口
type Observer interface {
	Update(message string)
}

// 异步版 MessageSubject
type AsyncMessageSubject struct {
	observers []Observer
	message   string
	mu        sync.Mutex
}

func (s *AsyncMessageSubject) RegisterObserver(observer Observer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.observers = append(s.observers, observer)
}

func (s *AsyncMessageSubject) RemoveObserver(observer Observer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, obs := range s.observers {
		if obs == observer {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

func (s *AsyncMessageSubject) NotifyObserversWithTimeout(timeout time.Duration) {
	s.mu.Lock()
	observers := s.observers
	s.mu.Unlock()

	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for _, observer := range observers {
		wg.Add(1)
		go func(obs Observer) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				fmt.Println("Notification timed out for an observer")
			default:
				obs.Update(s.message)
			}
		}(observer)
	}

	// 等待所有通知完成或超时
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("All notifications completed")
	case <-ctx.Done():
		fmt.Println("Notification process timed out")
	}
}

func (s *AsyncMessageSubject) SetMessage(message string) {
	s.message = message
}

type MessageObserver struct {
	name string
}

func NewMessageObserver(name string) *MessageObserver {
	return &MessageObserver{name: name}
}

func (o *MessageObserver) Update(message string) {
	fmt.Printf("%s received message: %s\n", o.name, message)
	time.Sleep(1 * time.Second) // 模拟处理时间
}

func main() {
	subject := &AsyncMessageSubject{}

	observer1 := NewMessageObserver("Observer 1")
	observer2 := NewMessageObserver("Observer 2")
	observer3 := NewMessageObserver("Observer 3")

	subject.RegisterObserver(observer1)
	subject.RegisterObserver(observer2)
	subject.RegisterObserver(observer3)

	subject.SetMessage("Hello Observers!")
	subject.NotifyObserversWithTimeout(2 * time.Second) // 设置超时时间为 2 秒
}
