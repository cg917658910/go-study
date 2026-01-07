# 同步模式 (Synchronization Patterns)

同步模式提供了线程同步和资源访问控制的机制。这些模式确保在并发环境中正确地访问共享资源。

## 📋 模式列表

### 1. [Mutex (互斥锁模式)](./mutex/)

**目的**: 提供互斥访问机制，确保同一时刻只有一个 goroutine 可以访问共享资源。

**使用场景**:
- 保护共享变量
- 临界区保护
- 单例模式的线程安全实现
- 计数器的原子操作

**Go 特有实现**: 使用 `sync.Mutex` 或 `sync.RWMutex`。

**示例**: ⏳ 待实现

---

### 2. [Semaphore (信号量模式)](./semaphore/)

**目的**: 控制对有限资源的并发访问数量。

**使用场景**:
- 限流控制
- 连接池管理
- 并发数限制
- 资源配额管理

**Go 特有实现**: 使用 buffered channel 或 `golang.org/x/sync/semaphore`。

**示例**: ✅ 已实现

---

### 3. [Barrier (屏障模式)](./barrier/)

**目的**: 让一组 goroutine 在某个点等待，直到所有 goroutine 都到达该点后才继续执行。

**使用场景**:
- 阶段性同步
- 并行算法的阶段分隔
- 多任务协同
- 测试场景同步

**Go 特有实现**: 使用 `sync.WaitGroup` 或 channel。

**示例**: ✅ 已实现（n_barrier）

---

### 4. [Read-Write Lock (读写锁模式)](./read_write_lock/)

**目的**: 允许多个读操作并发执行，但写操作独占访问。

**使用场景**:
- 缓存系统
- 配置管理
- 读多写少的数据结构
- 共享数据的并发访问

**Go 特有实现**: 使用 `sync.RWMutex`。

**示例**: ⏳ 待实现

---

### 5. [Condition Variable (条件变量模式)](./condition_variable/)

**目的**: 允许 goroutine 等待某个条件成立，避免忙等待。

**使用场景**:
- 生产者-消费者问题
- 资源可用性等待
- 事件通知
- 队列的空/满条件

**Go 特有实现**: 使用 `sync.Cond` 或 channel。

**示例**: ⏳ 待实现

---

### 6. [Monitor (监视器模式)](./monitor/)

**目的**: 封装共享数据和同步机制，确保线程安全的访问。

**使用场景**:
- 线程安全的数据结构
- 银行账户操作
- 资源管理
- 状态管理

**Go 特有实现**: 使用 mutex 和条件变量组合。

**示例**: ✅ 已实现

---

## 🎯 学习顺序建议

1. **Mutex** - 最基础的同步原语，理解互斥锁
2. **Read-Write Lock** - 学习读写分离的优化
3. **Semaphore** - 掌握资源数量控制
4. **Barrier** - 理解协程同步点
5. **Condition Variable** - 学习条件等待机制
6. **Monitor** - 综合应用多种同步原语

## 💡 Go 语言实现要点

### 1. Mutex 模式
基本的互斥锁：

```go
type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (c *SafeCounter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}

func (c *SafeCounter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.count
}
```

### 2. Read-Write Lock 模式
读写锁优化：

```go
type Cache struct {
    mu    sync.RWMutex
    data  map[string]string
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
```

### 3. Semaphore 模式
使用 buffered channel:

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

### 4. Barrier 模式
使用 WaitGroup:

```go
type Barrier struct {
    n  int
    wg sync.WaitGroup
}

func NewBarrier(n int) *Barrier {
    b := &Barrier{n: n}
    b.wg.Add(n)
    return b
}

func (b *Barrier) Wait() {
    b.wg.Done()
    b.wg.Wait()
}
```

### 5. Condition Variable 模式
使用 sync.Cond:

```go
type Queue struct {
    mu    sync.Mutex
    cond  *sync.Cond
    items []interface{}
}

func NewQueue() *Queue {
    q := &Queue{
        items: make([]interface{}, 0),
    }
    q.cond = sync.NewCond(&q.mu)
    return q
}

func (q *Queue) Enqueue(item interface{}) {
    q.mu.Lock()
    defer q.mu.Unlock()
    q.items = append(q.items, item)
    q.cond.Signal()
}

func (q *Queue) Dequeue() interface{} {
    q.mu.Lock()
    defer q.mu.Unlock()
    for len(q.items) == 0 {
        q.cond.Wait()
    }
    item := q.items[0]
    q.items = q.items[1:]
    return item
}
```

### 6. Channel 作为同步原语
Go 惯用方式：

```go
// 使用 channel 实现信号量
sem := make(chan struct{}, maxConcurrent)

// 获取
sem <- struct{}{}
// 释放
<-sem

// 使用 channel 实现屏障
barrier := make(chan struct{})
go func() {
    // 等待所有任务完成
    barrier <- struct{}{}
}()
<-barrier
```

## 🔄 模式对比

| 模式 | 同步粒度 | 适用场景 | 复杂度 |
|------|----------|----------|--------|
| Mutex | 互斥 | 简单的临界区 | ⭐ |
| Read-Write Lock | 读写分离 | 读多写少 | ⭐⭐ |
| Semaphore | 资源计数 | 有限资源控制 | ⭐⭐ |
| Barrier | 同步点 | 阶段同步 | ⭐⭐ |
| Condition Variable | 条件等待 | 复杂等待条件 | ⭐⭐⭐ |
| Monitor | 封装同步 | 线程安全对象 | ⭐⭐⭐ |

## 📚 相关模式

- **Mutex vs Monitor**: Monitor 是对 Mutex 的更高级封装
- **Semaphore vs Mutex**: Semaphore 可以看作是 Mutex 的推广（计数为1时等价）
- **Read-Write Lock vs Mutex**: Read-Write Lock 是 Mutex 的优化版本
- **Barrier vs WaitGroup**: Barrier 通常使用 WaitGroup 实现
- **Condition Variable + Mutex**: 通常组合使用

## ⚠️ 常见陷阱

1. **死锁**: 
   - 互相等待对方释放锁
   - 忘记释放锁
   - 加锁顺序不一致

2. **活锁**: 
   - 协程不断改变状态但无法前进

3. **饥饿**: 
   - 某些协程长时间得不到资源
   - 不公平的锁调度

4. **性能问题**:
   - 锁粒度过大导致并发度低
   - 频繁的锁竞争
   - 不必要的锁使用

5. **忘记 defer unlock**:
   - 导致锁未释放
   - 使用 `defer mu.Unlock()` 确保释放

## 🎓 最佳实践

### 1. 优先使用 Channel
Go 鼓励使用 channel 进行通信而非共享内存：

```go
// 推荐：使用 channel
done := make(chan bool)
go func() {
    // 工作
    done <- true
}()
<-done

// 避免：过度使用锁
var mu sync.Mutex
var done bool
go func() {
    // 工作
    mu.Lock()
    done = true
    mu.Unlock()
}()
```

### 2. 锁的最小化
只保护必要的临界区：

```go
func (c *Counter) Inc() {
    // 准备工作（不需要锁）
    newValue := calculateValue()
    
    // 最小临界区
    c.mu.Lock()
    c.value += newValue
    c.mu.Unlock()
    
    // 后续工作（不需要锁）
    logUpdate()
}
```

### 3. 使用 defer 确保释放
防止忘记释放锁：

```go
func (c *Counter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}
```

### 4. 避免在持有锁时调用外部代码
可能导致死锁：

```go
// 危险
func (c *Counter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
    externalCallback() // 可能导致死锁
}

// 安全
func (c *Counter) Inc() {
    c.mu.Lock()
    c.count++
    c.mu.Unlock()
    externalCallback()
}
```

### 5. 读写锁的适用条件
读多写少时才使用：

```go
// 适合使用 RWMutex
type Cache struct {
    mu   sync.RWMutex
    data map[string]string
}

// 读操作频繁
func (c *Cache) Get(key string) string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.data[key]
}

// 写操作较少
func (c *Cache) Set(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value
}
```

### 6. 使用竞态检测器
开发时启用竞态检测：

```bash
go test -race ./...
go run -race main.go
```

## 🚀 性能优化建议

1. **减少锁竞争**:
   - 减小临界区
   - 分段锁（sharding）
   - 使用原子操作

2. **选择合适的同步原语**:
   - 简单计数器用 `atomic`
   - 读多写少用 `RWMutex`
   - 固定数量资源用 `Semaphore`

3. **避免忙等待**:
   - 使用条件变量而非轮询
   - 使用 channel 的阻塞特性

4. **批量操作**:
   - 减少加锁次数
   - 批量处理数据

5. **无锁数据结构**:
   - 考虑使用 `atomic` 包
   - lock-free 数据结构
