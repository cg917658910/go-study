# Factory Method Pattern (工厂方法模式)

## 定义
工厂方法模式定义一个创建对象的接口,但让子类决定实例化哪个类。

## 目的
- 延迟实例化到子类
- 封装对象创建逻辑
- 提高代码可扩展性

## 使用场景
- 数据库驱动选择(MySQL, PostgreSQL, SQLite)
- 日志处理器(文件日志、控制台日志、网络日志)
- 消息队列客户端(Kafka, RabbitMQ, Redis)

## Go 特有实现
使用接口和工厂函数:

```go
type Product interface {
    Use() string
}

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

## 优点
1. **解耦** - 客户端与具体产品解耦
2. **易扩展** - 添加新产品不影响现有代码
3. **单一职责** - 创建逻辑集中管理

## 缺点
1. **类增多** - 每个产品需要一个工厂
2. **复杂度增加** - 引入额外的抽象层
