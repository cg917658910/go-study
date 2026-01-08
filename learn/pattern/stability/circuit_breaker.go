package stability

import (
	"fmt"
	"sync"
	"time"
)

type CircuitBreaker struct {
	mu               sync.Mutex
	state            string
	failureCount     int
	successCount     int
	failureThreshold int
	successThreshold int
	openTimeout      time.Duration
	lastOpened       time.Time
}

func NewCircuitBreaker(failureThreshold, successThreshold int, openTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:            "CLOSED",
		failureThreshold: failureThreshold,
		successThreshold: successThreshold,
		openTimeout:      openTimeout,
	}
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
		err := task()
		if err != nil {
			cb.failureCount++
			if cb.failureCount >= cb.failureThreshold {
				cb.state = "OPEN"
				cb.lastOpened = time.Now()
			}
			return err
		}
		cb.successCount++
		if cb.successCount >= cb.successThreshold {
			cb.state = "CLOSED"
			cb.failureCount = 0
			cb.successCount = 0
		}
		return nil

	case "CLOSED":
		err := task()
		if err != nil {
			cb.failureCount++
			if cb.failureCount >= cb.failureThreshold {
				cb.state = "OPEN"
				cb.lastOpened = time.Now()
			}
			return err
		}
		cb.failureCount = 0
		return nil

	default:
		return fmt.Errorf("unknown circuit breaker state")
	}
	return nil
}

func (cb *CircuitBreaker) GetState() string {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// Example usage
func ExampleCircuitBreaker() {
	cb := NewCircuitBreaker(3, 2, 5*time.Second)

	tasks := []func() error{
		func() error {
			// Simulate a successful task
			return nil
		},
		func() error {
			// Simulate a failed task
			return fmt.Errorf("task failed")
		},
	}

	for i, task := range tasks {
		err := cb.Execute(task)
		if err != nil {
			fmt.Printf("Task %d failed: %v\n", i+1, err)
		} else {
			fmt.Printf("Task %d succeeded\n", i+1)
		}
		fmt.Printf("Circuit Breaker State: %s\n", cb.GetState())
	}
}

//总结: 熔断器模式是一种稳定性设计模式 通过监控系统的运行状态 在检测到故障时 自动切断对故障组件的调用 从而防止故障蔓延
// 优点: 提高系统的稳定性和可靠性 防止级联故障
// 缺点: 可能导致部分请求被拒绝 需要合理配置熔断器的参数
// 适用场景: 分布式系统 微服务架构 需要防止某个服务或组件的故障影响到整个系统
// 对比隔离舱壁模式: 隔离舱壁模式关注资源隔离和限制 而熔断器模式关注故障检测和恢复
// 对比限流模式: 限流模式关注请求速率控制 而熔断器模式关注故障检测和恢复
// 举例: 在微服务架构中 为每个服务实例设置熔断器 防止某个实例的故障影响到整个服务
// Web领域例子: 为不同的Web应用设置熔断器 防止某个应用的故障影响到其他应用
// 数据库领域例子: 为不同的数据库连接设置熔断器 防止某个连接的故障影响到其他连接
