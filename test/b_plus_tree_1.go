package test

import (
	"fmt"
	"sort"
)

// B+ 树阶数（每个节点最多存 order-1 个 key）
const order = 4

// Student 模拟数据库记录
type Student struct {
	ID   int
	Name string
	Age  int
}

// Node 结构体（B+ 树节点）
type Node struct {
	keys     []int
	children []*Node
	isLeaf   bool
	next     *Node
	records  map[int]*Student // 叶子节点存储数据
}

// BPlusTree 结构体
type BPlusTree struct {
	root *Node
}

// 创建一个新节点
func newNode(isLeaf bool) *Node {
	return &Node{
		keys:     []int{},
		children: []*Node{},
		isLeaf:   isLeaf,
		next:     nil,
		records:  make(map[int]*Student),
	}
}

// 插入数据
func (tree *BPlusTree) Insert(id int, student *Student) {
	if tree.root == nil {
		// 初始化根节点
		tree.root = newNode(true)
	}
	root := tree.root
	if len(root.keys) == order-1 {
		// 根节点满了，分裂
		newRoot := newNode(false)
		newRoot.children = append(newRoot.children, root)
		tree.splitChild(newRoot, 0)
		tree.root = newRoot
	}
	tree.insertNonFull(tree.root, id, student)
}

// 在未满的节点中插入 key
func (tree *BPlusTree) insertNonFull(node *Node, id int, student *Student) {
	if node.isLeaf {
		// 叶子节点：插入记录
		node.keys = append(node.keys, id)
		sort.Ints(node.keys)
		node.records[id] = student
	} else {
		// 内部节点：找到合适的子树插入
		i := len(node.keys) - 1
		for i >= 0 && id < node.keys[i] {
			i--
		}
		i++
		if len(node.children[i].keys) == order-1 {
			// 子节点满了，先分裂
			tree.splitChild(node, i)
			if id > node.keys[i] {
				i++
			}
		}
		tree.insertNonFull(node.children[i], id, student)
	}
}

// 分裂子节点
func (tree *BPlusTree) splitChild(parent *Node, index int) {
	child := parent.children[index]
	newChild := newNode(child.isLeaf)
	mid := order / 2

	// 分裂 keys
	parent.keys = append(parent.keys, child.keys[mid])
	sort.Ints(parent.keys)

	newChild.keys = append(newChild.keys, child.keys[mid+1:]...)
	child.keys = child.keys[:mid]

	// 分裂数据
	if child.isLeaf {
		for _, k := range newChild.keys {
			newChild.records[k] = child.records[k]
			delete(child.records, k)
		}
		// 叶子节点链表连接
		newChild.next = child.next
		child.next = newChild
	}

	// 插入新子节点
	parent.children = append(parent.children[:index+1], append([]*Node{newChild}, parent.children[index+1:]...)...)
}

// 查询 ID 对应的学生
func (tree *BPlusTree) Search(id int) *Student {
	node := tree.root
	for node != nil {
		i := 0
		for i < len(node.keys) && id > node.keys[i] {
			i++
		}
		if node.isLeaf {
			if i < len(node.keys) && id == node.keys[i] {
				return node.records[id]
			}
			return nil
		}
		node = node.children[i]
	}
	return nil
}

// 打印 B+ 树（层次遍历）
func (tree *BPlusTree) Print() {
	queue := []*Node{tree.root}
	for len(queue) > 0 {
		nextQueue := []*Node{}
		for _, node := range queue {
			fmt.Print(node.keys, "  ")
			if !node.isLeaf {
				nextQueue = append(nextQueue, node.children...)
			}
		}
		fmt.Println()
		queue = nextQueue
	}
}

// 测试
func testBPlusTree() {
	tree := &BPlusTree{}

	// 插入学生数据
	students := []Student{
		{ID: 101, Name: "Alice", Age: 22},
		{ID: 102, Name: "Bob", Age: 23},
		{ID: 103, Name: "Charlie", Age: 21},
		{ID: 104, Name: "David", Age: 24},
		{ID: 105, Name: "Eve", Age: 22},
	}

	for _, student := range students {
		tree.Insert(student.ID, &student)
	}

	tree.Print()

	// 查询学生
	fmt.Println("Search 103:", tree.Search(103)) // Charlie
	fmt.Println("Search 999:", tree.Search(999)) // nil
}
