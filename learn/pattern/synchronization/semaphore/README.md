# Semaphore Pattern (信号量模式)

## 定义
信号量模式控制对有限资源的并发访问数量。

## 目的
- 限制并发数
- 资源配额管理
- 流量控制

## 使用场景
- 限流控制
- 连接池管理
- 并发数限制
- 资源配额管理

## Go 特有实现
使用 buffered channel 或 `golang.org/x/sync/semaphore`:

```go
type Semaphore struct {
    sem chan struct{}
}

func NewSemaphore(maxConcurrent int) *Semaphore {
    return &Semaphore{
        sem: make(chan struct{}, maxConcurrent),
    }
}

func (s *Semaphore) Acquire() {
    s.sem <- struct{}{}
}

func (s *Semaphore) Release() {
    <-s.sem
}
```

## 优点
1. **并发控制** - 有效控制并发数
2. **资源保护** - 保护有限资源
3. **简单实用** - 实现简单易懂

## 缺点
1. **可能死锁** - 使用不当可能死锁
2. **公平性** - 不保证公平性
