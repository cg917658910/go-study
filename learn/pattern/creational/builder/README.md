# Builder Pattern (生成器模式)

## 定义
生成器模式将复杂对象的构建与其表示分离,使得同样的构建过程可以创建不同的表示。

## 目的
- 分步骤构建复杂对象
- 隐藏对象构建细节
- 支持不同的表示

## 使用场景
- HTTP 请求构建
- SQL 查询构建
- 复杂配置对象创建
- 文档生成器(HTML、PDF、Markdown)

## Go 特有实现
使用方法链(Method Chaining):

```go
type Builder struct {
    product *Product
}

func (b *Builder) SetA(a string) *Builder {
    b.product.A = a
    return b
}

func (b *Builder) SetB(b int) *Builder {
    b.product.B = b
    return b
}

func (b *Builder) Build() *Product {
    return b.product
}
```

或使用函数选项模式:

```go
type Option func(*Product)

func WithA(a string) Option {
    return func(p *Product) {
        p.A = a
    }
}

func NewProduct(opts ...Option) *Product {
    p := &Product{}
    for _, opt := range opts {
        opt(p)
    }
    return p
}
```

## 优点
1. **分步骤构建** - 更清晰的构建过程
2. **代码复用** - 相同构建过程创建不同对象
3. **更好的控制** - 精细控制构建过程

## 缺点
1. **代码增多** - 需要创建额外的 Builder 类
2. **适用场景有限** - 简单对象不需要 Builder
