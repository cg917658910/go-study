# 结构型模式 (Structural Patterns)

结构型模式关注类和对象的组合，通过继承和组合来组织接口和实现。这些模式有助于确保当系统的某部分改变时，整个结构不需要改变。

## 📋 模式列表

### 1. [Adapter (适配器模式)](./adapter/)

**目的**: 将一个类的接口转换成客户期望的另一个接口，使原本接口不兼容的类可以一起工作。

**使用场景**:
- 集成第三方库（接口不匹配）
- 数据格式转换（JSON、XML、Protobuf）
- 旧系统与新系统对接
- 不同日志框架的统一接口

**Go 特有实现**: 使用接口适配和嵌入实现。

**示例**: ✅ 已实现

---

### 2. [Bridge (桥接模式)](./bridge/)

**目的**: 将抽象部分与它的实现部分分离，使它们都可以独立地变化。

**使用场景**:
- 跨平台应用（业务逻辑与平台实现分离）
- 数据库驱动（SQL 接口与具体数据库分离）
- 图形渲染（图形对象与渲染引擎分离）
- 消息发送（消息类型与发送方式分离）

**Go 特有实现**: 通过接口组合实现抽象与实现的分离。

**示例**: ⏳ 待实现

---

### 3. [Composite (组合模式)](./composite/)

**目的**: 将对象组合成树形结构以表示"部分-整体"的层次结构，使客户端对单个对象和组合对象的使用具有一致性。

**使用场景**:
- 文件系统（文件和文件夹）
- UI 组件树（容器和控件）
- 组织架构（部门和员工）
- 菜单系统（菜单和菜单项）

**Go 特有实现**: 使用统一接口处理叶子节点和组合节点。

**示例**: ✅ 已实现

---

### 4. [Decorator (装饰器模式)](./decorator/)

**目的**: 动态地给一个对象添加一些额外的职责，就增加功能来说，装饰器模式相比生成子类更为灵活。

**使用场景**:
- HTTP 中间件（日志、认证、压缩）
- I/O 流装饰（缓冲、加密、压缩）
- 缓存装饰（为服务添加缓存层）
- 功能增强（为基础服务添加新功能）

**Go 特有实现**: 利用接口和函数闭包实现装饰。

**示例**: ✅ 已实现

---

### 5. [Facade (外观模式)](./facade/)

**目的**: 为子系统中的一组接口提供一个一致的界面，使子系统更容易使用。

**使用场景**:
- 简化复杂系统的接口
- 第三方库的封装
- 微服务聚合接口
- 系统启动和关闭流程

**Go 特有实现**: 提供简单的函数或结构体封装复杂操作。

**示例**: ✅ 已实现

---

### 6. [Flyweight (享元模式)](./flyweight/)

**目的**: 运用共享技术有效地支持大量细粒度的对象，减少内存使用。

**使用场景**:
- 字符串池（String interning）
- 对象池（数据库连接池、goroutine 池）
- 缓存系统
- 游戏中的大量相似对象（子弹、粒子）

**Go 特有实现**: 使用 map 或 sync.Pool 实现对象共享。

**示例**: ⏳ 待实现

---

### 7. [Proxy (代理模式)](./proxy/)

**目的**: 为其他对象提供一种代理以控制对这个对象的访问。

**使用场景**:
- 远程代理（RPC 调用）
- 虚拟代理（延迟加载）
- 保护代理（访问控制）
- 缓存代理（结果缓存）

**Go 特有实现**: 使用接口实现透明代理。

**示例**: ✅ 已实现

---

## 🎯 学习顺序建议

1. **Adapter** - 最实用的模式，理解接口转换
2. **Decorator** - 学习动态功能扩展，理解中间件概念
3. **Facade** - 掌握系统简化技术
4. **Proxy** - 理解代理和访问控制
5. **Composite** - 学习树形结构处理
6. **Bridge** - 理解抽象与实现分离
7. **Flyweight** - 学习内存优化技术

## 💡 Go 语言实现要点

### 1. 适配器模式
使用接口适配和嵌入：

```go
// 目标接口
type Target interface {
    Request() string
}

// 被适配的类
type Adaptee struct{}

func (a *Adaptee) SpecificRequest() string {
    return "Specific request"
}

// 适配器
type Adapter struct {
    adaptee *Adaptee
}

func (a *Adapter) Request() string {
    return a.adaptee.SpecificRequest()
}
```

### 2. 装饰器模式
Go 中常用中间件模式：

```go
type Handler func(http.ResponseWriter, *http.Request)

func LoggingMiddleware(next Handler) Handler {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Request: %s", r.URL.Path)
        next(w, r)
    }
}
```

### 3. 组合模式
统一接口处理节点：

```go
type Component interface {
    Operation() string
}

type Leaf struct {
    name string
}

func (l *Leaf) Operation() string {
    return l.name
}

type Composite struct {
    children []Component
}

func (c *Composite) Operation() string {
    result := ""
    for _, child := range c.children {
        result += child.Operation()
    }
    return result
}
```

### 4. 享元模式
使用 sync.Pool 或 map:

```go
var flyweightPool = &sync.Pool{
    New: func() interface{} {
        return &Flyweight{}
    },
}

func GetFlyweight() *Flyweight {
    return flyweightPool.Get().(*Flyweight)
}

func PutFlyweight(f *Flyweight) {
    flyweightPool.Put(f)
}
```

## 🔄 模式对比

| 模式 | 关注点 | 主要用途 | 复杂度 |
|------|--------|----------|--------|
| Adapter | 接口转换 | 不兼容接口对接 | ⭐⭐ |
| Bridge | 抽象与实现分离 | 多维度变化 | ⭐⭐⭐ |
| Composite | 树形结构 | 部分-整体层次 | ⭐⭐⭐ |
| Decorator | 动态扩展 | 功能增强 | ⭐⭐ |
| Facade | 简化接口 | 子系统封装 | ⭐ |
| Flyweight | 对象共享 | 减少内存 | ⭐⭐⭐ |
| Proxy | 访问控制 | 代理访问 | ⭐⭐ |

## 📚 相关模式

- **Adapter vs Bridge**: Adapter 改变接口，Bridge 分离抽象和实现
- **Decorator vs Proxy**: Decorator 增强功能，Proxy 控制访问
- **Composite vs Decorator**: Composite 关注结构，Decorator 关注功能
- **Facade vs Adapter**: Facade 简化接口，Adapter 转换接口

## ⚠️ 常见陷阱

1. **过度装饰**: 太多层装饰器会降低性能和可读性
2. **忽略组合的性能**: 大型组合结构的遍历可能很耗时
3. **代理的透明性**: 确保代理和真实对象行为一致
4. **享元的线程安全**: 共享对象必须考虑并发访问

## 🎓 最佳实践

1. **使用接口实现灵活性**: 所有结构型模式都应基于接口
2. **保持装饰器的单一职责**: 每个装饰器只做一件事
3. **合理使用嵌入**: Go 的嵌入可以简化适配器和代理的实现
4. **注意内存管理**: Flyweight 和对象池需要正确管理对象生命周期
5. **文档化代理行为**: 清楚说明代理与真实对象的差异
