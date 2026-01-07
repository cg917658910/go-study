package parallelism

// Parallelism Pattern 并行模式 并发设计模式
// 通过多个协程同时执行任务 提高程序的执行效率
// 适用于计算密集型任务或需要同时处理多个任务的场景

import (
	"fmt"
	"sync"
	"time"
)

// Worker 函数类型 定义工作函数签名
type Worker func(id int, wg *sync.WaitGroup)

// SampleWorker 实现一个示例工作函数
func SampleWorker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Second) // 模拟工作耗时
	fmt.Printf("Worker %d done\n", id)
}

func main() {
	var wg sync.WaitGroup
	numWorkers := 5

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go SampleWorker(i, &wg)
	}

	wg.Wait()
}

// 优点: 提高程序执行效率 充分利用多核CPU资源
// 缺点: 增加代码复杂度 需要处理协程间的同步和通信问题
// 适用场景: 计算密集型任务 大规模数据处理 并发请求处理等
// 对比并发模式: 并发模式关注任务的调度和管理 而并行模式关注任务的同时执行
// 对比协程模式: 协程模式关注轻量级线程的使用 而并行模式关注多核CPU的利用
// 举例: 图像处理任务 可以将图像分割成多个部分 并行处理每个部分 提高处理速度
// Web领域例子: 并行处理多个HTTP请求 提高服务器响应速度
// 数据库领域例子: 并行执行多个查询任务 提高查询效率
