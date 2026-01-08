# Circuit Breaker Pattern (熔断器模式)

## 定义
熔断器模式是一种稳定性设计模式，通过监控系统的运行状态，在检测到故障时自动切断对故障组件的调用，从而防止故障蔓延，并在系统恢复后自动恢复调用。

## 目的
- 防止级联故障
- 快速失败，避免资源浪费
- 自动恢复，提高系统可用性
- 保护下游服务

## 使用场景
- **微服务架构**: 保护服务间调用
- **外部 API 调用**: 防止外部服务故障影响系统
- **数据库连接**: 保护数据库连接池
- **分布式系统**: 防止部分节点故障导致系统崩溃

## 工作原理

### 三种状态

1. **CLOSED (关闭状态)**:
   - 正常状态，所有请求正常通过
   - 记录失败次数
   - 失败次数达到阈值 → 转换到 OPEN 状态

2. **OPEN (打开/熔断状态)**:
   - 熔断状态，快速失败，直接拒绝请求
   - 不调用下游服务，避免资源浪费
   - 等待超时时间 → 转换到 HALF-OPEN 状态

3. **HALF-OPEN (半开状态)**:
   - 测试状态，允许少量请求通过
   - 如果成功次数达到阈值 → 转换到 CLOSED 状态
   - 如果失败 → 转换回 OPEN 状态

### 状态转换图

```
     失败次数 >= 阈值
CLOSED ──────────────→ OPEN
  ↑                      │
  │                      │ 超时时间到
  │                      ↓
  │                  HALF-OPEN
  │                      │
  └──────────────────────┘
    成功次数 >= 阈值
```

## Go 特有实现

使用 mutex 保护状态，确保线程安全：

```go
type CircuitBreaker struct {
    mu               sync.Mutex
    state            string        // 当前状态
    failureCount     int           // 失败计数
    successCount     int           // 成功计数
    failureThreshold int           // 失败阈值
    successThreshold int           // 成功阈值
    openTimeout      time.Duration // 打开状态超时
    lastOpened       time.Time     // 上次打开时间
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
    // ... 其他状态处理
    }
}
```

## 优点
1. **防止故障蔓延** - 快速隔离故障服务
2. **快速失败** - 避免等待超时，提升响应速度
3. **自动恢复** - 系统恢复后自动恢复调用
4. **保护资源** - 避免浪费线程和连接等资源
5. **提高可用性** - 部分功能降级，整体系统仍可用

## 缺点
1. **增加复杂性** - 需要配置和维护状态机
2. **误判风险** - 阈值设置不当可能导致误熔断
3. **部分请求被拒** - OPEN 状态下的正常请求也会被拒绝
4. **需要监控** - 需要监控熔断状态和指标

## 参数配置建议

### 失败阈值 (failureThreshold)
- 推荐值: 3-10
- 太小: 容易误判，正常波动就触发熔断
- 太大: 反应迟钝，故障已蔓延才熔断

### 成功阈值 (successThreshold)
- 推荐值: 2-5
- 确保服务确实恢复再转为 CLOSED

### 超时时间 (openTimeout)
- 推荐值: 5-60 秒
- 取决于下游服务的恢复时间
- 太短: 服务未恢复就尝试，增加负载
- 太长: 影响用户体验

## 实际应用示例

### 1. HTTP 客户端保护

```go
type HTTPClient struct {
    client *http.Client
    cb     *CircuitBreaker
}

func (c *HTTPClient) Get(url string) (*http.Response, error) {
    var resp *http.Response
    err := c.cb.Execute(func() error {
        var err error
        resp, err = c.client.Get(url)
        return err
    })
    return resp, err
}
```

### 2. 数据库调用保护

```go
type DBClient struct {
    db *sql.DB
    cb *CircuitBreaker
}

func (c *DBClient) Query(query string) ([]Row, error) {
    var rows []Row
    err := c.cb.Execute(func() error {
        var err error
        rows, err = c.db.Query(query)
        return err
    })
    return rows, err
}
```

### 3. 带降级的服务调用

```go
func GetUserInfo(userID string) (*User, error) {
    var user *User
    err := cb.Execute(func() error {
        var err error
        user, err = externalService.GetUser(userID)
        return err
    })
    
    if err != nil {
        // 降级：从缓存获取
        return getCachedUser(userID)
    }
    return user, nil
}
```

## 监控指标

建议记录以下指标：

1. **状态变更事件**: 记录每次状态转换
2. **请求总数**: 总请求量
3. **成功请求数**: 成功的请求数
4. **失败请求数**: 失败的请求数
5. **拒绝请求数**: 在 OPEN 状态被拒绝的请求数
6. **熔断时长**: 每次熔断持续的时间

## 与其他模式的关系

- **Retry Pattern**: 熔断器可以与重试模式结合，失败时先重试，多次失败后熔断
- **Timeout Pattern**: 配合超时模式，防止长时间等待
- **Bulkheads Pattern**: 隔离舱壁模式隔离资源，熔断器保护调用
- **Fallback Pattern**: 熔断后提供降级方案

## 最佳实践

1. **合理配置阈值**: 根据实际流量和错误率调整
2. **添加监控**: 实时监控熔断器状态
3. **提供降级方案**: 熔断时提供备用数据或功能
4. **记录日志**: 记录熔断事件，便于问题排查
5. **定期审查**: 定期审查配置是否合理
6. **测试熔断逻辑**: 在测试环境验证熔断器行为

## 反模式

❌ **不要做**:
- 不要对所有请求使用同一个熔断器实例
- 不要忽略 HALF-OPEN 状态
- 不要使用过短的超时时间
- 不要在熔断时静默失败（应该告知用户）

✅ **应该做**:
- 为不同的下游服务使用独立的熔断器
- 正确实现三种状态的转换
- 根据服务特性配置超时
- 提供清晰的错误信息和降级方案
