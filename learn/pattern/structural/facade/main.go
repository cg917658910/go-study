package main

// Facade Pattern 外观模式 结构型设计模式
// 为子系统中的一组接口提供一个一致的界面 外观模式定义了一个高层接口 使得子系统更易使用
// 适用于需要为复杂子系统提供一个简单接口的场景

type SubsystemA struct {
}

func (s *SubsystemA) OperationA() {
	println("Subsystem A Operation")
}

type SubsystemB struct {
}

func (s *SubsystemB) OperationB() {
	println("Subsystem B Operation")
}

type Facade struct {
	subsystemA *SubsystemA
	subsystemB *SubsystemB
}

func NewFacade() *Facade {
	return &Facade{
		subsystemA: &SubsystemA{},
		subsystemB: &SubsystemB{},
	}
}

func (f *Facade) Operation() {
	f.subsystemA.OperationA()
	f.subsystemB.OperationB()
}

func main() {
	facade := NewFacade()
	facade.Operation()
}

// 优点: 简化了客户端与子系统的交互 隐藏了子系统的复杂性 提高了代码的可维护性
// 缺点: 增加了系统的复杂度 可能导致子系统变得臃肿
// 适用场景: 需要为复杂子系统提供一个简单接口 需要减少客户端与子系统之间的依赖
// 对比代理模式: 代理模式关注控制对对象的访问 而外观模式关注简化接口
// 对比装饰器模式: 装饰器模式关注在不改变对象接口的情况下 动态地给对象添加功能 而外观模式关注提供一个简单的接口
// 举例: 为一个复杂的图形绘制库提供一个简单的绘图接口 使得客户端可以方便地绘制图形而不需要了解底层实现细节
// Web领域例子: 为一个复杂的Web服务提供一个简单的API接口 使得客户端可以方便地调用服务而不需要了解底层实现细节
// 数据库领域例子: 为一个复杂的数据库操作库提供一个简单的查询接口 使得客户端可以方便地执行查询操作而不需要了解底层实现细节
