# 消息传递模式 (Messaging Patterns)

消息传递模式提供了进程间通信和异步处理的机制，帮助构建松耦合、可扩展的系统。这些模式在分布式系统、事件驱动架构中广泛应用。

## 📋 模式列表

### 1. [Future/Promise (期约模式)](./futures_promises/)

**目的**: 表示异步操作的最终结果，允许程序在等待结果的同时继续执行其他任务。

**使用场景**:
- 异步 I/O 操作
- 并行计算
- 远程服务调用
- 长时间运行的任务

**Go 特有实现**: 使用 channel 实现 Future，支持链式调用和超时处理。

**示例**: ✅ 已实现

**核心概念**:
- **Future**: 代表一个尚未完成的计算结果
- **Promise**: 可以设置 Future 的值
- **链式调用**: 通过 Then 方法组合多个异步操作

---

### 2. [Push/Pull (推拉模式)](./push_pull/)

**目的**: 定义生产者和消费者之间的消息传递方式，支持推送（生产者主动）和拉取（消费者主动）两种模式。

**使用场景**:
- 消息队列
- 事件总线
- 数据流处理
- 发布-订阅系统

**Go 特有实现**: 使用 buffered channel 实现消息队列，天然支持推拉模式。

**示例**: ✅ 已实现

**两种模式对比**:
- **Push**: 生产者控制，消费者被动接收
- **Pull**: 消费者控制，按需获取数据

---

## 🎯 学习顺序建议

1. **Future/Promise** - 理解异步编程和延迟计算
2. **Push/Pull** - 掌握消息传递的两种基本模式

## 💡 Go 语言实现要点

### 1. Future/Promise 模式

基本实现框架：

```go
type Future struct {
    result chan interface{}
    err    chan error
}

func NewFuture() *Future {
    return &Future{
        result: make(chan interface{}, 1),
        err:    make(chan error, 1),
    }
}

func (f *Future) Get() (interface{}, error) {
    select {
    case res := <-f.result:
        return res, nil
    case err := <-f.err:
        return nil, err
    }
}

// 支持链式调用
func (f *Future) Then(fn func(interface{}) interface{}) *Future {
    nextFuture := NewFuture()
    go func() {
        res, err := f.Get()
        if err != nil {
            nextFuture.SetError(err)
            return
        }
        nextFuture.SetResult(fn(res))
    }()
    return nextFuture
}
```

### 2. Push/Pull 模式

使用 channel 实现：

```go
type Producer struct {
    messages chan Message
}

// Push: 生产者主动推送
func (p *Producer) Push(msg Message) {
    p.messages <- msg
}

type Consumer struct {
    messages chan Message
}

// Pull: 消费者主动拉取
func (c *Consumer) Pull() Message {
    return <-c.messages
}
```

## 🔄 模式对比

| 模式 | 控制方 | 优点 | 缺点 | 适用场景 |
|------|--------|------|------|----------|
| Future/Promise | 调用者 | 非阻塞，链式调用 | 复杂度高 | 异步操作 |
| Push | 生产者 | 实时性好 | 消费者压力大 | 实时通知 |
| Pull | 消费者 | 流量可控 | 可能延迟 | 批处理 |

## 🎓 实际应用场景

### Future/Promise 应用

1. **异步 HTTP 请求**
```go
func FetchURL(url string) *Future {
    future := NewFuture()
    go func() {
        resp, err := http.Get(url)
        if err != nil {
            future.SetError(err)
            return
        }
        body, _ := ioutil.ReadAll(resp.Body)
        future.SetResult(string(body))
    }()
    return future
}

// 使用
future := FetchURL("https://api.example.com/data")
result, err := future.GetWithTimeout(5000)
```

2. **并行计算**
```go
func ParallelCompute(tasks []func() interface{}) []*Future {
    futures := make([]*Future, len(tasks))
    for i, task := range tasks {
        futures[i] = NewFuture()
        go func(f *Future, t func() interface{}) {
            f.SetResult(t())
        }(futures[i], task)
    }
    return futures
}
```

### Push/Pull 应用

1. **消息队列**
```go
// Push 模式：实时推送
producer.Push(Message{Content: "urgent notification"})

// Pull 模式：批量处理
for i := 0; i < batchSize; i++ {
    msg := consumer.Pull()
    processBatch(msg)
}
```

2. **事件系统**
```go
// WebSocket Push
go func() {
    for event := range eventChan {
        websocket.WriteJSON(event)  // Push 给客户端
    }
}()

// HTTP Polling Pull
func (h *Handler) PollEvents(w http.ResponseWriter, r *http.Request) {
    events := eventQueue.Pull(10)  // Pull 最多10个事件
    json.NewEncoder(w).Encode(events)
}
```

## ⚠️ 常见陷阱

### Future/Promise

1. **忘记处理错误**:
   ```go
   // ❌ 错误
   result, _ := future.Get()
   
   // ✅ 正确
   result, err := future.Get()
   if err != nil {
       return err
   }
   ```

2. **无限等待**:
   ```go
   // ❌ 可能永久阻塞
   result, err := future.Get()
   
   // ✅ 使用超时
   result, err := future.GetWithTimeout(5000)
   ```

3. **Goroutine 泄露**:
   ```go
   // ✅ 确保 channel 被消费
   future := NewFuture()
   go func() {
       result := compute()
       future.SetResult(result)  // 必须被 Get() 消费
   }()
   ```

### Push/Pull

1. **缓冲区溢出**:
   ```go
   // ❌ 无缓冲 channel，可能阻塞
   messages := make(chan Message)
   
   // ✅ 合理的缓冲区
   messages := make(chan Message, 100)
   ```

2. **忘记关闭 channel**:
   ```go
   // ✅ 生产者关闭 channel
   defer close(producer.messages)
   
   // ✅ 消费者检查关闭
   for msg := range consumer.messages {
       process(msg)
   }
   ```

## 🎯 最佳实践

### Future/Promise 最佳实践

1. **提供超时机制**:
   ```go
   func (f *Future) GetWithTimeout(ms int) (interface{}, error) {
       select {
       case res := <-f.result:
           return res, nil
       case err := <-f.err:
           return nil, err
       case <-time.After(time.Duration(ms) * time.Millisecond):
           return nil, fmt.Errorf("timeout")
       }
   }
   ```

2. **支持取消操作**:
   ```go
   type CancellableFuture struct {
       *Future
       cancel context.CancelFunc
   }
   
   func (f *CancellableFuture) Cancel() {
       f.cancel()
   }
   ```

3. **错误处理链**:
   ```go
   future.
       Then(processData).
       Then(transformData).
       Catch(func(err error) {
           log.Printf("error: %v", err)
       })
   ```

### Push/Pull 最佳实践

1. **适当的缓冲区大小**:
   ```go
   // 根据生产速率和消费速率设置
   bufferSize := producerRate * avgProcessTime
   messages := make(chan Message, bufferSize)
   ```

2. **优雅关闭**:
   ```go
   // 生产者
   func (p *Producer) Close() {
       close(p.messages)
   }
   
   // 消费者
   for msg := range consumer.messages {
       if msg.IsShutdown() {
           break
       }
       process(msg)
   }
   ```

3. **背压处理**:
   ```go
   // 非阻塞推送
   select {
   case p.messages <- msg:
       // 成功
   default:
       // 队列满，丢弃或记录
       metrics.Dropped.Inc()
   }
   ```

## 🔗 相关模式

### Future/Promise 相关
- **Observer Pattern**: 通知结果就绪
- **Callback Pattern**: Future 可以看作类型安全的 callback
- **Pipeline Pattern**: 链式 Then 调用形成管道

### Push/Pull 相关
- **Producer-Consumer**: Push/Pull 是实现方式
- **Pub-Sub Pattern**: 多消费者的 Push 模式
- **Queue Pattern**: 缓冲的 Push/Pull

## 🚀 性能优化建议

### Future/Promise

1. **复用 Future 对象**:
   ```go
   var futurePool = sync.Pool{
       New: func() interface{} {
           return NewFuture()
       },
   }
   ```

2. **批量操作**:
   ```go
   func WaitAll(futures []*Future) ([]interface{}, error) {
       results := make([]interface{}, len(futures))
       for i, f := range futures {
           res, err := f.Get()
           if err != nil {
               return nil, err
           }
           results[i] = res
       }
       return results, nil
   }
   ```

### Push/Pull

1. **批量推送**:
   ```go
   func (p *Producer) PushBatch(msgs []Message) {
       for _, msg := range msgs {
           p.messages <- msg
       }
   }
   ```

2. **批量拉取**:
   ```go
   func (c *Consumer) PullBatch(n int) []Message {
       msgs := make([]Message, 0, n)
       for i := 0; i < n; i++ {
           select {
           case msg := <-c.messages:
               msgs = append(msgs, msg)
           default:
               return msgs
           }
       }
       return msgs
   }
   ```

## 📚 参考资源

- [Futures and Promises - Wikipedia](https://en.wikipedia.org/wiki/Futures_and_promises)
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Messaging Patterns - Enterprise Integration Patterns](https://www.enterpriseintegrationpatterns.com/)
