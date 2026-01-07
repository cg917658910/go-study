package main

// Command Pattern 命令模式 行为型设计模式
// 将一个请求封装为一个对象 从而使你可用不同的请求对客户进行参数化
// 对请求排队或记录请求日志 以及支持可撤销的操作
// 适用于需要将请求调用者与请求接收者解耦的场景

type Command interface {
	Execute()
}

type Receiver struct {
}

func (r *Receiver) Action() {
	println("Receiver Action Executed")
}

type ConcreteCommand struct {
	receiver *Receiver
}

func NewConcreteCommand(receiver *Receiver) *ConcreteCommand {
	return &ConcreteCommand{receiver: receiver}
}

func (c *ConcreteCommand) Execute() {
	c.receiver.Action()
}

type Invoker struct {
	command Command
}

func (i *Invoker) SetCommand(command Command) {
	i.command = command
}

func (i *Invoker) Invoke() {
	i.command.Execute()
}

func main() {
	receiver := &Receiver{}
	command := NewConcreteCommand(receiver)
	invoker := &Invoker{}
	invoker.SetCommand(command)
	invoker.Invoke()
}

// 优点: 降低了请求发送者与接收者之间的耦合度 提高了系统的灵活性和可扩展性 支持请求的撤销和重做
// 缺点: 增加了系统的复杂度 可能导致类的数量增加
// 适用场景: 需要将请求调用者与请求接收者解耦 需要支持请求的排队和日志记录 需要支持可撤销的操作
// 对比职责链模式: 职责链模式关注请求的传递和处理 而命令模式关注请求的封装和执行
// 对比策略模式: 策略模式关注算法的选择 而命令模式关注请求的封装
// 举例: 实现一个远程控制器 可以通过命令模式封装不同的操作 如开灯 关灯 调节温度等
// Web领域例子: 实现一个任务调度系统 可以通过命令模式封装不同的任务 如发送邮件 生成报告等
// 数据库领域例子: 实现一个事务管理系统 可以通过命令模式封装不同的数据库操作 如插入数据 更新数据等
