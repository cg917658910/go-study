package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	DEFAULT_ORDER = 4    // 默认阶数
	PAGE_SIZE     = 4096 // 页大小(字节)
	KEY_SIZE      = 8    // 键大小(字节)
	POINTER_SIZE  = 8    // 指针大小(字节)
)

// BPlusTree MySQL风格的B+树
type BPlusTree struct {
	root      *Page
	order     int
	pageCount int
}

// Page 表示B+树的页(节点)
type Page struct {
	id        int
	isLeaf    bool
	keys      []int
	pointers  []interface{} // 对于叶子节点存储值，内部节点存储子页指针
	next      *Page         // 叶子节点的链表指针
	parent    *Page
	freeSpace int
}

// Record 表示数据库记录
type Record struct {
	key    int
	value  string
	status byte // 用于事务: 0=正常, 1=已删除
	txID   int  // 事务ID
}

// NewBPlusTree 创建新的B+树
func NewBPlusTree(order int) *BPlusTree {
	if order < 3 {
		order = DEFAULT_ORDER
	}
	return &BPlusTree{
		order: order,
		root:  createNewPage(true),
	}
}

// createNewPage 创建新页
func createNewPage(isLeaf bool) *Page {
	return &Page{
		id:        generatePageID(),
		isLeaf:    isLeaf,
		keys:      make([]int, 0),
		pointers:  make([]interface{}, 0),
		freeSpace: PAGE_SIZE - 16, // 减去页头开销
	}
}

var pageIDCounter = 0

func generatePageID() int {
	pageIDCounter++
	return pageIDCounter
}

// Insert 插入键值对
func (t *BPlusTree) Insert(key int, value string) {
	// 查找合适的叶子页
	leaf := t.findLeafPage(key)

	// 检查键是否已存在
	for i, k := range leaf.keys {
		if k == key {
			// 更新现有记录
			if rec, ok := leaf.pointers[i].(*Record); ok {
				rec.value = value
			}
			return
		}
	}

	// 插入新记录
	record := &Record{key: key, value: value}
	t.insertIntoLeaf(leaf, key, record)
}

// findLeafPage 查找包含key的叶子页
func (t *BPlusTree) findLeafPage(key int) *Page {
	current := t.root
	for !current.isLeaf {
		// 二分查找合适的子页
		idx := 0
		for idx < len(current.keys) && key >= current.keys[idx] {
			idx++
		}
		if child, ok := current.pointers[idx].(*Page); ok {
			current = child
		}
	}
	return current
}

// insertIntoLeaf 向叶子页插入记录
func (t *BPlusTree) insertIntoLeaf(leaf *Page, key int, record *Record) {
	// 找到插入位置
	idx := 0
	for idx < len(leaf.keys) && key > leaf.keys[idx] {
		idx++
	}

	// 插入键和记录
	leaf.keys = append(leaf.keys[:idx], append([]int{key}, leaf.keys[idx:]...)...)
	leaf.pointers = append(leaf.pointers[:idx], append([]interface{}{record}, leaf.pointers[idx:]...)...)

	// 更新空闲空间
	leaf.freeSpace -= (KEY_SIZE + POINTER_SIZE)

	// 检查是否需要分裂
	if len(leaf.keys) > t.order {
		t.splitLeafPage(leaf)
	}
}

// splitLeafPage 分裂叶子页
func (t *BPlusTree) splitLeafPage(leaf *Page) {
	// 创建新页
	newPage := createNewPage(true)

	// 计算分裂点
	splitIdx := len(leaf.keys) / 2

	// 移动一半键和记录到新页
	newPage.keys = append(newPage.keys, leaf.keys[splitIdx:]...)
	newPage.pointers = append(newPage.pointers, leaf.pointers[splitIdx:]...)

	// 更新原页
	leaf.keys = leaf.keys[:splitIdx]
	leaf.pointers = leaf.pointers[:splitIdx]

	// 更新链表指针
	newPage.next = leaf.next
	leaf.next = newPage
	newPage.parent = leaf.parent

	// 更新空闲空间
	leaf.freeSpace = PAGE_SIZE - 16 - (len(leaf.keys)*KEY_SIZE + len(leaf.pointers)*POINTER_SIZE)
	newPage.freeSpace = PAGE_SIZE - 16 - (len(newPage.keys)*KEY_SIZE + len(newPage.pointers)*POINTER_SIZE)

	// 插入分隔键到父页
	t.insertIntoParent(leaf, newPage.keys[0], newPage)
}

// insertIntoParent 更新父页
func (t *BPlusTree) insertIntoParent(left *Page, key int, right *Page) {
	parent := left.parent

	// 如果是根页分裂
	if parent == nil {
		newRoot := createNewPage(false)
		newRoot.keys = append(newRoot.keys, key)
		newRoot.pointers = append(newRoot.pointers, left, right)
		left.parent = newRoot
		right.parent = newRoot
		t.root = newRoot
		return
	}

	// 找到插入位置
	idx := 0
	for idx < len(parent.keys) && key > parent.keys[idx] {
		idx++
	}

	// 插入键和指针
	parent.keys = append(parent.keys[:idx], append([]int{key}, parent.keys[idx:]...)...)
	parent.pointers = append(parent.pointers[:idx+1], append([]interface{}{right}, parent.pointers[idx+1:]...)...)

	// 更新空闲空间
	parent.freeSpace -= (KEY_SIZE + POINTER_SIZE)

	// 检查是否需要分裂
	if len(parent.keys) > t.order {
		t.splitInternalPage(parent)
	}
}

// splitInternalPage 分裂内部页
func (t *BPlusTree) splitInternalPage(page *Page) {
	// 创建新页
	newPage := createNewPage(false)

	// 计算分裂点
	splitIdx := len(page.keys) / 2
	promoteKey := page.keys[splitIdx]

	// 移动一半键和指针到新页
	newPage.keys = append(newPage.keys, page.keys[splitIdx+1:]...)
	newPage.pointers = append(newPage.pointers, page.pointers[splitIdx+1:]...)

	// 更新原页
	page.keys = page.keys[:splitIdx]
	page.pointers = page.pointers[:splitIdx+1]

	// 更新子页的父指针
	for _, p := range newPage.pointers {
		if child, ok := p.(*Page); ok {
			child.parent = newPage
		}
	}

	newPage.parent = page.parent

	// 更新空闲空间
	page.freeSpace = PAGE_SIZE - 16 - (len(page.keys)*KEY_SIZE + len(page.pointers)*POINTER_SIZE)
	newPage.freeSpace = PAGE_SIZE - 16 - (len(newPage.keys)*KEY_SIZE + len(newPage.pointers)*POINTER_SIZE)

	// 插入分隔键到父页
	t.insertIntoParent(page, promoteKey, newPage)
}

// Search 搜索记录
func (t *BPlusTree) Search(key int) *Record {
	leaf := t.findLeafPage(key)
	for i, k := range leaf.keys {
		if k == key {
			if rec, ok := leaf.pointers[i].(*Record); ok {
				return rec
			}
		}
	}
	return nil
}

// Print 打印B+树结构
func (t *BPlusTree) Print() {
	if t.root == nil {
		fmt.Println("B+树为空")
		return
	}

	queue := []*Page{t.root}
	level := 0

	for len(queue) > 0 {
		fmt.Printf("Level %d:\n", level)
		nextQueue := []*Page{}
		for _, page := range queue {
			fmt.Printf("  Page %d: %v ", page.id, page.keys)
			if page.isLeaf {
				fmt.Print("(leaf)")
				if page.next != nil {
					fmt.Printf(" -> Page %d", page.next.id)
				}
			} else {
				fmt.Print("(internal)")
				for _, p := range page.pointers {
					if child, ok := p.(*Page); ok {
						nextQueue = append(nextQueue, child)
					}
				}
			}
			fmt.Println()
		}
		queue = nextQueue
		level++
	}
}

// 插入相关方法保持不变...

// Delete 删除键
func (t *BPlusTree) Delete(key int) bool {
	leaf := t.findLeafPage(key)
	idx := -1
	for i, k := range leaf.keys {
		if k == key {
			idx = i
			break
		}
	}
	if idx == -1 {
		return false
	}

	// 标记为已删除（逻辑删除）
	if rec, ok := leaf.pointers[idx].(*Record); ok {
		rec.status = 1
	}

	// 物理删除（可选）
	// leaf.keys = append(leaf.keys[:idx], leaf.keys[idx+1:]...)
	// leaf.pointers = append(leaf.pointers[:idx], leaf.pointers[idx+1:]...)
	// leaf.freeSpace += KEY_SIZE + POINTER_SIZE
	// t.rebalance(leaf)

	return true
}

// rebalance 重新平衡页
func (t *BPlusTree) rebalance(p *Page) {
	if len(p.keys) < (t.order+1)/2 && p.parent != nil {
		parent := p.parent
		idx := t.getChildIndex(parent, p)

		// 尝试从左兄弟借
		if idx > 0 {
			leftSibling := parent.pointers[idx-1].(*Page)
			if len(leftSibling.keys) > (t.order+1)/2 {
				t.redistribute(leftSibling, p, parent, idx-1, false)
				return
			}
		}

		// 尝试从右兄弟借
		if idx < len(parent.pointers)-1 {
			rightSibling := parent.pointers[idx+1].(*Page)
			if len(rightSibling.keys) > (t.order+1)/2 {
				t.redistribute(p, rightSibling, parent, idx, true)
				return
			}
		}

		// 需要合并
		if idx > 0 {
			// 与左兄弟合并
			leftSibling := parent.pointers[idx-1].(*Page)
			t.merge(leftSibling, p, parent, idx-1)
		} else {
			// 与右兄弟合并
			rightSibling := parent.pointers[idx+1].(*Page)
			t.merge(p, rightSibling, parent, idx)
		}
	}
}

// redistribute 重新分配键
func (t *BPlusTree) redistribute(left, right *Page, parent *Page, parentIdx int, borrowFromRight bool) {
	if borrowFromRight {
		// 从右兄弟借一个键
		if left.isLeaf {
			// 叶子节点
			left.keys = append(left.keys, right.keys[0])
			left.pointers = append(left.pointers, right.pointers[0])
			right.keys = right.keys[1:]
			right.pointers = right.pointers[1:]
			parent.keys[parentIdx] = right.keys[0]
		} else {
			// 内部节点
			left.keys = append(left.keys, parent.keys[parentIdx])
			parent.keys[parentIdx] = right.keys[0]
			left.pointers = append(left.pointers, right.pointers[0])
			right.keys = right.keys[1:]
			right.pointers = right.pointers[1:]
		}
	} else {
		// 从左兄弟借一个键
		if left.isLeaf {
			// 叶子节点
			right.keys = append([]int{left.keys[len(left.keys)-1]}, right.keys...)
			right.pointers = append([]interface{}{left.pointers[len(left.pointers)-1]}, right.pointers...)
			left.keys = left.keys[:len(left.keys)-1]
			left.pointers = left.pointers[:len(left.pointers)-1]
			parent.keys[parentIdx] = right.keys[0]
		} else {
			// 内部节点
			right.keys = append([]int{parent.keys[parentIdx]}, right.keys...)
			parent.keys[parentIdx] = left.keys[len(left.keys)-1]
			right.pointers = append([]interface{}{left.pointers[len(left.pointers)-1]}, right.pointers...)
			left.keys = left.keys[:len(left.keys)-1]
			left.pointers = left.pointers[:len(left.pointers)-1]
		}
	}
}

// merge 合并页
func (t *BPlusTree) merge(left, right *Page, parent *Page, parentIdx int) {
	if left.isLeaf {
		// 合并叶子节点
		left.keys = append(left.keys, right.keys...)
		left.pointers = append(left.pointers, right.pointers...)
		left.next = right.next
	} else {
		// 合并内部节点
		left.keys = append(left.keys, parent.keys[parentIdx])
		left.keys = append(left.keys, right.keys...)
		left.pointers = append(left.pointers, right.pointers...)
	}

	// 更新父节点
	parent.keys = append(parent.keys[:parentIdx], parent.keys[parentIdx+1:]...)
	parent.pointers = append(parent.pointers[:parentIdx+1], parent.pointers[parentIdx+2:]...)

	// 如果父节点是根且变空
	if parent == t.root && len(parent.keys) == 0 {
		t.root = left
		left.parent = nil
	}
}

// getChildIndex 获取子页索引
func (t *BPlusTree) getChildIndex(parent, child *Page) int {
	for i, p := range parent.pointers {
		if p == child {
			return i
		}
	}
	return -1
}

// BatchInsert 批量插入
func (t *BPlusTree) BatchInsert(records map[int]string) {
	for key, value := range records {
		t.Insert(key, value)
	}
}

// BatchDelete 批量删除
func (t *BPlusTree) BatchDelete(keys []int) {
	for _, key := range keys {
		t.Delete(key)
	}
}

// 其他方法保持不变...

func main() {
	tree := NewBPlusTree(3)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("MySQL风格B+树实现（完整版）")
	fmt.Println("命令:")
	fmt.Println("  insert <key> <value> - 插入记录")
	fmt.Println("  delete <key> - 删除记录")
	fmt.Println("  batch_insert <key1:value1,key2:value2,...> - 批量插入")
	fmt.Println("  batch_delete <key1,key2,...> - 批量删除")
	fmt.Println("  search <key> - 搜索记录")
	fmt.Println("  print - 打印树结构")
	fmt.Println("  exit - 退出")

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		parts := strings.Split(input, " ")

		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "insert":
			if len(parts) < 3 {
				parts = append(parts, fmt.Sprint("v", parts[1]))
				/* fmt.Println("Usage: insert <key> <value>")
				continue */
			}
			key, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("无效键")
				continue
			}
			value := strings.Join(parts[2:], " ")
			tree.Insert(key, value)
			fmt.Printf("已插入记录: %d -> %s\n", key, value)
			tree.Print()

		case "delete":
			if len(parts) < 2 {
				fmt.Println("Usage: delete <key>")
				continue
			}
			key, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("无效键")
				continue
			}
			if tree.Delete(key) {
				fmt.Printf("已删除键 %d\n", key)
			} else {
				fmt.Printf("键 %d 未找到\n", key)
			}

		case "batch_insert":
			if len(parts) < 2 {
				fmt.Println("Usage: batch_insert <key1:value1,key2:value2,...>")
				continue
			}
			records := make(map[int]string)
			pairs := strings.Split(parts[1], ",")
			for _, pair := range pairs {
				kv := strings.Split(pair, ":")
				if len(kv) != 2 {
					fmt.Printf("无效键值对: %s\n", pair)
					continue
				}
				key, err := strconv.Atoi(kv[0])
				if err != nil {
					fmt.Printf("无效键: %s\n", kv[0])
					continue
				}
				records[key] = kv[1]
			}
			tree.BatchInsert(records)
			fmt.Printf("已批量插入 %d 条记录\n", len(records))
			tree.Print()

		case "batch_delete":
			if len(parts) < 2 {
				fmt.Println("Usage: batch_delete <key1,key2,...>")
				continue
			}
			keys := make([]int, 0)
			for _, k := range strings.Split(parts[1], ",") {
				key, err := strconv.Atoi(k)
				if err != nil {
					fmt.Printf("无效键: %s\n", k)
					continue
				}
				keys = append(keys, key)
			}
			tree.BatchDelete(keys)
			fmt.Printf("已批量删除 %d 个键\n", len(keys))

		case "search":
			if len(parts) < 2 {
				fmt.Println("Usage: search <key>")
				continue
			}
			key, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("无效键")
				continue
			}
			record := tree.Search(key)
			if record != nil {
				if record.status == 0 {
					fmt.Printf("找到记录: %d -> %s\n", record.key, record.value)
				} else {
					fmt.Printf("键 %d 已被删除\n", key)
				}
			} else {
				fmt.Printf("键 %d 未找到\n", key)
			}

		case "print":
			tree.Print()

		case "exit":
			fmt.Println("退出中...")
			return

		default:
			fmt.Println("未知命令")
		}
	}
}
