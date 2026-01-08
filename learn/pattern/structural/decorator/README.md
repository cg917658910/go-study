# Decorator Pattern (装饰器模式)

## 定义
装饰器模式动态地给一个对象添加一些额外的职责,就增加功能来说,装饰器模式相比生成子类更为灵活。

## 目的
- 动态添加功能
- 避免子类爆炸
- 遵循开闭原则

## 使用场景
- HTTP 中间件(日志、认证、压缩)
- I/O 流装饰(缓冲、加密、压缩)
- 缓存装饰(为服务添加缓存层)
- 功能增强(为基础服务添加新功能)

## Go 特有实现
使用接口和函数闭包:

```go
type Handler func(http.ResponseWriter, *http.Request)

func LoggingDecorator(next Handler) Handler {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Request: %s", r.URL.Path)
        next(w, r)
    }
}
```

## 优点
1. **灵活性** - 动态组合功能
2. **单一职责** - 每个装饰器职责单一
3. **避免子类** - 不需要大量子类

## 缺点
1. **调试困难** - 多层装饰难以调试
2. **顺序敏感** - 装饰顺序可能影响结果
