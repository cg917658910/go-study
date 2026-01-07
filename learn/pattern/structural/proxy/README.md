# Proxy Pattern (代理模式)

## 定义
代理模式为其他对象提供一种代理以控制对这个对象的访问。

## 目的
- 控制访问
- 延迟加载
- 增加额外功能

## 使用场景
- 远程代理(RPC 调用)
- 虚拟代理(延迟加载)
- 保护代理(访问控制)
- 缓存代理(结果缓存)

## Go 特有实现
使用接口实现透明代理:

```go
type Subject interface {
    Request() string
}

type RealSubject struct{}

func (r *RealSubject) Request() string {
    return "RealSubject"
}

type Proxy struct {
    realSubject *RealSubject
}

func (p *Proxy) Request() string {
    // 访问控制、缓存等
    return p.realSubject.Request()
}
```

## 优点
1. **控制访问** - 可以控制对真实对象的访问
2. **延迟初始化** - 提高性能
3. **额外功能** - 添加日志、缓存等

## 缺点
1. **响应延迟** - 可能增加响应时间
2. **复杂度** - 增加系统复杂度
