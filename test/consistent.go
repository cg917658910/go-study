package test

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"sort"
	"strconv"
	"sync"
)

// ConsistentHash 实现一致性哈希
type ConsistentHash struct {
	replicas   int               // 每个节点的虚拟节点数
	hashFunc   hash.Hash         // 哈希函数
	sortedKeys []uint32          // 哈希环（有序）
	hashMap    map[uint32]string // 虚拟节点 -> 真实节点映射
	nodes      map[string]bool   // 真实节点集合
	mu         sync.RWMutex      // 读写锁
}

// NewConsistentHash 创建一个一致性哈希
func NewConsistentHash(replicas int) *ConsistentHash {
	return &ConsistentHash{
		replicas:   replicas,
		hashFunc:   sha256.New(),
		hashMap:    make(map[uint32]string),
		nodes:      make(map[string]bool),
		sortedKeys: []uint32{},
	}
}

// hashKey 计算哈希值（返回 32 位整数）
func (c *ConsistentHash) hashKey(key string) uint32 {
	c.hashFunc.Reset()
	c.hashFunc.Write([]byte(key))
	hashBytes := c.hashFunc.Sum(nil)
	return uint32(hashBytes[0])<<24 | uint32(hashBytes[1])<<16 | uint32(hashBytes[2])<<8 | uint32(hashBytes[3])
}

// AddNode 添加一个真实节点（带虚拟节点）
func (c *ConsistentHash) AddNode(node string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.nodes[node] {
		return // 避免重复添加
	}

	// 真实节点加入
	c.nodes[node] = true

	// 生成虚拟节点
	for i := 0; i < c.replicas; i++ {
		vNodeKey := node + "#" + strconv.Itoa(i)
		hashValue := c.hashKey(vNodeKey)
		c.hashMap[hashValue] = node
		c.sortedKeys = append(c.sortedKeys, hashValue)
	}

	// 重新排序哈希环
	sort.Slice(c.sortedKeys, func(i, j int) bool { return c.sortedKeys[i] < c.sortedKeys[j] })
}

// RemoveNode 删除一个真实节点
func (c *ConsistentHash) RemoveNode(node string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.nodes[node] {
		return
	}

	// 删除虚拟节点
	for i := 0; i < c.replicas; i++ {
		vNodeKey := node + "#" + strconv.Itoa(i)
		hashValue := c.hashKey(vNodeKey)
		delete(c.hashMap, hashValue)

		// 从 sortedKeys 中移除
		for index, key := range c.sortedKeys {
			if key == hashValue {
				c.sortedKeys = append(c.sortedKeys[:index], c.sortedKeys[index+1:]...)
				break
			}
		}
	}

	// 删除真实节点
	delete(c.nodes, node)
}

// GetNode 根据 key 计算哈希并找到最近的服务器
func (c *ConsistentHash) GetNode(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.sortedKeys) == 0 {
		return ""
	}

	hashValue := c.hashKey(key)

	// 通过二分查找找到**顺时针最近的服务器**
	index := sort.Search(len(c.sortedKeys), func(i int) bool {
		return c.sortedKeys[i] >= hashValue
	})

	// 如果 index == len(sortedKeys)，说明 key 的哈希值大于所有节点，应该返回第一个节点
	if index == len(c.sortedKeys) {
		index = 0
	}

	return c.hashMap[c.sortedKeys[index]]
}

// PrintRing 打印哈希环
func (c *ConsistentHash) PrintRing() {
	fmt.Println("------ 哈希环 ------")
	for _, key := range c.sortedKeys {
		fmt.Printf("Hash: %d => Node: %s\n", key, c.hashMap[key])
	}
	fmt.Println("--------------------")
}

// 测试代码
func testConsistent() {
	ch := NewConsistentHash(3) // 每个节点有 3 个虚拟节点

	// 添加节点
	ch.AddNode("Server-A")
	ch.AddNode("Server-B")
	ch.AddNode("Server-C")

	ch.PrintRing()

	// 查找 Key
	keys := []string{"apple", "banana", "cherry", "date", "grape", "melon"}
	for _, key := range keys {
		server := ch.GetNode(key)
		fmt.Printf("Key %s (keyHash %d) is stored in %s\n", key, ch.hashKey(key), server)
	}

	// 删除节点
	fmt.Println("\n--- 删除 Server-B ---")
	ch.RemoveNode("Server-B")
	ch.PrintRing()

	// 重新查找 Key
	for _, key := range keys {
		server := ch.GetNode(key)
		fmt.Printf("Key %s (keyHash %d) is stored in %s\n", key, ch.hashKey(key), server)
	}
}
