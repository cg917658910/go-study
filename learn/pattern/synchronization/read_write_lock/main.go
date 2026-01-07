package main

import (
"fmt"
"sync"
)

type Cache struct {
mu   sync.RWMutex
data map[string]string
}

func NewCache() *Cache {
return &Cache{
data: make(map[string]string),
}
}

func (c *Cache) Get(key string) (string, bool) {
c.mu.RLock()
defer c.mu.RUnlock()
val, ok := c.data[key]
return val, ok
}

func (c *Cache) Set(key, value string) {
c.mu.Lock()
defer c.mu.Unlock()
c.data[key] = value
}

func main() {
cache := NewCache()
cache.Set("key1", "value1")

val, ok := cache.Get("key1")
if ok {
fmt.Printf("Value: %s\n", val)
}
}
