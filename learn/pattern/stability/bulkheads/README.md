# Bulkheads Pattern (隔离舱壁模式)

## 定义
隔离舱壁模式是一种稳定性设计模式，将系统划分为多个隔离的单元（舱壁），每个单元都有自己的资源池和限制，防止一个单元的故障影响到其他单元，类似船舶的防水舱设计。

## 目的
- 资源隔离，防止资源耗尽
- 故障隔离，防止级联失败
- 提高系统整体稳定性
- 限制故障影响范围

## 使用场景
- **微服务架构**: 为不同服务分配独立的资源池
- **多租户系统**: 隔离不同租户的资源
- **数据库连接池**: 为不同类型的查询分配独立连接池
- **线程池管理**: 为不同任务类型创建独立线程池
- **API 网关**: 为不同的后端服务设置独立的并发限制

## 工作原理

### 核心概念

隔离舱壁模式的核心是**资源分区**和**并发限制**：

1. **资源分区**: 将系统资源（goroutine、连接、内存等）划分为独立的池
2. **并发限制**: 每个池都有自己的并发限制
3. **故障隔离**: 一个池的资源耗尽不会影响其他池

### 船舶舱壁类比

```
┌─────────────────────────────────┐
│  ┌────┐  ┌────┐  ┌────┐  ┌────┐│
│  │舱1 │  │舱2 │  │舱3 │  │舱4 ││
│  │✓   │  │✓   │  │✗   │  │✓   ││ ← 舱3漏水，其他舱室不受影响
│  └────┘  └────┘  └────┘  └────┘│
└─────────────────────────────────┘
```

## Go 特有实现

使用 buffered channel 作为信号量控制并发：

```go
type Bulkhead struct {
    limit     int            // 并发限制
    requestCh chan func()    // 请求通道
    resultCh  chan interface{}
    errorCh   chan error
    quitCh    chan struct{} // 退出信号
}

func (b *Bulkhead) start() {
    sem := make(chan struct{}, b.limit)
    for {
        select {
        case req := <-b.requestCh:
            sem <- struct{}{}  // 获取信号量
            go func(r func()) {
                defer func() { <-sem }()  // 释放信号量
                r()
            }(req)
        case <-b.quitCh:
            return
        }
    }
}
```

## 优点
1. **故障隔离** - 限制单个组件故障的影响范围
2. **资源保护** - 防止某个组件耗尽所有资源
3. **提高稳定性** - 部分功能失败，系统整体仍可用
4. **公平性** - 保证不同组件都能获得资源
5. **可预测性** - 资源使用更加可控

## 缺点
1. **资源浪费** - 预留资源可能未充分利用
2. **配置复杂** - 需要合理规划资源分配
3. **增加复杂性** - 需要维护多个资源池
4. **可能限制吞吐量** - 过度隔离可能降低整体性能

## 资源配置建议

### 并发限制 (limit)
- 考虑因素: 
  - 系统总资源（CPU、内存、连接数）
  - 服务的重要性和优先级
  - 预期的流量模式

- 推荐策略:
  ```go
  // 基于 CPU 核心数
  cpuLimit := runtime.NumCPU() * 2
  
  // 基于内存
  memoryLimit := totalMemory / avgTaskMemory
  
  // 基于连接数
  connLimit := maxConnections / numServices
  ```

## 实际应用示例

### 1. 服务隔离

为不同的服务创建独立的隔离舱壁：

```go
type ServiceManager struct {
    userService    *Bulkhead
    orderService   *Bulkhead
    paymentService *Bulkhead
}

func NewServiceManager() *ServiceManager {
    return &ServiceManager{
        userService:    NewBulkhead(10),  // 用户服务：10个并发
        orderService:   NewBulkhead(20),  // 订单服务：20个并发
        paymentService: NewBulkhead(5),   // 支付服务：5个并发
    }
}
```

### 2. 数据库连接池隔离

```go
type DBManager struct {
    readPool  *Bulkhead  // 读操作池
    writePool *Bulkhead  // 写操作池
}

func (m *DBManager) Query(sql string) (interface{}, error) {
    return m.readPool.Execute(func() (interface{}, error) {
        return db.Query(sql)
    })
}

func (m *DBManager) Update(sql string) (interface{}, error) {
    return m.writePool.Execute(func() (interface{}, error) {
        return db.Exec(sql)
    })
}
```

### 3. 多租户隔离

```go
type TenantManager struct {
    tenants map[string]*Bulkhead
    mu      sync.RWMutex
}

func (m *TenantManager) Execute(tenantID string, task func() error) error {
    m.mu.RLock()
    bulkhead := m.tenants[tenantID]
    m.mu.RUnlock()
    
    if bulkhead == nil {
        // 为新租户创建隔离舱壁
        bulkhead = NewBulkhead(100)
        m.mu.Lock()
        m.tenants[tenantID] = bulkhead
        m.mu.Unlock()
    }
    
    _, err := bulkhead.Execute(func() (interface{}, error) {
        return nil, task()
    })
    return err
}
```

## 监控指标

建议记录以下指标：

1. **资源使用率**: 当前使用 / 总限制
2. **拒绝次数**: 达到限制被拒绝的请求数
3. **等待时间**: 请求等待获取资源的时间
4. **活跃任务数**: 当前正在执行的任务数
5. **任务完成时间**: 任务执行的平均时间

## 与其他模式的关系

### Bulkheads vs Circuit Breaker
- **Bulkheads**: 关注资源隔离和限制，防止资源耗尽
- **Circuit Breaker**: 关注故障检测和恢复，防止故障蔓延

### Bulkheads vs Rate Limiter
- **Bulkheads**: 限制并发数（同时执行的请求数）
- **Rate Limiter**: 限制请求速率（单位时间内的请求数）

### 组合使用

```go
type ProtectedService struct {
    bulkhead *Bulkhead        // 资源隔离
    cb       *CircuitBreaker  // 故障保护
    limiter  *RateLimiter     // 速率限制
}

func (s *ProtectedService) Call(task func() error) error {
    // 1. 先检查速率限制
    if !s.limiter.Allow() {
        return errors.New("rate limit exceeded")
    }
    
    // 2. 通过隔离舱壁
    _, err := s.bulkhead.Execute(func() (interface{}, error) {
        // 3. 通过熔断器
        return nil, s.cb.Execute(task)
    })
    return err
}
```

## 最佳实践

### 1. 合理分配资源

```go
// 根据服务优先级分配
highPriority := NewBulkhead(50)  // 高优先级服务
mediumPriority := NewBulkhead(30) // 中优先级服务
lowPriority := NewBulkhead(20)   // 低优先级服务
```

### 2. 添加超时控制

```go
func (b *Bulkhead) ExecuteWithTimeout(
    task func() (interface{}, error),
    timeout time.Duration,
) (interface{}, error) {
    resultCh := make(chan interface{})
    errorCh := make(chan error)
    
    go func() {
        result, err := b.Execute(task)
        if err != nil {
            errorCh <- err
        } else {
            resultCh <- result
        }
    }()
    
    select {
    case result := <-resultCh:
        return result, nil
    case err := <-errorCh:
        return nil, err
    case <-time.After(timeout):
        return nil, errors.New("timeout")
    }
}
```

### 3. 动态调整限制

```go
type DynamicBulkhead struct {
    *Bulkhead
    mu sync.RWMutex
}

func (b *DynamicBulkhead) UpdateLimit(newLimit int) {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.limit = newLimit
    // 重新创建信号量通道
}
```

### 4. 优雅关闭

```go
func (b *Bulkhead) Close() {
    close(b.quitCh)
    // 等待所有正在执行的任务完成
}
```

## 常见陷阱

### 1. Goroutine 泄露

❌ **错误做法**:
```go
go func() {
    // 忘记 defer 释放资源
    sem <- struct{}{}
    doWork()
}()
```

✅ **正确做法**:
```go
go func() {
    sem <- struct{}{}
    defer func() { <-sem }()
    doWork()
}()
```

### 2. 过度隔离

❌ **错误做法**: 为每个小功能都创建隔离舱壁
✅ **正确做法**: 根据业务重要性和资源需求合理分组

### 3. 死锁风险

❌ **错误做法**: 在持有锁时等待 channel
✅ **正确做法**: 避免嵌套的资源获取

## 测试建议

```go
func TestBulkheadLimit(t *testing.T) {
    limit := 2
    bulkhead := NewBulkhead(limit)
    defer bulkhead.Close()
    
    // 测试并发限制
    var running int32
    var wg sync.WaitGroup
    
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            bulkhead.Execute(func() (interface{}, error) {
                current := atomic.AddInt32(&running, 1)
                defer atomic.AddInt32(&running, -1)
                
                // 验证不超过限制
                if current > int32(limit) {
                    t.Errorf("exceeded limit: %d > %d", current, limit)
                }
                time.Sleep(10 * time.Millisecond)
                return nil, nil
            })
        }()
    }
    
    wg.Wait()
}
```

## 参考资源

- [Release It! - Michael Nygard](https://pragprog.com/titles/mnee2/release-it-second-edition/)
- [Microsoft Azure - Bulkhead Pattern](https://docs.microsoft.com/en-us/azure/architecture/patterns/bulkhead)
