package main

import "container/list"

type LRUCache struct {
	capacity int
	cache    map[int]*list.Element
	list     *list.List // 双向链表：头是最近用的，尾是最久没用的
}

type pair struct {
	key int
	val int
}

func Constructor(capacity int) LRUCache {
	return LRUCache{
		capacity: capacity,
		cache:    make(map[int]*list.Element),
		list:     list.New(),
	}
}

func (c *LRUCache) Get(key int) int {
	if node, ok := c.cache[key]; ok {
		c.list.MoveToFront(node)
		return node.Value.(pair).val
	}
	return -1
}

func (c *LRUCache) Put(key, value int) {
	if node, ok := c.cache[key]; ok {
		c.list.MoveToFront(node)
		node.Value = pair{key, value}
	} else {
		if c.list.Len() == c.capacity {
			back := c.list.Back()
			c.list.Remove(back)
			delete(c.cache, back.Value.(pair).key)
		}
		node := c.list.PushFront(pair{key, value})
		c.cache[key] = node
	}
}

func main() {
	// Example usage
	cache := Constructor(2)
	cache.Put(1, 1)
	cache.Put(2, 2)
	println(cache.Get(1)) // returns 1
	cache.Put(3, 3)       // evicts key 2
	println(cache.Get(2)) // returns -1 (not found)
	cache.Put(4, 4)       // evicts key 1
	println(cache.Get(1)) // returns -1 (not found)
	println(cache.Get(3)) // returns 3
	println(cache.Get(4)) // returns 4
}
