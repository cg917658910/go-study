# 创建型模式 (Creational Patterns)

创建型模式关注对象的创建机制，旨在以适当的方式创建对象。这些模式提供了创建对象的灵活性，使系统独立于对象的创建、组合和表示方式。

## 📋 模式列表

### 1. [Singleton (单例模式)](./singleton/)

**目的**: 确保一个类只有一个实例，并提供一个全局访问点。

**使用场景**:
- 数据库连接池
- 配置管理器
- 日志记录器
- 缓存管理器

**Go 特有实现**: 使用 `sync.Once` 确保线程安全的延迟初始化。

**示例**: ✅ 已实现

---

### 2. [Factory Method (工厂方法模式)](./factory_method/)

**目的**: 定义一个创建对象的接口，但让子类决定实例化哪个类。

**使用场景**:
- 数据库驱动选择（MySQL, PostgreSQL, SQLite）
- 日志处理器（文件日志、控制台日志、网络日志）
- 消息队列客户端（Kafka, RabbitMQ, Redis）

**Go 特有实现**: 利用接口和函数返回不同的实现类型。

**示例**: ✅ 已实现

---

### 3. [Abstract Factory (抽象工厂模式)](./abstract_factory/)

**目的**: 提供一个接口，用于创建相关或依赖对象的家族，而不需要指定它们的具体类。

**使用场景**:
- UI 组件库（不同主题的按钮、输入框、对话框）
- 跨平台应用（Windows、Linux、MacOS 组件）
- 数据库访问层（不同数据库的连接、命令、事务对象）

**Go 特有实现**: 使用接口定义产品族，工厂返回具体实现。

**示例**: ✅ 已实现

---

### 4. [Builder (生成器模式)](./builder/)

**目的**: 将复杂对象的构建与其表示分离，使得同样的构建过程可以创建不同的表示。

**使用场景**:
- HTTP 请求构建
- SQL 查询构建
- 复杂配置对象创建
- 文档生成器（HTML、PDF、Markdown）

**Go 特有实现**: 使用方法链（Method Chaining）和可选参数模式。

**示例**: ✅ 已实现

---

### 5. [Prototype (原型模式)](./prototype/)

**目的**: 通过复制现有对象来创建新对象，而不是通过实例化类。

**使用场景**:
- 对象克隆（深拷贝、浅拷贝）
- 配置对象的复制
- 游戏对象的克隆（敌人、道具）
- 文档模板系统

**Go 特有实现**: 由于 Go 没有类继承，通常通过定义 `Clone()` 方法或使用序列化/反序列化实现。

**示例**: ⏳ 待实现

---

## 🎯 学习顺序建议

1. **Singleton** - 最简单的创建型模式，理解单例和线程安全
2. **Factory Method** - 学习对象创建的封装和多态
3. **Builder** - 掌握复杂对象的构建过程
4. **Abstract Factory** - 理解产品族的创建
5. **Prototype** - 学习对象克隆技术

## 💡 Go 语言实现要点

### 1. 使用接口而非抽象类
Go 没有抽象类，使用接口定义产品类型：

```go
type Product interface {
    Use() string
}

type ConcreteProductA struct{}

func (p *ConcreteProductA) Use() string {
    return "Product A"
}
```

### 2. 工厂函数
Go 中工厂通常是函数而非类：

```go
// 简单工厂函数
func NewProduct(productType string) Product {
    switch productType {
    case "A":
        return &ConcreteProductA{}
    case "B":
        return &ConcreteProductB{}
    default:
        return nil
    }
}
```

### 3. 函数式选项模式
Builder 模式的 Go 惯用实现：

```go
type Server struct {
    host string
    port int
}

type Option func(*Server)

func WithHost(host string) Option {
    return func(s *Server) {
        s.host = host
    }
}

func NewServer(opts ...Option) *Server {
    s := &Server{
        host: "localhost",
        port: 8080,
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}
```

### 4. 单例的线程安全实现
使用 `sync.Once`:

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

## 🔄 模式对比

| 模式 | 关注点 | 使用时机 | 复杂度 |
|------|--------|----------|--------|
| Singleton | 唯一实例 | 需要全局唯一对象 | ⭐ |
| Factory Method | 对象创建 | 不确定具体类型 | ⭐⭐ |
| Abstract Factory | 产品族创建 | 需要创建相关对象组 | ⭐⭐⭐ |
| Builder | 复杂对象构建 | 对象创建步骤复杂 | ⭐⭐ |
| Prototype | 对象克隆 | 创建成本高，需复制 | ⭐⭐ |

## 📚 相关模式

- **Singleton + Factory**: 工厂本身可以是单例
- **Abstract Factory + Singleton**: 抽象工厂通常实现为单例
- **Builder + Composite**: Builder 可以构建 Composite 对象
- **Prototype + Abstract Factory**: 可以用原型替代工厂创建对象

## ⚠️ 常见陷阱

1. **过度使用 Singleton**: 可能导致全局状态和测试困难
2. **工厂过于复杂**: 简单场景不需要工厂模式
3. **Builder 过度设计**: 对于简单对象，直接构造更清晰
4. **忽略线程安全**: 在并发环境中必须考虑线程安全

## 🎓 最佳实践

1. **优先使用简单工厂函数**: 除非确实需要复杂的工厂模式
2. **使用接口定义产品**: 保持灵活性和可测试性
3. **考虑零值可用性**: Go 的零值特性可以简化对象创建
4. **合理使用选项模式**: 对于配置较多的对象，选项模式更优雅
5. **注意内存和性能**: Prototype 模式中注意深拷贝的性能开销
