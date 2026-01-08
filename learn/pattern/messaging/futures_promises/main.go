package messaging

import (
	"fmt"
	"time"
)

// Future 表示异步操作的结果
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

func (f *Future) SetResult(result interface{}) {
	f.result <- result
}

func (f *Future) SetError(err error) {
	f.err <- err
}

func (f *Future) Get() (interface{}, error) {
	select {
	case res := <-f.result:
		return res, nil
	case err := <-f.err:
		return nil, err
	}
}

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

func (f *Future) Catch(fn func(error)) *Future {
	go func() {
		_, err := f.Get()
		if err != nil {
			fn(err)
		}
	}()
	return f
}

func useFutureExample() {
	future := NewFuture()

	// 异步任务
	go func() {
		time.Sleep(2 * time.Second) // 模拟任务处理
		future.SetResult("Task Completed")
	}()

	// 主线程继续执行其他任务...

	// 获取结果
	result, err := future.GetWithTimeout(3000) // 设置超时时间为 3000 毫秒
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result.(string))
	}

	// 链式调用示例
	NewFuture().
		Then(func(res interface{}) interface{} {
			fmt.Println("Processing result:", res)
			return "Processed Result"
		}).
		Catch(func(err error) {
			fmt.Println("Error occurred:", err)
		})
}
