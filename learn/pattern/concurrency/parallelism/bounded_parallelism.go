package parallelism

// Bounded Parallelism Pattern 有界并行模式 并发设计模式
// 限制同时运行的协程数量 避免资源耗尽或过载
// 适用于需要控制并发度的场景 如数据库连接池 限流等

import (
	"fmt"
	"sync"
	"time"
)

// BoundedWorker 函数类型 定义有界工作函数签名
type BoundedWorker func(id int, wg *sync.WaitGroup)

// SampleBoundedWorker 实现一个示例有界工作函数
func SampleBoundedWorker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Bounded Worker %d starting\n", id)
	time.Sleep(time.Second) // 模拟工作耗时
	fmt.Printf("Bounded Worker %d done\n", id)
}

func main2() {
	var wg sync.WaitGroup
	numWorkers := 10
	maxConcurrent := 3
	sem := make(chan struct{}, maxConcurrent)

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		sem <- struct{}{} // 获取信号量
		go func(id int) {
			defer func() { <-sem }() // 释放信号量
			SampleBoundedWorker(id, &wg)
		}(i)
	}

	wg.Wait()
}

// 优点: 控制并发度 防止资源耗尽 提高系统稳定性
// 缺点: 增加代码复杂度 需要管理信号量的获取和释放
// 适用场景: 限制数据库连接池访问 控制并发请求数 限制文件访问等
// 对比信号量模式: 信号量模式提供了更通用的同步机制 而有界并行模式专注于限制协程数量
// 对比工作池模式: 工作池模式通过预创建固定数量的工作协程来处理任务 而有界并行模式动态控制协程数量
// 举例: 限制同时处理的HTTP请求数 使用有界并行模式可以确保不会超过设定的并发请求数
// Web领域例子: 限制对API的并发调用数 使用有界并行模式可以防止服务器过载
// 数据库领域例子: 限制对数据库连接的并发访问 使用有界并行模式可以确保不会超过连接池的最大连接数
