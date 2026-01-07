package main

// Strategy Pattern 策略模式 行为型设计模式
// 定义一系列算法 将每一个算法封装起来 并使它们可以相互替换
// 使得算法可以独立于使用它的客户而变化

type Strategy interface {
	Execute(a, b int) int
}

type AddStrategy struct {
}

func (s *AddStrategy) Execute(a, b int) int {
	return a + b
}

type SubtractStrategy struct {
}

func (s *SubtractStrategy) Execute(a, b int) int {
	return a - b
}

type Context struct {
	strategy Strategy
}

func (c *Context) SetStrategy(strategy Strategy) {
	c.strategy = strategy
}

func (c *Context) ExecuteStrategy(a, b int) int {
	return c.strategy.Execute(a, b)
}

func main() {
	context := &Context{}

	context.SetStrategy(&AddStrategy{})
	result1 := context.ExecuteStrategy(5, 3)
	println("Addition Result:", result1)

	context.SetStrategy(&SubtractStrategy{})
	result2 := context.ExecuteStrategy(5, 3)
	println("Subtraction Result:", result2)
}

// 优点: 遵循开闭原则 可以在不修改现有代码的情况下添加新算法 提高代码的灵活性
// 缺点: 增加了类的数量和复杂度 客户端必须了解不同的策略
// 适用场景: 需要在运行时选择算法 需要避免使用大量的条件语句来选择算法
// 对比状态模式: 状态模式关注对象的状态变化 而策略模式关注算法的选择
// 对比简单工厂模式: 简单工厂模式集中管理对象创建 但违反开闭原则 不易扩展 而策略模式通过封装算法 遵循开闭原则
// 举例: 不同的排序算法 如快速排序 归并排序等 使用策略模式可以动态地选择排序算法
// Web领域例子: 不同的压缩算法 如Gzip Brotli等 使用策略模式可以动态地选择压缩算法
// 数据库领域例子: 不同的查询优化策略 使用策略模式可以动态地选择查询优化策略
