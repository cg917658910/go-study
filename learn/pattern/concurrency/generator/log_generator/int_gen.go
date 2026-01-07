package gen_ext

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// IntGenerator 生成大规模整数数据
func IntGenerator(done <-chan struct{}, count int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 0; i < count; i++ {
			select {
			case <-done:
				return
			case out <- rand.Intn(100): // 生成随机整数
			}
		}
	}()
	return out
}

// Worker 并发处理数据（如计算平方）
func Worker(done <-chan struct{}, in <-chan int, wg *sync.WaitGroup) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		defer wg.Done()
		for num := range in {
			select {
			case <-done:
				return
			case out <- num * num: // 计算平方
			}
		}
	}()
	return out
}

// Merge 合并多个通道的数据
func Merge(done <-chan struct{}, channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		defer wg.Done()
		for val := range c {
			select {
			case <-done:
				return
			case out <- val:
			}
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	rand.Seed(time.Now().UnixNano())
	done := make(chan struct{})
	defer close(done)

	// 生成大规模数据
	data := IntGenerator(done, 100)

	// 启动多个 Worker 并发处理数据
	var workers []<-chan int
	var wg sync.WaitGroup
	numWorkers := 4
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		workers = append(workers, Worker(done, data, &wg))
	}

	// 合并所有 Worker 的结果
	results := Merge(done, workers...)

	// 消费结果
	for result := range results {
		fmt.Println(result)
	}
}
