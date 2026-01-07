# 行为型模式 (Behavioral Patterns)

行为型模式关注对象之间的通信和职责分配。这些模式定义了对象之间的通信方式，使得系统更加灵活和可扩展。

## 📋 模式列表

### 1. [Chain of Responsibility (责任链模式)](./chain_of_responsibility/)

**目的**: 使多个对象都有机会处理请求，从而避免请求的发送者和接收者之间的耦合关系。

**使用场景**:
- HTTP 中间件链
- 事件处理链
- 审批流程
- 日志处理（按级别过滤）

**Go 特有实现**: 使用接口和链表结构，或者函数链。

**示例**: ⏳ 待实现

---

### 2. [Command (命令模式)](./command/)

**目的**: 将请求封装成对象，从而使你可用不同的请求对客户进行参数化，对请求排队或记录请求日志。

**使用场景**:
- 任务队列
- 撤销/重做操作
- 宏命令（批处理）
- 事务系统

**Go 特有实现**: 使用接口和函数类型实现命令。

**示例**: ✅ 已实现

---

### 3. [Iterator (迭代器模式)](./iterator/)

**目的**: 提供一种方法顺序访问一个聚合对象中的各个元素，而又不暴露该对象的内部表示。

**使用场景**:
- 集合遍历
- 数据库结果集
- 文件行读取
- 树/图的遍历

**Go 特有实现**: 使用 channel 或者实现 Next() 方法。

**示例**: ⏳ 待实现

---

### 4. [Mediator (中介者模式)](./mediator/)

**目的**: 用一个中介对象来封装一系列的对象交互，使各对象不需要显式地相互引用。

**使用场景**:
- 聊天室（用户之间的消息中介）
- GUI 组件交互
- 微服务间的消息总线
- 事件总线系统

**Go 特有实现**: 使用 channel 作为通信中介。

**示例**: ⏳ 待实现

---

### 5. [Memento (备忘录模式)](./memento/)

**目的**: 在不破坏封装性的前提下，捕获一个对象的内部状态，并在该对象之外保存这个状态。

**使用场景**:
- 编辑器的撤销/重做
- 游戏存档
- 数据库事务回滚
- 配置快照

**Go 特有实现**: 使用结构体快照或序列化。

**示例**: ⏳ 待实现

---

### 6. [Observer (观察者模式)](./observer/)

**目的**: 定义对象间的一种一对多的依赖关系，当一个对象的状态发生改变时，所有依赖于它的对象都得到通知。

**使用场景**:
- 事件系统
- 发布-订阅系统
- 数据绑定
- MVC 架构中的模型-视图更新

**Go 特有实现**: 使用 channel 实现异步通知。

**示例**: ✅ 已实现

---

### 7. [State (状态模式)](./state/)

**目的**: 允许对象在内部状态改变时改变它的行为，对象看起来好像修改了它的类。

**使用场景**:
- 工作流引擎
- TCP 连接状态
- 订单状态管理
- 游戏角色状态

**Go 特有实现**: 使用接口表示状态，上下文持有当前状态。

**示例**: ⏳ 待实现

---

### 8. [Strategy (策略模式)](./strategy/)

**目的**: 定义一系列算法，把它们一个个封装起来，并且使它们可相互替换。

**使用场景**:
- 排序算法选择
- 压缩算法选择
- 支付方式选择
- 数据验证策略

**Go 特有实现**: 使用接口或函数类型实现策略。

**示例**: ✅ 已实现

---

### 9. [Template Method (模板方法模式)](./template_method/)

**目的**: 定义一个操作中的算法骨架，而将一些步骤延迟到子类中。

**使用场景**:
- 框架和库的钩子方法
- 数据处理流程
- 测试框架
- 算法框架

**Go 特有实现**: 使用接口和嵌入实现模板方法。

**示例**: ✅ 已实现

---

### 10. [Visitor (访问者模式)](./visitor/)

**目的**: 表示一个作用于某对象结构中的各元素的操作，使你可以在不改变各元素的类的前提下定义作用于这些元素的新操作。

**使用场景**:
- AST（抽象语法树）遍历
- 对象结构的序列化
- 报表生成
- 编译器中的语义分析

**Go 特有实现**: 使用接口和双分派技术。

**示例**: ⏳ 待实现

---

## 🎯 学习顺序建议

1. **Strategy** - 最简单实用，理解算法封装
2. **Observer** - 学习事件通知机制，理解发布订阅
3. **Command** - 掌握请求封装
4. **Template Method** - 理解算法骨架定义
5. **State** - 学习状态转换管理
6. **Chain of Responsibility** - 理解责任链处理
7. **Iterator** - 掌握集合遍历
8. **Mediator** - 学习对象解耦
9. **Memento** - 理解状态保存和恢复
10. **Visitor** - 学习双分派技术

## 💡 Go 语言实现要点

### 1. 策略模式
使用接口或函数类型：

```go
// 接口方式
type Strategy interface {
    Execute(a, b int) int
}

// 函数类型方式
type StrategyFunc func(a, b int) int

type Context struct {
    strategy StrategyFunc
}

func (c *Context) Execute(a, b int) int {
    return c.strategy(a, b)
}
```

### 2. 观察者模式
使用 channel 实现异步通知：

```go
type Observer interface {
    Update(data interface{})
}

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

### 3. 状态模式
使用接口表示状态：

```go
type State interface {
    Handle(context *Context)
}

type Context struct {
    state State
}

func (c *Context) SetState(state State) {
    c.state = state
}

func (c *Context) Request() {
    c.state.Handle(c)
}
```

### 4. 命令模式
函数作为命令：

```go
type Command interface {
    Execute()
}

// 或者使用函数类型
type CommandFunc func()

func (f CommandFunc) Execute() {
    f()
}
```

## 🔄 模式对比

| 模式 | 关注点 | 主要用途 | 复杂度 |
|------|--------|----------|--------|
| Chain of Responsibility | 请求传递 | 解耦发送者和接收者 | ⭐⭐ |
| Command | 请求封装 | 操作参数化、排队 | ⭐⭐ |
| Iterator | 遍历访问 | 集合遍历 | ⭐⭐ |
| Mediator | 对象交互 | 减少对象间耦合 | ⭐⭐⭐ |
| Memento | 状态保存 | 撤销/恢复 | ⭐⭐ |
| Observer | 状态变化通知 | 一对多依赖 | ⭐⭐ |
| State | 状态切换 | 状态相关行为 | ⭐⭐⭐ |
| Strategy | 算法选择 | 算法可替换 | ⭐ |
| Template Method | 算法骨架 | 固定流程变化步骤 | ⭐⭐ |
| Visitor | 操作与结构分离 | 新操作添加 | ⭐⭐⭐⭐ |

## 📚 相关模式

- **Strategy vs State**: Strategy 关注算法替换，State 关注状态转换
- **Observer vs Mediator**: Observer 一对多，Mediator 多对多
- **Command vs Strategy**: Command 封装请求，Strategy 封装算法
- **Chain of Responsibility vs Decorator**: 都是链式结构，前者可中断，后者全部执行
- **Iterator vs Visitor**: Iterator 遍历，Visitor 操作

## ⚠️ 常见陷阱

1. **观察者过多**: 太多观察者会影响性能
2. **责任链过长**: 可能导致性能问题和调试困难
3. **状态爆炸**: 状态过多会增加复杂度
4. **命令队列内存**: 大量命令对象可能占用大量内存
5. **访问者的脆弱性**: 添加新元素类型需要修改所有访问者

## 🎓 最佳实践

1. **优先使用函数**: Go 中函数是一等公民，很多模式可用函数简化
2. **利用 channel**: 观察者、中介者等模式可用 channel 实现
3. **接口最小化**: 定义小而精的接口
4. **考虑并发安全**: 多个 goroutine 访问时需要同步
5. **文档化状态转换**: 状态模式要清楚文档化所有可能的状态转换
