# Singleton Pattern (单例模式)

## 定义
单例模式确保一个类只有一个实例,并提供一个全局访问点。

## 目的
- 控制实例数量
- 节省系统资源
- 提供全局访问点

## 使用场景
- 数据库连接池
- 配置管理器
- 日志记录器
- 缓存管理器

## Go 特有实现
使用 `sync.Once` 确保线程安全的延迟初始化:

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

## 优点
1. **唯一实例** - 确保只有一个实例
2. **全局访问** - 提供全局访问点
3. **延迟初始化** - 第一次使用时才创建

## 缺点
1. **全局状态** - 可能导致测试困难
2. **违反单一职责** - 既管理实例又执行业务逻辑
3. **并发问题** - 需要考虑线程安全
