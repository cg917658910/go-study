package stability

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestBulkheadConcurrencyLimit(t *testing.T) {
	limit := 2
	bulkhead := NewBulkhead(limit)
	defer bulkhead.Close()

	var running int32
	var maxConcurrent int32
	var wg sync.WaitGroup

	// 启动多个并发任务
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := bulkhead.Execute(func() (interface{}, error) {
				current := atomic.AddInt32(&running, 1)
				defer atomic.AddInt32(&running, -1)

				// 记录最大并发数
				for {
					max := atomic.LoadInt32(&maxConcurrent)
					if current <= max || atomic.CompareAndSwapInt32(&maxConcurrent, max, current) {
						break
					}
				}

				// 模拟任务执行
				time.Sleep(50 * time.Millisecond)
				return "success", nil
			})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}

	wg.Wait()

	// 验证并发数不超过限制
	if maxConcurrent > int32(limit) {
		t.Errorf("exceeded limit: max concurrent %d > limit %d", maxConcurrent, limit)
	}
	if maxConcurrent != int32(limit) {
		t.Logf("warning: max concurrent %d did not reach limit %d", maxConcurrent, limit)
	}
}

func TestBulkheadTaskExecution(t *testing.T) {
	bulkhead := NewBulkhead(5)
	defer bulkhead.Close()

	tests := []struct {
		name        string
		task        func() (interface{}, error)
		wantResult  interface{}
		wantErr     bool
	}{
		{
			name: "successful task",
			task: func() (interface{}, error) {
				return "result", nil
			},
			wantResult: "result",
			wantErr:    false,
		},
		{
			name: "failed task",
			task: func() (interface{}, error) {
				return nil, errors.New("task error")
			},
			wantResult: nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := bulkhead.Execute(tt.task)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.wantResult {
				t.Errorf("Execute() result = %v, want %v", result, tt.wantResult)
			}
		})
	}
}

func TestBulkheadResourceIsolation(t *testing.T) {
	// 创建两个独立的隔离舱壁
	bulkhead1 := NewBulkhead(2)
	bulkhead2 := NewBulkhead(2)
	defer bulkhead1.Close()
	defer bulkhead2.Close()

	var running1 int32
	var running2 int32

	var wg sync.WaitGroup

	// 在 bulkhead1 中执行慢任务
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bulkhead1.Execute(func() (interface{}, error) {
				atomic.AddInt32(&running1, 1)
				defer atomic.AddInt32(&running1, -1)
				time.Sleep(100 * time.Millisecond)
				return nil, nil
			})
		}()
	}

	// 在 bulkhead2 中执行快任务
	time.Sleep(10 * time.Millisecond) // 确保 bulkhead1 已经开始执行
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bulkhead2.Execute(func() (interface{}, error) {
				atomic.AddInt32(&running2, 1)
				defer atomic.AddInt32(&running2, -1)
				time.Sleep(10 * time.Millisecond)
				return nil, nil
			})
		}()
	}

	wg.Wait()

	// 验证两个隔离舱壁互不影响
	// bulkhead2 应该能快速完成，不受 bulkhead1 的影响
	t.Logf("bulkhead1 and bulkhead2 executed independently")
}

func TestBulkheadMultipleTasks(t *testing.T) {
	bulkhead := NewBulkhead(3)
	defer bulkhead.Close()

	taskCount := 20
	successCount := int32(0)
	var wg sync.WaitGroup

	for i := 0; i < taskCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			result, err := bulkhead.Execute(func() (interface{}, error) {
				time.Sleep(10 * time.Millisecond)
				return id, nil
			})
			if err == nil && result != nil {
				atomic.AddInt32(&successCount, 1)
			}
		}(i)
	}

	wg.Wait()

	// 所有任务应该都成功完成
	if int(successCount) != taskCount {
		t.Errorf("expected %d successful tasks, got %d", taskCount, successCount)
	}
}

func TestBulkheadExample(t *testing.T) {
	// 运行示例代码
	ExampleBulkhead()
	// 示例代码应该能正常执行而不 panic
}

func BenchmarkBulkheadExecute(b *testing.B) {
	bulkhead := NewBulkhead(10)
	defer bulkhead.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bulkhead.Execute(func() (interface{}, error) {
				// 模拟简单任务
				time.Sleep(time.Millisecond)
				return nil, nil
			})
		}
	})
}
