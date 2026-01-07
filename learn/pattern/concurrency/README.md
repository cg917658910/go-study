# 并发模式 (Concurrency Patterns)

并发模式是 Go 语言的特色之一，利用 goroutine 和 channel 实现高效的并发编程。这些模式帮助开发者构建可扩展的并发应用。

## 📋 模式列表

### 1. [Active Object (活动对象模式)](./active_object/)

**目的**: 将方法调用与方法执行解耦，每个对象在自己的控制线程中运行。

**使用场景**:
- 异步任务处理
- Actor 模型实现
- 消息队列处理
- 后台服务

**Go 特有实现**: 使用 goroutine 和 channel 实现消息队列。

**示例**: ⏳ 待实现

---

### 2. [Monitor Object (监视对象模式)](./monitor_object/)

**目的**: 同步对共享资源的并发访问，确保在任何时刻只有一个线程可以执行对象的方法。

**使用场景**:
- 线程安全的计数器
- 资源池管理
- 共享状态管理
- 银行账户操作

**Go 特有实现**: 使用 mutex 或 channel 实现同步。

**示例**: ✅ 已实现（位于 synchronization 目录）

---

### 3. [Thread Pool (线程池模式)](./thread_pool/)

**目的**: 维护多个线程等待分配任务，避免频繁创建和销毁线程的开销。

**使用场景**:
- Web 服务器请求处理
- 批量任务处理
- 并发下载
- 数据并行处理

**Go 特有实现**: 使用固定数量的 goroutine 和任务 channel。

**示例**: ⏳ 待实现

---

### 4. [Producer-Consumer (生产者-消费者模式)](./producer_consumer/)

**目的**: 通过缓冲区解耦生产者和消费者，平衡不同速度的生产和消费。

**使用场景**:
- 日志处理系统
- 消息队列
- 数据处理管道
- 爬虫系统

**Go 特有实现**: 使用 buffered channel 作为缓冲区。

**示例**: ✅ 已实现

---

### 5. [Pipeline (管道模式)](./pipeline/)

**目的**: 将数据处理分解为一系列阶段，每个阶段由独立的 goroutine 处理。

**使用场景**:
- 数据处理流水线
- 图像处理管道
- 日志分析
- ETL（提取、转换、加载）

**Go 特有实现**: 使用多个 channel 连接各个处理阶段。

**示例**: ⏳ 待实现

---

### 6. [Fan-Out/Fan-In (扇出/扇入模式)](./fan_out_fan_in/)

**目的**: Fan-Out 将任务分发给多个 worker，Fan-In 收集所有 worker 的结果。

**使用场景**:
- 并行计算
- 多路数据聚合
- 分布式任务处理
- 并发 API 调用

**Go 特有实现**: 使用多个 goroutine 和 select 语句。

**示例**: ⏳ 待实现

---

### 7. [Worker Pool (工作池模式)](./worker_pool/)

**目的**: 创建固定数量的 worker goroutine 处理任务队列。

**使用场景**:
- 限制并发数
- 任务队列处理
- 资源受限的并发处理
- 批量操作

**Go 特有实现**: 使用 semaphore 或固定数量的 goroutine。

**示例**: ⏳ 待实现

---

### 8. [Generator (生成器模式)](./generator/)

**目的**: 按需生成数据序列，而不是一次性生成所有数据。

**使用场景**:
- 无限序列生成
- 惰性求值
- 数据流生成
- 测试数据生成

**Go 特有实现**: 使用 channel 和 goroutine 实现。

**示例**: ✅ 已实现

---

## 🎯 学习顺序建议

1. **Producer-Consumer** - 最基础的并发模式，理解 channel 的使用
2. **Pipeline** - 学习数据流处理
3. **Worker Pool** - 掌握并发控制
4. **Fan-Out/Fan-In** - 理解任务分发和结果聚合
5. **Generator** - 学习惰性求值
6. **Active Object** - 理解异步方法调用
7. **Thread Pool** - 掌握协程池管理

## 💡 Go 语言实现要点

### 1. 生产者-消费者模式
使用 buffered channel:

```go
func ProducerConsumer() {
    buffer := make(chan int, 10)
    
    // 生产者
    go func() {
        for i := 0; i < 100; i++ {
            buffer <- i
        }
        close(buffer)
    }()
    
    // 消费者
    for item := range buffer {
        fmt.Println(item)
    }
}
```

### 2. Pipeline 模式
连接多个处理阶段：

```go
func Generator(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

func Square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

// 使用
func main() {
    c := Generator(2, 3, 4)
    out := Square(c)
    for result := range out {
        fmt.Println(result)
    }
}
```

### 3. Worker Pool 模式
固定数量的 worker:

```go
func WorkerPool(tasks <-chan Task, numWorkers int) {
    var wg sync.WaitGroup
    
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for task := range tasks {
                task.Execute()
            }
        }()
    }
    
    wg.Wait()
}
```

### 4. Fan-Out/Fan-In 模式
任务分发和结果收集：

```go
func FanOut(in <-chan int, workers int) []<-chan int {
    channels := make([]<-chan int, workers)
    for i := 0; i < workers; i++ {
        channels[i] = worker(in)
    }
    return channels
}

func FanIn(channels ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    
    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for n := range c {
                out <- n
            }
        }(ch)
    }
    
    go func() {
        wg.Wait()
        close(out)
    }()
    
    return out
}
```

## 🔄 模式对比

| 模式 | 关注点 | 主要用途 | 复杂度 |
|------|--------|----------|--------|
| Active Object | 异步执行 | 方法调用解耦 | ⭐⭐⭐ |
| Monitor Object | 同步访问 | 线程安全 | ⭐⭐ |
| Thread Pool | 资源复用 | 限制并发数 | ⭐⭐ |
| Producer-Consumer | 解耦生产消费 | 速度匹配 | ⭐⭐ |
| Pipeline | 数据流处理 | 阶段化处理 | ⭐⭐ |
| Fan-Out/Fan-In | 并行处理 | 任务分发聚合 | ⭐⭐⭐ |
| Worker Pool | 任务处理 | 并发控制 | ⭐⭐ |
| Generator | 数据生成 | 惰性求值 | ⭐⭐ |

## 📚 相关模式

- **Pipeline + Fan-Out/Fan-In**: 可以在 Pipeline 的某个阶段使用 Fan-Out 并行处理
- **Worker Pool + Producer-Consumer**: Worker Pool 是 Producer-Consumer 的特例
- **Active Object + Command**: Active Object 可以使用 Command 模式封装请求
- **Generator + Iterator**: Generator 是一种特殊的 Iterator

## ⚠️ 常见陷阱

1. **Goroutine 泄漏**: 忘记关闭 channel 或等待 goroutine 结束
2. **死锁**: Channel 读写不匹配导致死锁
3. **竞态条件**: 共享变量未正确同步
4. **过度并发**: 创建太多 goroutine 反而降低性能
5. **Channel 缓冲区大小**: 缓冲区太小导致阻塞，太大浪费内存

## 🎓 最佳实践

1. **关闭 channel 的责任**: 由发送方关闭 channel
2. **使用 sync.WaitGroup**: 等待所有 goroutine 完成
3. **context 取消**: 使用 context 控制 goroutine 生命周期
4. **有限的并发数**: 使用 Worker Pool 限制并发
5. **select 的 default**: 非阻塞的 channel 操作
6. **检测 goroutine 泄漏**: 使用工具检测未结束的 goroutine

## 🚀 性能优化建议

1. **合理设置 buffer 大小**: 根据生产和消费速度调整
2. **减少 channel 操作**: 批量处理减少 channel 通信
3. **避免频繁创建 goroutine**: 使用 worker pool
4. **使用 sync.Pool**: 复用对象减少 GC 压力
5. **profile 分析**: 使用 pprof 分析并发性能瓶颈
