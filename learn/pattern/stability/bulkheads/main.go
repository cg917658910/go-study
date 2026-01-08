package stability

import "fmt"

//Bulkheads Pattern 隔离舱壁模式 稳定性设计模式
//将系统划分为多个隔离的单元 每个单元都有自己的资源和限制
//防止一个单元的故障影响到其他单元 提高系统的整体稳定性和可靠性
//适用于需要隔离不同服务或组件以防止级联故障的场景

type Bulkhead struct {
	limit     int
	current   int
	requestCh chan func()
	resultCh  chan interface{}
	errorCh   chan error
	quitCh    chan struct{}
}

func NewBulkhead(limit int) *Bulkhead {
	b := &Bulkhead{
		limit:     limit,
		requestCh: make(chan func()),
		resultCh:  make(chan interface{}),
		errorCh:   make(chan error),
		quitCh:    make(chan struct{}),
	}
	go b.start()
	return b
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

func (b *Bulkhead) Execute(task func() (interface{}, error)) (interface{}, error) {
	resultCh := make(chan interface{})
	errorCh := make(chan error)

	b.requestCh <- func() {
		result, err := task()
		if err != nil {
			errorCh <- err
			return
		}
		resultCh <- result
	}

	select {
	case res := <-resultCh:
		return res, nil
	case err := <-errorCh:
		return nil, err
	}
}

func (b *Bulkhead) Close() {
	close(b.quitCh)
}

func ExampleBulkhead() {
	bulkhead := NewBulkhead(2)
	defer bulkhead.Close()

	tasks := []func() (interface{}, error){
		func() (interface{}, error) {
			// 模拟任务1
			return "Task 1 completed", nil
		},
		func() (interface{}, error) {
			// 模拟任务2
			return "Task 2 completed", nil
		},
		func() (interface{}, error) {
			// 模拟任务3
			return nil, fmt.Errorf("Task 3 failed")
		},
	}

	for _, task := range tasks {
		result, err := bulkhead.Execute(task)
		if err != nil {
			fmt.Println("Error occurred:", err)
		} else {
			fmt.Println("Result:", result)
		}
	}
}

//总结: 隔离舱壁模式是一种稳定性设计模式 通过将系统划分为多个隔离的单元 来防止一个单元的故障影响

// 优点: 提高系统的稳定性和可靠性 防止级联故障
// 缺点: 增加了系统的复杂性 需要合理配置隔离单元的资源和限制
// 适用场景: 分布式系统 微服务架构 需要隔离不同服务或组件以防止级联故障
// 对比熔断器模式: 熔断器模式关注故障检测和恢复 而隔离舱壁模式关注资源隔离和限制
// 对比限流模式: 限流模式关注请求速率控制 而隔离舱壁模式关注资源隔离和限制
// 举例: 在微服务架构中 为每个服务实例设置隔离舱壁 防止某个实例的故障影响到整个服务
// Web领域例子: 为不同的Web应用设置隔离舱壁 防止某个应用的故障影响到其他应用
// 数据库领域例子: 为不同的数据库连接设置隔离舱壁 防止某个连接的故障影响到其他连接
