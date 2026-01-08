# Push/Pull Pattern (推拉模式)

## 定义
Push/Pull 模式定义了生产者和消费者之间的消息传递方式。在 Push 模式中，生产者主动将消息推送给消费者；在 Pull 模式中，消费者主动从生产者拉取消息。

## 目的
- 解耦生产者和消费者
- 支持不同的消息传递策略
- 控制消息流量和处理速率
- 实现灵活的消息传递机制

## 使用场景

### Push 模式适用于：
- **实时通知**: WebSocket、Server-Sent Events
- **事件驱动**: 事件总线、发布-订阅
- **流式处理**: 实时数据流
- **高时效性**: 需要立即处理的消息

### Pull 模式适用于：
- **批量处理**: 按批次处理消息
- **流量控制**: 消费者控制处理速率
- **定时任务**: 定期拉取数据
- **负载均衡**: 多消费者竞争拉取

## 两种模式对比

| 特性 | Push 模式 | Pull 模式 |
|------|-----------|-----------|
| 控制方 | 生产者 | 消费者 |
| 实时性 | 高 | 中等 |
| 流量控制 | 困难 | 容易 |
| 消费者压力 | 可能过载 | 可控 |
| 适用场景 | 实时通知 | 批量处理 |
| 实现复杂度 | 简单 | 中等 |

## Go 特有实现

### Push 模式

生产者主动推送消息：

```go
type Producer struct {
    messages chan Message
}

func NewProducer(bufferSize int) *Producer {
    return &Producer{
        messages: make(chan Message, bufferSize),
    }
}

// Push: 生产者推送消息
func (p *Producer) Push(msg Message) {
    p.messages <- msg
}
```

### Pull 模式

消费者主动拉取消息：

```go
type Consumer struct {
    messages chan Message
}

func NewConsumer(messages chan Message) *Consumer {
    return &Consumer{
        messages: messages,
    }
}

// Pull: 消费者拉取消息
func (c *Consumer) Pull() Message {
    return <-c.messages
}
```

## 优点

### Push 模式优点
1. **实时性好** - 消息立即推送给消费者
2. **实现简单** - 直接发送到 channel
3. **低延迟** - 无需轮询等待
4. **资源利用** - 无需消费者主动轮询

### Pull 模式优点
1. **流量可控** - 消费者按需拉取
2. **避免过载** - 消费者控制处理速率
3. **批量处理** - 可以批量拉取消息
4. **灵活性高** - 消费者决定何时拉取

## 缺点

### Push 模式缺点
1. **消费者压力** - 可能被消息淹没
2. **流量控制难** - 生产者不知道消费者状态
3. **背压问题** - 需要处理消息堆积

### Pull 模式缺点
1. **延迟较高** - 需要轮询或定时拉取
2. **资源浪费** - 空闲时也可能轮询
3. **复杂度高** - 需要实现拉取逻辑

## 实际应用示例

### 1. Push 模式 - 实时通知系统

```go
type NotificationService struct {
    subscribers map[string]chan Notification
    mu          sync.RWMutex
}

func (s *NotificationService) Subscribe(userID string) chan Notification {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    ch := make(chan Notification, 10)
    s.subscribers[userID] = ch
    return ch
}

// Push: 推送通知给所有订阅者
func (s *NotificationService) Notify(notification Notification) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    for _, ch := range s.subscribers {
        select {
        case ch <- notification:
            // 成功推送
        default:
            // 消费者太慢，丢弃消息或记录
            log.Println("Notification dropped")
        }
    }
}
```

### 2. Pull 模式 - 任务队列

```go
type TaskQueue struct {
    tasks chan Task
}

func NewTaskQueue(size int) *TaskQueue {
    return &TaskQueue{
        tasks: make(chan Task, size),
    }
}

// Worker Pull: 工作者拉取任务
func (q *TaskQueue) Worker(id int) {
    for {
        task := <-q.tasks
        if task.IsShutdown() {
            break
        }
        
        fmt.Printf("Worker %d processing task %d\n", id, task.ID)
        task.Process()
    }
}

// 启动多个 worker
func (q *TaskQueue) Start(numWorkers int) {
    for i := 0; i < numWorkers; i++ {
        go q.Worker(i)
    }
}
```

### 3. 混合模式 - WebSocket

```go
type WebSocketServer struct {
    clients   map[*Client]bool
    broadcast chan Message  // Push: 广播通道
    mu        sync.RWMutex
}

type Client struct {
    conn     *websocket.Conn
    send     chan Message  // Push: 发送通道
    receive  chan Message  // Pull: 接收通道
}

// Push: 服务器推送消息给客户端
func (c *Client) writePump() {
    for msg := range c.send {
        c.conn.WriteJSON(msg)
    }
}

// Pull: 客户端从连接拉取消息
func (c *Client) readPump() {
    for {
        var msg Message
        if err := c.conn.ReadJSON(&msg); err != nil {
            break
        }
        c.receive <- msg
    }
}
```

### 4. 批量拉取优化

```go
type BatchConsumer struct {
    messages chan Message
}

// Pull: 批量拉取
func (c *BatchConsumer) PullBatch(maxSize int, timeout time.Duration) []Message {
    batch := make([]Message, 0, maxSize)
    timer := time.NewTimer(timeout)
    defer timer.Stop()
    
    for i := 0; i < maxSize; i++ {
        select {
        case msg := <-c.messages:
            batch = append(batch, msg)
        case <-timer.C:
            return batch  // 超时返回已拉取的消息
        }
    }
    return batch
}

// 使用批量拉取
func (c *BatchConsumer) Process() {
    for {
        batch := c.PullBatch(100, time.Second)
        if len(batch) == 0 {
            continue
        }
        
        // 批量处理
        processBatch(batch)
    }
}
```

## 背压处理（Backpressure）

### Push 模式的背压处理

```go
// 1. 使用非阻塞发送
func (p *Producer) PushNonBlocking(msg Message) bool {
    select {
    case p.messages <- msg:
        return true
    default:
        // 队列满，拒绝消息
        return false
    }
}

// 2. 使用带超时的发送
func (p *Producer) PushWithTimeout(msg Message, timeout time.Duration) error {
    select {
    case p.messages <- msg:
        return nil
    case <-time.After(timeout):
        return errors.New("send timeout")
    }
}

// 3. 动态调整缓冲区
type DynamicProducer struct {
    messages chan Message
    mu       sync.RWMutex
}

func (p *DynamicProducer) Resize(newSize int) {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    newChan := make(chan Message, newSize)
    close(p.messages)
    for msg := range p.messages {
        newChan <- msg
    }
    p.messages = newChan
}
```

### Pull 模式的流量控制

```go
// 令牌桶限流
type RateLimitedConsumer struct {
    messages chan Message
    limiter  *rate.Limiter
}

func (c *RateLimitedConsumer) Pull() (Message, error) {
    if err := c.limiter.Wait(context.Background()); err != nil {
        return Message{}, err
    }
    return <-c.messages, nil
}
```

## 最佳实践

### 1. 合理设置缓冲区大小

```go
// ✅ 根据生产和消费速率计算
producerRate := 1000  // 每秒1000条消息
consumerRate := 800   // 每秒800条消息
avgDelay := 2         // 平均延迟2秒

bufferSize := (producerRate - consumerRate) * avgDelay
messages := make(chan Message, bufferSize)
```

### 2. 优雅关闭

```go
type GracefulProducer struct {
    messages chan Message
    done     chan struct{}
}

func (p *GracefulProducer) Close() {
    close(p.done)
    close(p.messages)
}

func (p *GracefulProducer) Push(msg Message) bool {
    select {
    case p.messages <- msg:
        return true
    case <-p.done:
        return false
    }
}
```

### 3. 监控指标

```go
type MonitoredQueue struct {
    messages   chan Message
    pushCount  int64
    pullCount  int64
    dropCount  int64
}

func (q *MonitoredQueue) Push(msg Message) {
    select {
    case q.messages <- msg:
        atomic.AddInt64(&q.pushCount, 1)
    default:
        atomic.AddInt64(&q.dropCount, 1)
    }
}

func (q *MonitoredQueue) Pull() Message {
    msg := <-q.messages
    atomic.AddInt64(&q.pullCount, 1)
    return msg
}

func (q *MonitoredQueue) Stats() (push, pull, drop int64) {
    return atomic.LoadInt64(&q.pushCount),
           atomic.LoadInt64(&q.pullCount),
           atomic.LoadInt64(&q.dropCount)
}
```

### 4. 错误处理

```go
// Push 错误处理
func (p *Producer) PushWithRetry(msg Message, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        select {
        case p.messages <- msg:
            return nil
        case <-time.After(time.Second):
            log.Printf("Retry %d/%d", i+1, maxRetries)
        }
    }
    return errors.New("max retries exceeded")
}

// Pull 错误处理
func (c *Consumer) PullWithTimeout(timeout time.Duration) (Message, error) {
    select {
    case msg := <-c.messages:
        return msg, nil
    case <-time.After(timeout):
        return Message{}, errors.New("pull timeout")
    }
}
```

## 常见陷阱

### 1. 忘记关闭 Channel

```go
// ❌ 错误：忘记关闭
producer := NewProducer(10)
// ... 使用后未关闭

// ✅ 正确：使用 defer
producer := NewProducer(10)
defer close(producer.messages)
```

### 2. Channel 已关闭仍然写入

```go
// ❌ 错误：可能 panic
p.messages <- msg

// ✅ 正确：检查是否关闭
select {
case p.messages <- msg:
    // 成功
case <-p.done:
    // 已关闭
}
```

### 3. 无缓冲 Channel 导致阻塞

```go
// ❌ 可能阻塞
messages := make(chan Message)

// ✅ 使用缓冲
messages := make(chan Message, 100)
```

## 与其他模式的关系

- **Producer-Consumer**: Push/Pull 是具体实现方式
- **Pub-Sub**: 多消费者的 Push 模式
- **Observer Pattern**: Push 模式的一种应用
- **Queue Pattern**: 提供缓冲的 Push/Pull

## 参考资源

- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Reactive Streams](http://www.reactive-streams.org/)
- [Enterprise Integration Patterns](https://www.enterpriseintegrationpatterns.com/)
