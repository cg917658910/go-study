package main

// Template Method Pattern 模板方法模式 行为型设计模式
// 定义一个操作中的算法骨架 将一些步骤延迟到子类中
// 使得子类可以在不改变算法结构的情况下重新定义算法的某些特定步骤
// 适用于多个子类有相似行为 但某些细节不同的场景

type AbstractClass interface {
	TemplateMethod()
	PrimitiveOperation1()
	PrimitiveOperation2()
}

type BaseClass struct {
}

func (b *BaseClass) TemplateMethod() {
	b.PrimitiveOperation1()
	b.PrimitiveOperation2()
}

func (b *BaseClass) PrimitiveOperation1() {
	// 默认实现 可以被子类覆盖
	println("BaseClass PrimitiveOperation1")
}

func (b *BaseClass) PrimitiveOperation2() {
	// 默认实现 可以被子类覆盖
	println("BaseClass PrimitiveOperation2")
}

type ConcreteClassA struct {
	BaseClass
}

func (c *ConcreteClassA) PrimitiveOperation1() {
	println("ConcreteClassA PrimitiveOperation1")
}

type ConcreteClassB struct {
	BaseClass
}

func (c *ConcreteClassB) PrimitiveOperation2() {
	println("ConcreteClassB PrimitiveOperation2")
}

func main() {
	var classA AbstractClass = &ConcreteClassA{}
	classA.TemplateMethod()

	var classB AbstractClass = &ConcreteClassB{}
	classB.TemplateMethod()
}

// 优点: 遵循开闭原则 可以在不修改现有代码的情况下添加新行为 提高代码的复用性
// 缺点: 增加了类的数量和复杂度 子类必须实现所有的抽象方法
// 适用场景: 多个子类有相似行为 但某些细节不同 需要通过子类来实现这些细节
// 对比策略模式: 策略模式关注算法的选择 而模板方法模式关注算法的结构
// 对比工厂方法模式: 工厂方法模式关注对象的创建 而模板方法模式关注算法的实现
// 举例: 制作饮料的过程 包括煮沸水 泡茶或咖啡 倒入杯中等步骤 使用模板方法模式可以定义一个制作饮料的模板方法 具体的泡茶或咖啡步骤由子类实现
// Web领域例子: 处理HTTP请求的过程 包括解析请求 验证身份 处理业务逻辑 返回响应等步骤 使用模板方法模式可以定义一个处理请求的模板方法 具体的业务逻辑由子类实现
// 数据库领域例子: 执行数据库操作的过程 包括建立连接 开始事务 执行操作 提交或回滚事务等步骤 使用模板方法模式可以定义一个执行数据库操作的模板方法 具体的操作由子类实现
