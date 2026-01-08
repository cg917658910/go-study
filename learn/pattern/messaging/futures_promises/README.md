# Future/Promise Pattern (期约模式)

## 定义
Future/Promise 模式表示一个异步操作的最终结果。Future 是一个占位符，代表一个可能还未完成的计算结果。程序可以在等待结果的同时继续执行其他任务，这是异步编程的核心概念。

## 目的
- 实现非阻塞的异步操作
- 延迟计算和懒加载
- 支持链式异步调用
- 提供超时和错误处理机制

## 使用场景
- **异步 I/O 操作**: 网络请求、文件读写
- **并行计算**: 多任务并发执行
- **远程服务调用**: RPC、HTTP API 调用
- **长时间运行的任务**: 数据处理、图像渲染
- **事件驱动编程**: 等待事件发生

## 核心概念

### Future vs Promise

虽然在 Go 实现中通常合并为一个类型，但概念上：

- **Future (未来)**: 只读的结果占位符，用于获取结果
- **Promise (承诺)**: 可写的结果提供者，用于设置结果

```
Promise (生产者)  ──设置结果──>  Future (消费者)
```

### 状态转换

```
Pending (等待中)
    │
    ├─> Fulfilled (已完成) ──> 返回结果
    │
    └─> Rejected (已拒绝)  ──> 返回错误
```

## Go 特有实现

使用 channel 实现 Future：

```go
type Future struct {
    result chan interface{}  // 结果通道
    err    chan error         // 错误通道
}

func NewFuture() *Future {
    return &Future{
        result: make(chan interface{}, 1),  // 缓冲1，防止阻塞
        err:    make(chan error, 1),
    }
}

// 设置结果（Promise 部分）
func (f *Future) SetResult(result interface{}) {
    f.result <- result
}

// 设置错误（Promise 部分）
func (f *Future) SetError(err error) {
    f.err <- err
}

// 获取结果（Future 部分）
func (f *Future) Get() (interface{}, error) {
    select {
    case res := <-f.result:
        return res, nil
    case err := <-f.err:
        return nil, err
    }
}
```

## 优点
1. **非阻塞执行** - 不等待结果，可以继续执行其他任务
2. **链式调用** - 通过 Then 方法组合多个异步操作
3. **统一的错误处理** - Catch 方法集中处理错误
4. **超时控制** - 防止无限等待
5. **代码清晰** - 比回调函数更易读

## 缺点
1. **复杂度增加** - 需要理解异步编程概念
2. **调试困难** - 异步执行使调试更复杂
3. **内存开销** - 每个 Future 需要额外的 channel
4. **可能的 Goroutine 泄露** - 未正确处理可能导致泄露

## 高级特性

### 1. 超时处理

```go
func (f *Future) GetWithTimeout(timeout int) (interface{}, error) {
    select {
    case res := <-f.result:
        return res, nil
    case err := <-f.err:
        return nil, err
    case <-time.After(time.Duration(timeout) * time.Millisecond):
        return nil, fmt.Errorf("timeout occurred while waiting for result")
    }
}
```

### 2. 链式调用 (Then)

```go
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

// 使用示例
future.
    Then(func(data interface{}) interface{} {
        return processData(data)
    }).
    Then(func(data interface{}) interface{} {
        return transformData(data)
    })
```

### 3. 错误处理 (Catch)

```go
func (f *Future) Catch(fn func(error)) *Future {
    go func() {
        _, err := f.Get()
        if err != nil {
            fn(err)
        }
    }()
    return f
}

// 使用示例
future.Catch(func(err error) {
    log.Printf("Error occurred: %v", err)
})
```

## 实际应用示例

### 1. 异步 HTTP 请求

```go
func FetchURLAsync(url string) *Future {
    future := NewFuture()
    go func() {
        resp, err := http.Get(url)
        if err != nil {
            future.SetError(err)
            return
        }
        defer resp.Body.Close()
        
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            future.SetError(err)
            return
        }
        future.SetResult(string(body))
    }()
    return future
}

// 使用
future := FetchURLAsync("https://api.example.com/data")
result, err := future.GetWithTimeout(5000)
if err != nil {
    log.Fatal(err)
}
fmt.Println(result)
```

### 2. 并行任务执行

```go
func WaitAll(futures ...*Future) ([]interface{}, error) {
    results := make([]interface{}, len(futures))
    for i, future := range futures {
        res, err := future.Get()
        if err != nil {
            return nil, err
        }
        results[i] = res
    }
    return results, nil
}

// 使用
future1 := FetchURLAsync("https://api1.example.com/data")
future2 := FetchURLAsync("https://api2.example.com/data")
future3 := FetchURLAsync("https://api3.example.com/data")

results, err := WaitAll(future1, future2, future3)
```

### 3. 数据库异步查询

```go
func QueryAsync(db *sql.DB, query string) *Future {
    future := NewFuture()
    go func() {
        rows, err := db.Query(query)
        if err != nil {
            future.SetError(err)
            return
        }
        defer rows.Close()
        
        var results []Row
        for rows.Next() {
            var row Row
            if err := rows.Scan(&row); err != nil {
                future.SetError(err)
                return
            }
            results = append(results, row)
        }
        future.SetResult(results)
    }()
    return future
}
```

### 4. 文件异步读取

```go
func ReadFileAsync(filename string) *Future {
    future := NewFuture()
    go func() {
        data, err := ioutil.ReadFile(filename)
        if err != nil {
            future.SetError(err)
            return
        }
        future.SetResult(string(data))
    }()
    return future
}

// 链式处理
ReadFileAsync("config.json").
    Then(func(data interface{}) interface{} {
        var config Config
        json.Unmarshal([]byte(data.(string)), &config)
        return config
    }).
    Then(func(config interface{}) interface{} {
        return validateConfig(config.(Config))
    }).
    Catch(func(err error) {
        log.Printf("Failed to load config: %v", err)
    })
```

## 与 Go Context 结合

```go
type ContextFuture struct {
    *Future
    ctx    context.Context
    cancel context.CancelFunc
}

func NewContextFuture(ctx context.Context) *ContextFuture {
    ctx, cancel := context.WithCancel(ctx)
    return &ContextFuture{
        Future: NewFuture(),
        ctx:    ctx,
        cancel: cancel,
    }
}

func (f *ContextFuture) Get() (interface{}, error) {
    select {
    case res := <-f.result:
        return res, nil
    case err := <-f.err:
        return nil, err
    case <-f.ctx.Done():
        return nil, f.ctx.Err()
    }
}

func (f *ContextFuture) Cancel() {
    f.cancel()
}
```

## 最佳实践

### 1. 总是设置超时

```go
// ❌ 可能永久阻塞
result, err := future.Get()

// ✅ 使用超时
result, err := future.GetWithTimeout(5000)
```

### 2. 正确处理错误

```go
// ✅ 完整的错误处理
future := AsyncOperation()
result, err := future.GetWithTimeout(3000)
if err != nil {
    switch err.Error() {
    case "timeout occurred while waiting for result":
        // 处理超时
    default:
        // 处理其他错误
    }
    return err
}
```

### 3. 避免 Goroutine 泄露

```go
// ✅ 确保 Future 被消费
future := NewFuture()
go func() {
    result := compute()
    select {
    case future.result <- result:
        // 成功发送
    case <-time.After(time.Minute):
        // 超时，放弃
        log.Println("Future result abandoned")
    }
}()
```

### 4. 使用带缓冲的 Channel

```go
// ✅ 缓冲为 1，防止 goroutine 阻塞
type Future struct {
    result chan interface{} // buffered
    err    chan error        // buffered
}

func NewFuture() *Future {
    return &Future{
        result: make(chan interface{}, 1),
        err:    make(chan error, 1),
    }
}
```

## 性能优化

### 1. Future 对象池

```go
var futurePool = sync.Pool{
    New: func() interface{} {
        return &Future{
            result: make(chan interface{}, 1),
            err:    make(chan error, 1),
        }
    },
}

func GetFuture() *Future {
    return futurePool.Get().(*Future)
}

func PutFuture(f *Future) {
    // 清空 channel
    select {
    case <-f.result:
    default:
    }
    select {
    case <-f.err:
    default:
    }
    futurePool.Put(f)
}
```

### 2. 批量等待优化

```go
func WaitAllFast(futures []*Future) ([]interface{}, error) {
    results := make([]interface{}, len(futures))
    errs := make([]error, len(futures))
    
    var wg sync.WaitGroup
    for i, future := range futures {
        wg.Add(1)
        go func(idx int, f *Future) {
            defer wg.Done()
            results[idx], errs[idx] = f.Get()
        }(i, future)
    }
    wg.Wait()
    
    for _, err := range errs {
        if err != nil {
            return nil, err
        }
    }
    return results, nil
}
```

## 与其他模式的关系

- **Callback Pattern**: Future 是类型安全的 callback 替代
- **Observer Pattern**: Future 可以通知结果就绪
- **Pipeline Pattern**: Then 链形成数据处理管道
- **Promise/A+**: JavaScript Promise 的 Go 实现

## 参考资源

- [Futures and Promises - Wikipedia](https://en.wikipedia.org/wiki/Futures_and_promises)
- [JavaScript Promises](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise)
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
