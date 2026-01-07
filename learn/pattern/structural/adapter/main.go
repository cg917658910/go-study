package main

// Adapter Pattern 适配器模式 结构型设计模式
// 将一个类的接口转换成客户希望的另一个接口 使得原本由于接口不兼容而不能一起工作的那些类可以一起工作
// 适用于系统需要使用现有类 而这些类的接口不符合系统的需求的场景

type Target interface {
	Request()
}

type Adaptee struct {
}

func (a *Adaptee) SpecificRequest() {
	println("Specific request from Adaptee")
}

type Adapter struct {
	adaptee *Adaptee
}

func NewAdapter(adaptee *Adaptee) *Adapter {
	return &Adapter{adaptee: adaptee}
}

func (adapter *Adapter) Request() {
	adapter.adaptee.SpecificRequest()
}

func main() {
	var target Target = NewAdapter(&Adaptee{})
	target.Request()
}

// 优点: 复用现有类 遵循开闭原则 提高代码的灵活性
// 缺点: 增加系统复杂度 可能影响性能
// 适用场景: 需要使用现有类 而这些类的接口不符合系统的需求 需要将多个不兼容的接口统一到一个接口中
// 对比装饰器模式: 装饰器模式关注在不改变对象接口的情况下 动态地给对象添加功能 而适配器模式关注将一个接口转换成另一个接口
// 对比桥接模式: 桥接模式关注将抽象部分与实现部分分离 使它们可以独立变化 而适配器模式关注将一个接口转换成另一个接口
// 举例: 将一个旧系统的接口适配到新系统中 使得新系统可以使用旧系统的功能
// Web领域例子: 适配不同的支付网关接口 使得系统可以使用多种支付方式
// 数据库领域例子: 适配不同的数据库驱动接口 使得系统可以支持多种数据库
