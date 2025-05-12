package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// RadixNode 表示Radix树的节点
type RadixNode struct {
	prefix   string                // 节点存储的前缀
	children map[string]*RadixNode // 子节点
	isEnd    bool                  // 是否单词结束
	value    string                // 存储的值
}

// NewRadixNode 创建新的Radix节点
func NewRadixNode(prefix string) *RadixNode {
	return &RadixNode{
		prefix:   prefix,
		children: make(map[string]*RadixNode),
		isEnd:    false,
		value:    "",
	}
}

// RadixTree 表示Radix树结构
type RadixTree struct {
	root *RadixNode
}

// NewRadixTree 创建新的Radix树
func NewRadixTree() *RadixTree {
	return &RadixTree{root: NewRadixNode("")}
}

// Insert 插入单词和值
func (rt *RadixTree) Insert(key, value string) {
	node := rt.root
	for {
		// 查找与当前key有共同前缀的子节点
		_, child := rt.findChild(node, key)
		if child == nil {
			// 没有共同前缀的子节点，直接添加新节点
			newNode := NewRadixNode(key)
			newNode.isEnd = true
			newNode.value = value
			node.children[key] = newNode
			return
		}

		commonPrefix := rt.longestCommonPrefix(key, child.prefix)
		if commonPrefix == len(child.prefix) {
			// 完全匹配当前子节点的前缀
			key = key[commonPrefix:]
			node = child
			if len(key) == 0 {
				// 完全匹配，更新值
				child.isEnd = true
				child.value = value
				return
			}
		} else {
			// 部分匹配，需要分裂节点
			rt.splitNode(child, commonPrefix)
			key = key[commonPrefix:]
			node = child
			if len(key) == 0 {
				// 分裂后正好匹配
				child.isEnd = true
				child.value = value
				return
			}
			// 添加剩余部分的新节点
			newNode := NewRadixNode(key)
			newNode.isEnd = true
			newNode.value = value
			child.children[key] = newNode
			return
		}
	}
}

// findChild 查找与key有共同前缀的子节点
func (rt *RadixTree) findChild(node *RadixNode, key string) (string, *RadixNode) {
	for prefix, child := range node.children {
		if commonLen := rt.longestCommonPrefix(key, prefix); commonLen > 0 {
			return prefix, child
		}
	}
	return "", nil
}

// longestCommonPrefix 计算两个字符串的最长公共前缀长度
func (rt *RadixTree) longestCommonPrefix(a, b string) int {
	i := 0
	for ; i < len(a) && i < len(b) && a[i] == b[i]; i++ {
	}
	return i
}

// splitNode 分裂节点
func (rt *RadixTree) splitNode(node *RadixNode, splitPos int) {
	// 创建新节点保存分裂后的部分
	newChild := NewRadixNode(node.prefix[splitPos:])
	newChild.children = node.children
	newChild.isEnd = node.isEnd
	newChild.value = node.value

	// 更新当前节点
	node.prefix = node.prefix[:splitPos]
	node.children = make(map[string]*RadixNode)
	node.children[newChild.prefix] = newChild
	node.isEnd = false
	node.value = ""
}

// Search 精确查找单词
func (rt *RadixTree) Search(key string) (string, bool) {
	node := rt.root
	for len(key) > 0 {
		match, child := rt.findChild(node, key)
		if child == nil {
			return "", false
		}

		commonPrefix := rt.longestCommonPrefix(key, match)
		if commonPrefix != len(match) {
			return "", false
		}

		key = key[commonPrefix:]
		node = child
	}
	return node.value, node.isEnd
}

// StartsWith 查找前缀匹配的所有单词
func (rt *RadixTree) StartsWith(prefix string) []string {
	node := rt.root
	remaining := prefix
	for len(remaining) > 0 {
		match, child := rt.findChild(node, remaining)
		if child == nil {
			return nil
		}

		commonPrefix := rt.longestCommonPrefix(remaining, match)
		if commonPrefix < len(match) {
			if strings.HasPrefix(match, remaining) {
				// 当前节点前缀以剩余前缀开头，收集所有子节点
				var results []string
				rt.collectWords(child, remaining, &results)
				return results
			}
			return nil
		}

		remaining = remaining[commonPrefix:]
		node = child
	}

	// 完全匹配前缀，收集所有子节点
	var results []string
	rt.collectWords(node, prefix[:len(prefix)-len(remaining)], &results)
	return results
}

// collectWords 递归收集单词
func (rt *RadixTree) collectWords(node *RadixNode, prefix string, results *[]string) {
	if node.isEnd {
		*results = append(*results, prefix+node.prefix+" -> "+node.value)
	}

	for _, child := range node.children {
		rt.collectWords(child, prefix+node.prefix, results)
	}
}

// Delete 删除单词
func (rt *RadixTree) Delete(key string) bool {
	path := []*RadixNode{rt.root}
	pathKeys := []string{""}
	node := rt.root

	// 查找要删除的节点路径
	for len(key) > 0 {
		match, child := rt.findChild(node, key)
		if child == nil {
			return false
		}

		commonPrefix := rt.longestCommonPrefix(key, match)
		if commonPrefix != len(match) {
			return false
		}

		path = append(path, child)
		pathKeys = append(pathKeys, match)
		key = key[commonPrefix:]
		node = child
	}

	if !node.isEnd {
		return false
	}

	// 标记为非结束节点
	node.isEnd = false
	node.value = ""

	// 向上合并无用节点
	for i := len(path) - 1; i > 0; i-- {
		node = path[i]
		parent := path[i-1]

		if len(node.children) == 0 && !node.isEnd {
			delete(parent.children, pathKeys[i])
		} else if len(node.children) == 1 && !node.isEnd {
			// 合并单一子节点
			var child *RadixNode
			for _, v := range node.children {
				child = v
			}
			mergedPrefix := node.prefix + child.prefix
			mergedNode := NewRadixNode(mergedPrefix)
			mergedNode.children = child.children
			mergedNode.isEnd = child.isEnd
			mergedNode.value = child.value
			parent.children[mergedPrefix] = mergedNode
			delete(parent.children, node.prefix)
		} else {
			break
		}
	}

	return true
}

// BatchInsert 批量插入
func (rt *RadixTree) BatchInsert(items map[string]string) {
	for key, value := range items {
		rt.Insert(key, value)
	}
}

// BatchDelete 批量删除
func (rt *RadixTree) BatchDelete(keys []string) {
	for _, key := range keys {
		rt.Delete(key)
	}
}

// Print 打印Radix树结构
func (rt *RadixTree) Print() {
	fmt.Println("Radix树结构:")
	rt.printNode(rt.root, 0)
}

// printNode 递归打印节点
func (rt *RadixTree) printNode(node *RadixNode, level int) {
	prefix := strings.Repeat("  ", level)
	endMark := ""
	if node.isEnd {
		endMark = "*"
	}
	fmt.Printf("%s%s%s\n", prefix, node.prefix, endMark)
	for _, child := range node.children {
		rt.printNode(child, level+1)
	}
}

func main() {
	rt := NewRadixTree()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Radix树命令行工具")
	fmt.Println("命令:")
	fmt.Println("  insert <key> <value> - 插入键值")
	fmt.Println("  batch_insert <key1:value1,key2:value2,...> - 批量插入")
	fmt.Println("  search <key>        - 搜索键")
	fmt.Println("  prefix <prefix>     - 前缀搜索")
	fmt.Println("  delete <key>        - 删除键")
	fmt.Println("  batch_delete <key1,key2,...> - 批量删除")
	fmt.Println("  print               - 打印树结构")
	fmt.Println("  exit                - 退出")

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
				parts = append(parts, "v"+parts[1])
			}
			key := parts[1]
			value := strings.Join(parts[2:], " ")
			rt.Insert(key, value)
			fmt.Printf("已插入: %s -> %s\n", key, value)
			rt.Print()

		case "batch_insert":
			if len(parts) < 2 {
				fmt.Println("Usage: batch_insert <key1:value1,key2:value2,...>")
				continue
			}
			items := make(map[string]string)
			pairs := strings.Split(parts[1], ",")
			for _, pair := range pairs {
				kv := strings.Split(pair, ":")
				if len(kv) != 2 {
					fmt.Printf("无效键值对: %s\n", pair)
					continue
				}
				items[kv[0]] = kv[1]
			}
			rt.BatchInsert(items)
			fmt.Printf("已批量插入 %d 个键值对\n", len(items))
			rt.Print()

		case "search":
			if len(parts) < 2 {
				fmt.Println("Usage: search <key>")
				continue
			}
			key := parts[1]
			if value, found := rt.Search(key); found {
				fmt.Printf("找到: %s -> %s\n", key, value)
			} else {
				fmt.Printf("未找到: %s\n", key)
			}

		case "prefix":
			if len(parts) < 2 {
				fmt.Println("Usage: prefix <prefix>")
				continue
			}
			prefix := parts[1]
			results := rt.StartsWith(prefix)
			if len(results) > 0 {
				fmt.Printf("前缀匹配 %s:\n", prefix)
				for _, res := range results {
					fmt.Println("  ", res)
				}
			} else {
				fmt.Printf("没有找到以 %s 开头的键\n", prefix)
			}

		case "delete":
			if len(parts) < 2 {
				fmt.Println("Usage: delete <key>")
				continue
			}
			key := parts[1]
			if rt.Delete(key) {
				fmt.Printf("已删除: %s\n", key)
				rt.Print()
			} else {
				fmt.Printf("删除失败，未找到: %s\n", key)
			}

		case "batch_delete":
			if len(parts) < 2 {
				fmt.Println("Usage: batch_delete <key1,key2,...>")
				continue
			}
			keys := strings.Split(parts[1], ",")
			rt.BatchDelete(keys)
			fmt.Printf("已批量删除 %d 个键\n", len(keys))
			rt.Print()

		case "print":
			rt.Print()

		case "exit":
			fmt.Println("退出...")
			return

		default:
			fmt.Println("未知命令")
		}
	}
}
