package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// TrieNode 表示Trie树的节点
type TrieNode struct {
	children map[rune]*TrieNode // 子节点
	isEnd    bool               // 是否单词结束
	value    string             // 存储的值
}

// NewTrieNode 创建新的Trie节点
func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[rune]*TrieNode),
		isEnd:    false,
		value:    "",
	}
}

// Trie 表示Trie树结构
type Trie struct {
	root *TrieNode
}

// NewTrie 创建新的Trie树
func NewTrie() *Trie {
	return &Trie{root: NewTrieNode()}
}

// Insert 插入单词和值
func (t *Trie) Insert(word string, value string) {
	node := t.root
	for _, ch := range word {
		if _, ok := node.children[ch]; !ok {
			node.children[ch] = NewTrieNode()
		}
		node = node.children[ch]
	}
	node.isEnd = true
	node.value = value
}

// Search 精确查找单词
func (t *Trie) Search(word string) (string, bool) {
	node := t.root
	for _, ch := range word {
		if _, ok := node.children[ch]; !ok {
			return "", false
		}
		node = node.children[ch]
	}
	return node.value, node.isEnd
}

// StartsWith 查找前缀匹配的所有单词
func (t *Trie) StartsWith(prefix string) []string {
	node := t.root
	// 先定位到前缀的最后一个节点
	for _, ch := range prefix {
		if _, ok := node.children[ch]; !ok {
			return nil
		}
		node = node.children[ch]
	}

	// 收集所有以该前缀开头的单词
	var results []string
	t.collectWords(node, prefix, &results)
	return results
}

// collectWords 递归收集单词
func (t *Trie) collectWords(node *TrieNode, prefix string, results *[]string) {
	if node.isEnd {
		*results = append(*results, prefix+" -> "+node.value)
	}

	for ch, child := range node.children {
		t.collectWords(child, prefix+string(ch), results)
	}
}

// Delete 删除单词
func (t *Trie) Delete(word string) bool {
	nodes := make([]*TrieNode, 0, len(word)+1)
	node := t.root
	nodes = append(nodes, node)

	// 查找单词路径上的所有节点
	for _, ch := range word {
		if _, ok := node.children[ch]; !ok {
			return false // 单词不存在
		}
		node = node.children[ch]
		nodes = append(nodes, node)
	}

	if !node.isEnd {
		return false // 不是完整单词
	}

	// 标记为非单词结尾
	node.isEnd = false
	node.value = ""

	// 从后向前删除无用节点
	for i := len(nodes) - 1; i > 0; i-- {
		current := nodes[i]
		parent := nodes[i-1]
		ch := rune(word[i-1])

		if len(current.children) == 0 && !current.isEnd {
			delete(parent.children, ch)
		} else {
			break
		}
	}

	return true
}

// Print 打印Trie树结构
func (t *Trie) Print() {
	fmt.Println("Trie结构:")
	t.printNode(t.root, 0)
}

// printNode 递归打印节点
func (t *Trie) printNode(node *TrieNode, level int) {
	prefix := strings.Repeat("  ", level)
	for ch, child := range node.children {
		endMark := ""
		if child.isEnd {
			endMark = "*"
		}
		fmt.Printf("%s├─ %c%s\n", prefix, ch, endMark)
		t.printNode(child, level+1)
	}
}

func main() {
	trie := NewTrie()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Trie树命令行工具")
	fmt.Println("命令:")
	fmt.Println("  insert <word> <value> - 插入单词")
	fmt.Println("  search <word>        - 搜索单词")
	fmt.Println("  prefix <prefix>      - 前缀搜索")
	fmt.Println("  delete <word>        - 删除单词")
	fmt.Println("  print                - 打印Trie结构")
	fmt.Println("  exit                 - 退出")

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
				/* fmt.Println("Usage: insert <word> <value>")
				continue */
				parts = append(parts, "v"+parts[1])
			}
			word := parts[1]
			value := strings.Join(parts[2:], " ")
			trie.Insert(word, value)
			fmt.Printf("已插入: %s -> %s\n", word, value)
			trie.Print()

		case "search":
			if len(parts) < 2 {
				fmt.Println("Usage: search <word>")
				continue
			}
			word := parts[1]
			if value, found := trie.Search(word); found {
				fmt.Printf("找到: %s -> %s\n", word, value)
			} else {
				fmt.Printf("未找到: %s\n", word)
			}

		case "prefix":
			if len(parts) < 2 {
				fmt.Println("Usage: prefix <prefix>")
				continue
			}
			prefix := parts[1]
			results := trie.StartsWith(prefix)
			if len(results) > 0 {
				fmt.Printf("前缀匹配 %s:\n", prefix)
				for _, res := range results {
					fmt.Println("  ", res)
				}
			} else {
				fmt.Printf("没有找到以 %s 开头的单词\n", prefix)
			}

		case "delete":
			if len(parts) < 2 {
				fmt.Println("Usage: delete <word>")
				continue
			}
			word := parts[1]
			if trie.Delete(word) {
				fmt.Printf("已删除: %s\n", word)
			} else {
				fmt.Printf("删除失败，未找到: %s\n", word)
			}
			trie.Print()
		case "print":
			trie.Print()

		case "exit":
			fmt.Println("退出...")
			return

		default:
			fmt.Println("未知命令")
		}
	}
}
