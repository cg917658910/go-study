package stability

import (
	"fmt"
	"testing"
	"time"
)

func TestCircuitBreaker(t *testing.T) {
	cb := NewCircuitBreaker(3, 2, 2*time.Second)

	tasks := []func() error{
		func() error {
			return fmt.Errorf("Task 1 failed")
		},
		func() error {
			return fmt.Errorf("Task 2 failed")
		},
		func() error {
			return fmt.Errorf("Task 3 failed")
		},
		func() error {
			return nil
		},
		func() error {
			return nil
		},
	}

	for i, task := range tasks {
		err := cb.Execute(task)
		if err != nil {
			fmt.Printf("Task %d: Error occurred: %v\n", i+1, err)
		} else {
			fmt.Printf("Task %d: Task succeeded\n", i+1)
		}
		fmt.Printf("Circuit Breaker State: %s\n", cb.GetState())
	}

	// 等待断路器从 OPEN 状态恢复到 HALF-OPEN
	time.Sleep(3 * time.Second)

	// 测试恢复逻辑
	for i := 0; i < 3; i++ {
		taskIndex := i // 捕获循环变量
		err := cb.Execute(func() error {
			// 第一个恢复任务失败，测试 HALF-OPEN 状态下的失败处理
			if taskIndex == 0 {
				return fmt.Errorf("Task failed in HALF-OPEN")
			}
			return nil
		})
		if err != nil {
			fmt.Printf("Recovery Task %d: Error occurred: %v\n", i+1, err)
		} else {
			fmt.Printf("Recovery Task %d: Task succeeded\n", i+1)
		}
		fmt.Printf("Circuit Breaker State: %s\n", cb.GetState())
	}
}
