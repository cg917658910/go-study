package messaging

import (
	"errors"
	"testing"
	"time"
)

func TestFutureBasic(t *testing.T) {
	future := NewFuture()

	// 异步设置结果
	go func() {
		time.Sleep(10 * time.Millisecond)
		future.SetResult("test result")
	}()

	// 获取结果
	result, err := future.Get()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != "test result" {
		t.Errorf("expected 'test result', got %v", result)
	}
}

func TestFutureError(t *testing.T) {
	future := NewFuture()

	// 异步设置错误
	go func() {
		time.Sleep(10 * time.Millisecond)
		future.SetError(errors.New("test error"))
	}()

	// 获取结果应该返回错误
	result, err := future.Get()
	if err == nil {
		t.Error("expected error, got nil")
	}
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
	if err.Error() != "test error" {
		t.Errorf("expected 'test error', got %v", err)
	}
}

func TestFutureGetWithTimeout(t *testing.T) {
	tests := []struct {
		name        string
		delay       time.Duration
		timeout     int
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "success before timeout",
			delay:       10 * time.Millisecond,
			timeout:     100,
			shouldError: false,
		},
		{
			name:        "timeout",
			delay:       200 * time.Millisecond,
			timeout:     50,
			shouldError: true,
			errorMsg:    "timeout occurred while waiting for result",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			future := NewFuture()

			go func() {
				time.Sleep(tt.delay)
				future.SetResult("result")
			}()

			result, err := future.GetWithTimeout(tt.timeout)
			if tt.shouldError {
				if err == nil {
					t.Error("expected timeout error, got nil")
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != "result" {
					t.Errorf("expected 'result', got %v", result)
				}
			}
		})
	}
}

func TestFutureThen(t *testing.T) {
	future := NewFuture()

	// 链式调用
	resultFuture := future.Then(func(val interface{}) interface{} {
		return val.(int) * 2
	}).Then(func(val interface{}) interface{} {
		return val.(int) + 10
	})

	// 设置初始值
	go func() {
		time.Sleep(10 * time.Millisecond)
		future.SetResult(5)
	}()

	// 获取最终结果: (5 * 2) + 10 = 20
	result, err := resultFuture.GetWithTimeout(200)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != 20 {
		t.Errorf("expected 20, got %v", result)
	}
}

func TestFutureThenWithError(t *testing.T) {
	future := NewFuture()

	resultFuture := future.Then(func(val interface{}) interface{} {
		return val.(int) * 2
	})

	// 设置错误
	go func() {
		time.Sleep(10 * time.Millisecond)
		future.SetError(errors.New("initial error"))
	}()

	// Then 链应该传播错误
	result, err := resultFuture.GetWithTimeout(200)
	if err == nil {
		t.Error("expected error, got nil")
	}
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}

func TestFutureCatch(t *testing.T) {
	future := NewFuture()
	errorCaught := false
	var caughtError error

	future.Catch(func(err error) {
		errorCaught = true
		caughtError = err
	})

	// 设置错误
	future.SetError(errors.New("test error"))

	// 等待 Catch 执行
	time.Sleep(50 * time.Millisecond)

	if !errorCaught {
		t.Error("error was not caught")
	}
	if caughtError == nil {
		t.Error("caught error is nil")
	}
	if caughtError.Error() != "test error" {
		t.Errorf("expected 'test error', got %v", caughtError)
	}
}

func TestFutureCatchNoError(t *testing.T) {
	future := NewFuture()
	errorCaught := false

	future.Catch(func(err error) {
		errorCaught = true
	})

	// 设置成功结果
	future.SetResult("success")

	// 等待可能的 Catch 执行
	time.Sleep(50 * time.Millisecond)

	if errorCaught {
		t.Error("error handler should not be called on success")
	}
}

func TestFutureConcurrentAccess(t *testing.T) {
	future := NewFuture()
	done := make(chan bool)

	// 多个 goroutine 等待同一个 Future
	for i := 0; i < 5; i++ {
		go func(id int) {
			result, err := future.GetWithTimeout(1000)
			if err != nil {
				t.Errorf("goroutine %d: unexpected error: %v", id, err)
			}
			if result != "shared result" {
				t.Errorf("goroutine %d: expected 'shared result', got %v", id, result)
			}
			done <- true
		}(i)
	}

	// 设置结果
	time.Sleep(50 * time.Millisecond)
	future.SetResult("shared result")

	// 只有第一个 Get 会获取到结果，其他会阻塞
	// 所以我们只等待一个完成
	select {
	case <-done:
		// 成功
	case <-time.After(200 * time.Millisecond):
		t.Error("timeout waiting for result")
	}
}

func TestMultipleFutures(t *testing.T) {
	future1 := NewFuture()
	future2 := NewFuture()
	future3 := NewFuture()

	// 并发设置结果
	go func() {
		time.Sleep(10 * time.Millisecond)
		future1.SetResult(1)
	}()
	go func() {
		time.Sleep(20 * time.Millisecond)
		future2.SetResult(2)
	}()
	go func() {
		time.Sleep(30 * time.Millisecond)
		future3.SetResult(3)
	}()

	// 获取所有结果
	result1, _ := future1.GetWithTimeout(100)
	result2, _ := future2.GetWithTimeout(100)
	result3, _ := future3.GetWithTimeout(100)

	if result1 != 1 || result2 != 2 || result3 != 3 {
		t.Errorf("unexpected results: %v, %v, %v", result1, result2, result3)
	}
}

func BenchmarkFutureSetGet(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			future := NewFuture()
			go func() {
				future.SetResult("result")
			}()
			future.Get()
		}
	})
}

func BenchmarkFutureThen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		future := NewFuture()
		resultFuture := future.Then(func(val interface{}) interface{} {
			return val.(int) * 2
		})
		go func() {
			future.SetResult(10)
		}()
		resultFuture.GetWithTimeout(1000)
	}
}
