# Observer Pattern (观察者模式)

## 定义
观察者模式定义对象间的一种一对多的依赖关系,当一个对象的状态发生改变时,所有依赖于它的对象都得到通知。

## 目的
- 建立一对多依赖
- 自动通知更新
- 松耦合设计

## 使用场景
- 事件系统
- 发布-订阅系统
- 数据绑定
- MVC 架构中的模型-视图更新

## Go 特有实现
使用 channel 实现异步通知:

```go
type Subject struct {
    observers []Observer
    updates   chan interface{}
}

func (s *Subject) Notify(data interface{}) {
    s.updates <- data
}

func (s *Subject) Start() {
    go func() {
        for data := range s.updates {
            for _, observer := range s.observers {
                observer.Update(data)
            }
        }
    }()
}
```

## 优点
1. **松耦合** - 主题和观察者松耦合
2. **动态关系** - 运行时建立关系
3. **广播通信** - 支持广播通信

## 缺点
1. **性能问题** - 大量观察者影响性能
2. **内存泄漏** - 未移除的观察者可能导致内存泄漏
