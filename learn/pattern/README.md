# Go 设计模式学习项目

这个项目提供了完整的 Go 语言设计模式学习资源，包括经典设计模式和 Go 语言特有的并发模式。

## 📚 设计模式分类

### 1. [创建型模式 (Creational Patterns)](./creational/README.md)

创建型模式关注对象的创建机制，旨在以适当的方式创建对象。这些模式提供了创建对象的灵活性，使系统独立于对象的创建、组合和表示方式。

**包含的模式:**
- Singleton (单例模式)
- Factory Method (工厂方法模式)
- Abstract Factory (抽象工厂模式)
- Builder (生成器模式)
- Prototype (原型模式)

### 2. [结构型模式 (Structural Patterns)](./structural/README.md)

结构型模式关注类和对象的组合，通过继承和组合来组织接口和实现。这些模式有助于确保当系统的某部分改变时，整个结构不需要改变。

**包含的模式:**
- Adapter (适配器模式)
- Bridge (桥接模式)
- Composite (组合模式)
- Decorator (装饰器模式)
- Facade (外观模式)
- Flyweight (享元模式)
- Proxy (代理模式)

### 3. [行为型模式 (Behavioral Patterns)](./behavioral/README.md)

行为型模式关注对象之间的通信和职责分配。这些模式定义了对象之间的通信方式，使得系统更加灵活和可扩展。

**包含的模式:**
- Chain of Responsibility (责任链模式)
- Command (命令模式)
- Iterator (迭代器模式)
- Mediator (中介者模式)
- Memento (备忘录模式)
- Observer (观察者模式)
- State (状态模式)
- Strategy (策略模式)
- Template Method (模板方法模式)
- Visitor (访问者模式)

### 4. [并发模式 (Concurrency Patterns)](./concurrency/README.md)

并发模式是 Go 语言的特色之一，利用 goroutine 和 channel 实现高效的并发编程。这些模式帮助开发者构建可扩展的并发应用。

**包含的模式:**
- Active Object (活动对象模式)
- Monitor Object (监视对象模式)
- Thread Pool (线程池模式)
- Producer-Consumer (生产者-消费者模式)
- Pipeline (管道模式)

### 5. [同步模式 (Synchronization Patterns)](./synchronization/README.md)

同步模式提供了线程同步和资源访问控制的机制。这些模式确保在并发环境中正确地访问共享资源。

**包含的模式:**
- Mutex (互斥锁模式)
- Semaphore (信号量模式)
- Barrier (屏障模式)
- Read-Write Lock (读写锁模式)
- Condition Variable (条件变量模式)

## 🎯 学习路径

### 初学者路径
1. **从创建型模式开始**
   - Singleton - 理解单例的实现和线程安全
   - Factory Method - 学习对象创建的封装
   - Builder - 掌握复杂对象的构建过程

2. **学习基础结构型模式**
   - Adapter - 理解接口适配
   - Decorator - 学习动态扩展功能
   - Proxy - 掌握代理和访问控制

3. **掌握常用行为型模式**
   - Strategy - 理解算法封装
   - Observer - 学习事件通知机制
   - Template Method - 掌握算法骨架定义

### 进阶路径
1. **深入理解 Go 并发特性**
   - Producer-Consumer - 理解 channel 的使用
   - Pipeline - 学习数据流处理
   - Monitor Object - 掌握同步监视器

2. **学习高级结构型模式**
   - Bridge - 理解抽象与实现分离
   - Composite - 学习树形结构处理
   - Flyweight - 掌握内存优化技术

3. **掌握复杂行为型模式**
   - Chain of Responsibility - 理解请求链处理
   - Mediator - 学习对象解耦
   - State - 掌握状态转换管理

### 专家路径
1. **并发模式深度学习**
   - Active Object - 异步方法调用
   - Thread Pool - 协程池管理
   - 各种同步原语的组合使用

2. **设计模式组合应用**
   - 学习多个模式的协同使用
   - 理解模式之间的关系
   - 在实际项目中应用模式

## 🔧 Go 语言特有的设计模式考虑

### 1. 接口的重要性
Go 语言采用隐式接口实现，这使得许多设计模式的实现更加简洁和灵活。接口是 Go 设计模式的核心。

```go
// 不需要显式声明实现接口
type Reader interface {
    Read(p []byte) (n int, err error)
}

// 只需要实现方法即可
type MyReader struct{}

func (r MyReader) Read(p []byte) (n int, err error) {
    // 实现
    return 0, nil
}
```

### 2. 组合优于继承
Go 不支持传统的类继承，而是通过嵌入（embedding）实现代码复用，这影响了许多模式的实现方式。

```go
type Base struct {
    name string
}

type Extended struct {
    Base  // 嵌入
    extra string
}
```

### 3. 并发原语
Go 的 goroutine 和 channel 为并发模式提供了原生支持，使得并发编程更加简单和安全。

```go
// 使用 channel 实现生产者-消费者
ch := make(chan int, 10)

// 生产者
go func() {
    ch <- 1
}()

// 消费者
go func() {
    value := <-ch
}()
```

### 4. 函数是一等公民
Go 中函数可以作为参数和返回值，这使得策略模式、命令模式等行为型模式的实现更加灵活。

```go
// 函数作为参数
func Process(data []int, handler func(int) int) []int {
    result := make([]int, len(data))
    for i, v := range data {
        result[i] = handler(v)
    }
    return result
}
```

### 5. 零值可用性
Go 的零值初始化特性影响了某些模式的实现，如 Singleton 可以利用 sync.Once 实现。

```go
var (
    instance *Singleton
    once     sync.Once
)

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{}
    })
    return instance
}
```

### 6. 错误处理
Go 的显式错误处理方式影响了许多模式的实现，需要在设计时考虑错误传播。

```go
type Result struct {
    Value interface{}
    Error error
}

// 在模式实现中明确处理错误
func (f *Factory) Create(name string) (*Product, error) {
    if name == "" {
        return nil, errors.New("name cannot be empty")
    }
    // ...
}
```

## 📖 如何使用本项目

### 查看模式文档
每个设计模式都有详细的文档说明，包括：
- 模式定义和目的
- 使用场景和示例
- Go 语言特有实现
- 优缺点分析
- 相关模式对比

### 运行示例代码
```bash
# 运行特定模式的示例
cd learn/pattern/creational/singleton
go run main.go

# 运行测试
go test ./...
```

### 学习建议
1. **阅读文档** - 先理解模式的概念和使用场景
2. **查看代码** - 研究具体实现细节
3. **运行示例** - 实际执行代码，观察行为
4. **编写测试** - 尝试修改和扩展代码
5. **实践应用** - 在实际项目中应用学到的模式

## 🤝 贡献

欢迎贡献新的设计模式示例或改进现有实现！

## 📄 许可证

本项目采用开源许可证，具体请查看项目根目录的 LICENSE 文件。

## 📚 参考资源

- [Design Patterns: Elements of Reusable Object-Oriented Software](https://en.wikipedia.org/wiki/Design_Patterns) - GoF 设计模式经典书籍
- [Go 语言设计与实现](https://draveness.me/golang/) - Go 语言深度解析
- [Effective Go](https://golang.org/doc/effective_go) - Go 语言官方最佳实践
- [Go Concurrency Patterns](https://go.dev/blog/pipelines) - Go 官方并发模式博客
