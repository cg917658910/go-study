package main

// Composite Pattern 组合模式 结构型设计模式
// 将对象组合成树形结构以表示“部分-整体”的层次结构 使得用户对单个对象和组合对象的使用具有一致性
// 适用于需要表示对象的部分-整体层次结构的场景

type MsgComponent interface {
	GetId() int
	SetId(id int)
	GetType() int
	SetType(msgType int)
}

type MsgModel struct {
	msgId   int
	msgType int
}

func (msg *MsgModel) GetId() int {
	return msg.msgId
}

func (msg *MsgModel) SetId(id int) {
	msg.msgId = id
}

func (msg *MsgModel) GetType() int {
	return msg.msgType
}

func (msg *MsgModel) SetType(msgType int) {
	msg.msgType = msgType
}

type GroupMsgModel struct {
	MsgModel
	children []MsgComponent
}

func (group *GroupMsgModel) Add(child MsgComponent) {
	group.children = append(group.children, child)
}
func (group *GroupMsgModel) Remove(child MsgComponent) {
	for i, c := range group.children {
		if c.GetId() == child.GetId() {
			group.children = append(group.children[:i], group.children[i+1:]...)
			return
		}
	}
}
func (group *GroupMsgModel) GetChildren() []MsgComponent {
	return group.children
}
func main() {
	// Example usage
	root := &GroupMsgModel{}
	root.SetId(1)
	root.SetType(0)

	child1 := &MsgModel{}
	child1.SetId(2)
	child1.SetType(1)

	child2 := &GroupMsgModel{}
	child2.SetId(3)
	child2.SetType(0)

	grandChild := &MsgModel{}
	grandChild.SetId(4)
	grandChild.SetType(1)

	child2.Add(grandChild)
	root.Add(child1)
	root.Add(child2)

	// Traverse the composite structure
	println("Root ID:", root.GetId())
	for _, child := range root.GetChildren() {
		println(" Child ID:", child.GetId())
		if groupChild, ok := child.(*GroupMsgModel); ok {
			for _, grandChild := range groupChild.GetChildren() {
				println("  GrandChild ID:", grandChild.GetId())
			}
		}
	}
}

// 优点: 简化客户端代码 统一对待单个对象和组合对象 易于扩展和维护
// 缺点: 设计复杂 可能导致过度设计
// 适用场景: 需要表示对象的部分-整体层次结构 需要统一对待单个对象和组合对象
// 对比装饰器模式: 装饰器模式关注在不改变对象接口的情况下动态添加行为 而组合模式关注表示对象的部分-整体层次结构
// 对比享元模式: 享元模式关注共享对象以节省内存 而组合模式关注表示对象的部分-整体层次结构
// 举例: 文件系统中文件和文件夹的关系 文件夹可以包含文件和子文件夹 通过组合模式可以统一处理文件和文件夹
// Web领域例子: HTML元素的嵌套关系 div 可以包含 p span 等元素 通过组合模式可以统一处理这些元素
// 数据库领域例子: SQL查询语句的嵌套关系 子查询可以作为主查询的一部分 通过组合模式可以统一处理查询语句
