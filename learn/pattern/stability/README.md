# 稳定性模式 (Stability Patterns)

稳定性模式旨在提高系统的可靠性和容错能力，防止故障蔓延，确保系统在面对错误和异常时能够优雅地降级而不是完全崩溃。

## 📋 模式列表

### 1. [Circuit Breaker (熔断器模式)](./circuit_breaker/)

**目的**: 通过监控系统运行状态，在检测到故障时自动切断对故障组件的调用，防止故障蔓延，并在系统恢复后自动恢复调用。

**使用场景**:
- 微服务架构中的服务保护
- 外部 API 调用保护
- 数据库连接保护
- 防止级联故障

**Go 特有实现**: 使用状态机管理三种状态（CLOSED、OPEN、HALF-OPEN），配合 mutex 确保线程安全。

**示例**: ✅ 已实现

**关键状态**:
- **CLOSED**: 正常状态，请求正常通过
- **OPEN**: 熔断状态，快速失败，拒绝请求
- **HALF-OPEN**: 半开状态，允许少量请求测试服务是否恢复

---

### 2. [Bulkheads (隔离舱壁模式)](./bulkheads/)

**目的**: 将系统划分为多个隔离的单元，每个单元都有自己的资源池和限制，防止一个单元的故障影响到其他单元。

**使用场景**:
- 服务资源隔离
- 线程池隔离
- 数据库连接池隔离
- 防止资源耗尽

**Go 特有实现**: 使用 buffered channel 作为信号量，控制并发访问数量。

**示例**: ✅ 已实现

**核心思想**: 类似船舶的隔离舱设计，即使一个舱室进水，也不会导致整艘船沉没。

---

## 🎯 学习顺序建议

1. **Circuit Breaker** - 理解熔断机制和状态转换
2. **Bulkheads** - 掌握资源隔离和并发控制

## 💡 Go 语言实现要点

### 1. Circuit Breaker 模式

基本实现框架：

```go
type CircuitBreaker struct {
    mu               sync.Mutex
    state            string // CLOSED, OPEN, HALF-OPEN
    failureCount     int
    successCount     int
    failureThreshold int
    successThreshold int
    openTimeout      time.Duration
    lastOpened       time.Time
}

func (cb *CircuitBreaker) Execute(task func() error) error {
    cb.mu.Lock()
    defer cb.mu.Unlock()
    
    switch cb.state {
    case "OPEN":
        if time.Since(cb.lastOpened) > cb.openTimeout {
            cb.state = "HALF-OPEN"
        } else {
            return fmt.Errorf("circuit breaker is open")
        }
    case "HALF-OPEN":
        // 允许请求通过，测试服务是否恢复
    case "CLOSED":
        // 正常执行请求
    }
    // 执行任务并更新状态
}
```

### 2. Bulkheads 模式

使用信号量控制并发：

```go
type Bulkhead struct {
    limit     int
    requestCh chan func()
    quitCh    chan struct{}
}

func (b *Bulkhead) start() {
    sem := make(chan struct{}, b.limit)
    for {
        select {
        case req := <-b.requestCh:
            sem <- struct{}{}
            go func(r func()) {
                defer func() { <-sem }()
                r()
            }(req)
        case <-b.quitCh:
            return
        }
    }
}
```

## 🔄 模式对比

| 模式 | 关注点 | 故障处理方式 | 适用场景 |
|------|--------|--------------|----------|
| Circuit Breaker | 故障检测与恢复 | 快速失败，自动恢复 | 外部服务调用 |
| Bulkheads | 资源隔离 | 限制影响范围 | 多租户系统 |

**Circuit Breaker vs Bulkheads**:
- Circuit Breaker 关注**故障检测和恢复**，通过监控错误率来决定是否熔断
- Bulkheads 关注**资源隔离和限制**，通过资源分区来防止资源耗尽

**两者结合使用**:
在实际应用中，这两种模式经常一起使用：
- 使用 Bulkheads 隔离不同服务的资源
- 在每个隔离区内使用 Circuit Breaker 保护外部调用

## ⚠️ 常见陷阱

### Circuit Breaker

1. **阈值设置不当**:
   - 太敏感：正常波动就触发熔断
   - 太迟钝：故障已经蔓延才熔断

2. **超时时间配置**:
   - 太短：服务未恢复就尝试调用
   - 太长：影响用户体验

3. **并发安全**:
   - 必须使用 mutex 保护状态变更

### Bulkheads

1. **资源限制设置**:
   - 太小：正常请求被拒绝
   - 太大：失去隔离效果

2. **goroutine 泄露**:
   - 必须确保 goroutine 正确退出
   - 使用 defer 确保释放资源

3. **死锁风险**:
   - 避免在持有锁时等待 channel

## 🎓 最佳实践

### 1. Circuit Breaker 配置

```go
// 推荐配置
cb := NewCircuitBreaker(
    5,              // 失败阈值：5次失败后熔断
    2,              // 成功阈值：2次成功后恢复
    30*time.Second, // 超时：30秒后尝试恢复
)
```

### 2. Bulkheads 配置

```go
// 根据系统资源合理配置
bulkhead := NewBulkhead(
    runtime.NumCPU() * 2, // 并发限制：CPU核心数的2倍
)
```

### 3. 监控和指标

```go
// 添加监控指标
type CircuitBreaker struct {
    // ... 其他字段
    metrics struct {
        totalRequests   int64
        failedRequests  int64
        rejectedRequests int64
    }
}
```

### 4. 优雅降级

```go
err := cb.Execute(func() error {
    return callExternalService()
})
if err != nil {
    // 使用备用方案
    return useCachedData()
}
```

### 5. 组合使用

```go
// 创建隔离的服务调用者
type ServiceCaller struct {
    bulkhead *Bulkhead
    cb       *CircuitBreaker
}

func (s *ServiceCaller) Call(task func() error) error {
    // 先通过隔离舱壁
    return s.bulkhead.Execute(func() (interface{}, error) {
        // 再通过熔断器
        err := s.cb.Execute(task)
        return nil, err
    })
}
```

## 🚀 性能优化建议

1. **减少锁竞争**:
   - 使用 atomic 操作记录简单计数
   - 只在状态变更时加锁

2. **合理的超时配置**:
   - 根据实际服务响应时间调整
   - 考虑添加自适应超时机制

3. **监控和告警**:
   - 记录熔断事件
   - 统计拒绝率和成功率
   - 设置告警阈值

4. **资源复用**:
   - 复用 goroutine 而非每次创建
   - 使用对象池减少 GC 压力

## 🔗 相关模式

- **Rate Limiter (限流器)**: 控制请求速率，与 Bulkheads 互补
- **Retry Pattern (重试模式)**: 与 Circuit Breaker 配合使用
- **Timeout Pattern (超时模式)**: 防止长时间阻塞
- **Fallback Pattern (降级模式)**: 提供备用方案

## 📚 参考资源

- [Release It! - Michael Nygard](https://pragprog.com/titles/mnee2/release-it-second-edition/)
- [Microsoft Azure - Circuit Breaker Pattern](https://docs.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker)
- [Netflix Hystrix](https://github.com/Netflix/Hystrix)
