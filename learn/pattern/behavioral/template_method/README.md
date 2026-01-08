# Template Method Pattern (模板方法模式)

## 定义
模板方法模式定义一个操作中的算法骨架,而将一些步骤延迟到子类中。

## 目的
- 定义算法骨架
- 延迟部分实现
- 代码复用

## 使用场景
- 框架和库的钩子方法
- 数据处理流程
- 测试框架
- 算法框架

## Go 特有实现
使用接口和嵌入:

```go
type Template interface {
    Step1()
    Step2()
}

type AbstractClass struct {
    impl Template
}

func (a *AbstractClass) TemplateMethod() {
    a.impl.Step1()
    a.impl.Step2()
}
```

## 优点
1. **代码复用** - 提取公共代码
2. **控制扩展** - 控制子类扩展点
3. **符合开闭原则** - 不修改模板扩展功能

## 缺点
1. **类增多** - 每个变体需要一个子类
2. **违反里氏替换** - 可能违反里氏替换原则
