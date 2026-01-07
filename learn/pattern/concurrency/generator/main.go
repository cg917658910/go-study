package main

// Generator Pattern 生成器模式 并发设计模式
// 通过协程和通道实现数据的惰性生成和消费
// 适用于需要处理大量数据或流式数据的场景

import (
	"fmt"
	"time"
)

// Generator 函数类型 定义生成器函数签名
type Generator func(done <-chan struct{}) <-chan int

// IntGenerator 实现一个整数生成器
func IntGenerator(done <-chan struct{}) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 0; ; i++ {
			select {
			case <-done:
				return
			case out <- i:
			}
		}
	}()
	return out
}

func main() {
	done := make(chan struct{})
	defer close(done)

	gen := IntGenerator(done)

	for i := 0; i < 5; i++ {
		fmt.Println(<-gen)
		time.Sleep(500 * time.Millisecond)
	}
}

// 优点: 实现数据的惰性生成和消费 提高内存效率 适用于流式数据处理
// 缺点: 需要处理协程和通道的同步问题 增加代码复杂度
// 适用场景: 处理大量数据 流式数据处理 数据管道等
// 对比迭代器模式: 迭代器模式通过对象封装遍历逻辑 而生成器模式通过协程和通道实现数据生成
// 对比发布-订阅模式: 发布-订阅模式关注消息的分发和订阅 而生成器模式关注数据的生成和消费
// 举例: 生成一个无限整数序列 使用生成器模式可以按需生成整数 避免一次性加载大量数据
// Web领域例子: 处理实时日志流 使用生成器模式可以按需生成日志条目 并通过通道传递给消费者
// 数据库领域例子: 处理大规模查询结果 使用生成器模式可以按需生成查询结果 避免一次性加载大量数据到内存
