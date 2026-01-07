package concurrency

// Producer-Consumer Pattern 生产者-消费者模式 并发设计模式
// 通过一个缓冲区连接生产者和消费者 生产者将数据放入缓冲区 消费者从缓冲区取出数据
// 适用于生产和消费速度不匹配的场景

import (
	"sync"
)

type ProducerConsumer struct {
	buffer       chan interface{}
	wg           sync.WaitGroup
	numProducers int
	numConsumers int
}

func NewProducerConsumer(bufferSize, numProducers, numConsumers int) *ProducerConsumer {
	return &ProducerConsumer{
		buffer:       make(chan interface{}, bufferSize),
		numProducers: numProducers,
		numConsumers: numConsumers,
	}
}

func (pc *ProducerConsumer) Start(produceFunc func() interface{}, consumeFunc func(interface{})) {
	pc.wg.Add(pc.numProducers + pc.numConsumers)

	for i := 0; i < pc.numProducers; i++ {
		go func() {
			defer pc.wg.Done()
			for {
				item := produceFunc()
				if item == nil {
					break
				}
				pc.buffer <- item
			}
		}()
	}

	for i := 0; i < pc.numConsumers; i++ {
		go func() {
			defer pc.wg.Done()
			for item := range pc.buffer {
				consumeFunc(item)
			}
		}()
	}
}

func (pc *ProducerConsumer) Stop() {
	close(pc.buffer)
	pc.wg.Wait()
}

// 优点: 解耦生产者和消费者 提高系统吞吐量 适应不同速度的生产和消费
// 缺点: 需要处理缓冲区大小和阻塞问题 增加系统复杂度
// 适用场景: 数据处理管道 任务调度 系统负载均衡等
// 对比发布-订阅模式: 发布-订阅模式关注消息的分发和订阅 而生产者-消费者模式关注数据的生成和消费
// 对比工作池模式: 工作池模式通过固定数量的工作协程处理任务 而生产者-消费者模式通过缓冲区连接生产者和消费者
// 举例: 日志处理系统 生产者生成日志条目 消费者处理和存储日志
// Web领域例子: 请求处理系统 生产者接收请求 消费者处理请求并返回响应
// 数据库领域例子: 数据导入系统 生产者读取数据源 消费者将数据写入数据库
