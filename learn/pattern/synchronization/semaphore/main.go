package main

// Semaphore Pattern 信号量模式 同步模式
// 信号量是一种用于控制对共享资源的访问的同步机制
// 它维护一个计数器 来表示可用资源的数量
// 适用于需要限制对资源的并发访问的场景

import (
	"sync"
	"time"
)

type Semaphore struct {
	capacity chan struct{}
}

func NewSemaphore(max int) *Semaphore {
	return &Semaphore{
		capacity: make(chan struct{}, max),
	}
}

func (s *Semaphore) Acquire() {
	s.capacity <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.capacity
}

func worker(id int, sem *Semaphore, wg *sync.WaitGroup) {
	defer wg.Done()
	sem.Acquire()
	println("Worker", id, "acquired semaphore")
	time.Sleep(1 * time.Second) // 模拟工作
	println("Worker", id, "releasing semaphore")
	sem.Release()
}

func main() {
	const maxConcurrentWorkers = 3
	sem := NewSemaphore(maxConcurrentWorkers)
	var wg sync.WaitGroup

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go worker(i, sem, &wg)
	}

	wg.Wait()
}

// 优点: 控制对共享资源的并发访问 防止资源争用 提高系统稳定性
// 缺点: 可能导致死锁 需要正确管理信号量的获取和释放
// 适用场景: 限制对数据库连接池的访问 控制并发请求数 限制对文件的访问等
// 对比互斥锁: 互斥锁用于保护临界区 只能允许一个线程访问 而信号量允许多个线程访问共享资源
// 对比读写锁: 读写锁允许多个读者或一个写者访问 而信号量通过计数器控制访问数量
// 举例: 限制同时处理的HTTP请求数 使用信号量可以确保不会超过设定的并发请求数
// Web领域例子: 限制对API的并发调用数 使用信号量可以防止服务器过载
// 数据库领域例子: 限制对数据库连接的并发访问 使用信号量可以确保不会超过连接池的最大连接数
// 函数编程领域例子: 限制对某个函数的并发调用数 使用信号量可以控制同时执行该函数的协程数量
