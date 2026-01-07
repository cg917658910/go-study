# Adapter Pattern (适配器模式)

## 定义
适配器模式将一个类的接口转换成客户期望的另一个接口,使原本接口不兼容的类可以一起工作。

## 目的
- 接口转换
- 复用现有类
- 提高兼容性

## 使用场景
- 集成第三方库(接口不匹配)
- 数据格式转换(JSON、XML、Protobuf)
- 旧系统与新系统对接
- 不同日志框架的统一接口

## Go 特有实现
使用接口适配和嵌入:

```go
type Target interface {
    Request() string
}

type Adaptee struct{}

func (a *Adaptee) SpecificRequest() string {
    return "Specific request"
}

type Adapter struct {
    adaptee *Adaptee
}

func (a *Adapter) Request() string {
    return a.adaptee.SpecificRequest()
}
```

## 优点
1. **解耦** - 客户端与被适配者解耦
2. **复用** - 复用现有功能
3. **灵活性** - 适配多个被适配者

## 缺点
1. **复杂度** - 增加系统复杂度
2. **性能** - 可能有性能开销
