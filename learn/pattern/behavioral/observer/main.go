package main

// Observer Pattern 观察者模式 行为型设计模式
// 定义对象间的一种一对多的依赖关系 使得每当一个对象改变状态 则所有依赖于它的对象都会得到通知并被自动更新
// 适用于一个对象的改变需要同时影响其他对象 而且不知道有多少对象需要被影响的场景

type Subject interface {
	RegisterObserver(observer Observer)
	RemoveObserver(observer Observer)
	NotifyObservers()
}

type Observer interface {
	Update(message string)
}

type MessageSubject struct {
	observers []Observer
	message   string
}

func (s *MessageSubject) RegisterObserver(observer Observer) {
	s.observers = append(s.observers, observer)
}

func (s *MessageSubject) RemoveObserver(observer Observer) {
	for i, obs := range s.observers {
		if obs == observer {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

func (s *MessageSubject) NotifyObservers() {
	for _, observer := range s.observers {
		observer.Update(s.message)
	}
}

func (s *MessageSubject) SetMessage(message string) {
	s.message = message
	s.NotifyObservers()
}

type MessageObserver struct {
	name string
}

func NewMessageObserver(name string) *MessageObserver {
	return &MessageObserver{name: name}
}

func (o *MessageObserver) Update(message string) {
	println(o.name, "received message:", message)
}

func main() {
	subject := &MessageSubject{}

	observer1 := NewMessageObserver("Observer 1")
	observer2 := NewMessageObserver("Observer 2")

	subject.RegisterObserver(observer1)
	subject.RegisterObserver(observer2)

	subject.SetMessage("Hello Observers!")

	subject.RemoveObserver(observer1)

	subject.SetMessage("Observer 1 has been removed.")
}

// 优点: 观察者和主题之间的耦合度低 观察者可以动态注册和注销 便于扩展
// 缺点: 可能导致过多的通知影响性能 观察者之间可能存在依赖关系 导致复杂性增加
// 适用场景: 一个对象的改变需要同时影响其他对象 而且不知道有多少对象需要被影响
// 对比发布-订阅模式: 发布-订阅模式通过消息队列实现解耦 观察者模式直接在对象间建立依赖关系
// 举例: GUI事件处理系统 当按钮被点击时 所有注册的监听器都会收到通知
// Web领域例子: 实时通知系统 当有新消息时 所有在线用户都会收到通知
// 数据库领域例子: 数据库触发器 当表中的数据被修改时 相关的触发器会被自动执行
// 函数编程领域例子: 响应式编程 当数据源发生变化时 相关的计算会自动更新

// 异步版观察者模式 可以使用消息队列或事件总线来实现异步通知 进一步降低耦合度 提高系统的可扩展性
