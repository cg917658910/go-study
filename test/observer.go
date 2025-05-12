package test

import "fmt"

// 观察者模式
type Observer interface {
	Update(news string)
}

type AObserver struct {
	Name string
}

type BObserver struct {
	Name string
}

func (o *AObserver) Update(news string) {
	fmt.Printf("%s 接收到新消息：%s\n", o.Name, news)
}
func (o *BObserver) Update(news string) {
	fmt.Printf("%s 接收到新消息：%s\n", o.Name, news)
}

type Subject interface {
	Attach(observer Observer)
	Detach(observer Observer)
	Notify(news string)
}

type NewsSubject struct {
	observers []Observer
}

func (s *NewsSubject) Attach(observer Observer) {
	s.observers = append(s.observers, observer)
}

func (s *NewsSubject) Detach(observer Observer) {
	for i, o := range s.observers {
		if o == observer {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

func (s *NewsSubject) Notify(news string) {
	for _, o := range s.observers {
		o.Update(news)
	}
}

func doObserver() {
	subject := &NewsSubject{}
	a := &AObserver{Name: "A"}
	b := &BObserver{Name: "B"}
	subject.Attach(a)
	subject.Attach(b)

	subject.Notify("今日天气��转多云")
	subject.Detach(a)
	subject.Notify("今日天气有��")
}
