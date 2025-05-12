package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	MaxLevel    = 6   // 跳表最大层数
	Probability = 0.5 // 每升一层的概率
)

type Node struct {
	key   int
	value string
	next  []*Node // 每一层的 next 指针
}

type SkipList struct {
	head  *Node
	level int
}

func NewNode(level int, key int, value string) *Node {
	return &Node{
		key:   key,
		value: value,
		next:  make([]*Node, level),
	}
}

func NewSkipList() *SkipList {
	rand.NewSource(time.Now().UnixNano())
	return &SkipList{
		head:  NewNode(MaxLevel, -1, ""), // head 为哨兵节点
		level: 1,
	}
}

// 随机生成节点的层数
func randomLevel() int {
	level := 1
	for rand.Float64() < Probability && level < MaxLevel {
		level++
	}
	return level
}

// 查找 key
func (sl *SkipList) Search(key int) (string, bool) {
	curr := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for curr.next[i] != nil && curr.next[i].key < key {
			curr = curr.next[i]
		}
	}
	curr = curr.next[0]
	if curr != nil && curr.key == key {
		return curr.value, true
	}
	return "", false
}

// 插入 key-value
func (sl *SkipList) Insert(key int, value string) {
	update := make([]*Node, MaxLevel)
	curr := sl.head

	// 记录每一层应该更新的前驱节点
	for i := sl.level - 1; i >= 0; i-- {
		for curr.next[i] != nil && curr.next[i].key < key {
			curr = curr.next[i]
		}
		update[i] = curr
	}

	// 如果 key 已存在，则更新
	if curr.next[0] != nil && curr.next[0].key == key {
		curr.next[0].value = value
		return
	}

	newLevel := randomLevel()
	if newLevel > sl.level {
		for i := sl.level; i < newLevel; i++ {
			update[i] = sl.head
		}
		sl.level = newLevel
	}

	newNode := NewNode(newLevel, key, value)
	for i := 0; i < newLevel; i++ {
		newNode.next[i] = update[i].next[i]
		update[i].next[i] = newNode
	}
}

// 删除节点
func (sl *SkipList) Delete(key int) bool {
	update := make([]*Node, MaxLevel)
	curr := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for curr.next[i] != nil && curr.next[i].key < key {
			curr = curr.next[i]
		}
		update[i] = curr
	}

	target := curr.next[0]
	if target == nil || target.key != key {
		return false
	}

	for i := 0; i < sl.level; i++ {
		if update[i].next[i] != target {
			break
		}
		update[i].next[i] = target.next[i]
	}

	// 如果最高层没有节点了，降低 level
	for sl.level > 1 && sl.head.next[sl.level-1] == nil {
		sl.level--
	}

	return true
}

// 打印结构（从高层到底层）
func (sl *SkipList) Print() {
	fmt.Println("\n--- SkipList Structure ---")
	for i := sl.level - 1; i >= 0; i-- {
		curr := sl.head
		fmt.Printf("Level %d: ", i)
		for curr.next[i] != nil {
			fmt.Printf("%d:%s → ", curr.next[i].key, curr.next[i].value)
			curr = curr.next[i]
		}
		fmt.Println("NULL")
	}
	fmt.Println()
}

func main() {
	skipList := NewSkipList()
	skipList.Insert(1, "one")
	skipList.Insert(2, "two")
	skipList.Insert(3, "three")
	skipList.Insert(4, "four")
	skipList.Insert(5, "five")

	skipList.Print()

	value, found := skipList.Search(3)
	if found {
		fmt.Printf("Found: 3 -> %s\n", value)
	} else {
		fmt.Println("Not Found: 3")
	}

	skipList.Delete(3)
	skipList.Print()
}
