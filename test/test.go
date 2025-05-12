package test

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

func MainGo() {
	go func() {
		i := 0
		for {
			i++
			fmt.Println("Hello from goroutine! i=", i)
			time.Sleep(time.Second * 1)
		}
	}()

	i := 0

	for {
		i++
		fmt.Println("Hello from main! i=", i)
		time.Sleep(time.Second * 1)
		if i == 2 {
			//break
		}
	}
}

func TestRuntimeGoshced() {
	go func(s string) {
		for i := 0; i < 2; i++ {
			fmt.Println(s)
		}
	}("world")
	// 主协程
	for i := 0; i < 2; i++ {
		// 切一下，再次分配任务
		runtime.Gosched()
		fmt.Println("hello")
	}
}

func TestRuntimeGoexit() {
	runtime.GOMAXPROCS(2)
	go func() func() {
		defer fmt.Println("a.defer")
		f := func() {
			defer fmt.Println("b.defer")
			runtime.Goexit()
			defer fmt.Println("c.defer")
			fmt.Println("b")
		}
		fmt.Println("a")
		return f
	}()
	for {

	}
}

func TestTimer() {
	timer := time.NewTimer(time.Second * 1)
	//<-timer.C
	//<-timer.C
	i := 0
	go func() {
		for {
			i++
			fmt.Println(<-timer.C)
			if i == 3 {
				timer.Stop()
			}
		}
	}()
	for {
	}

}

func countDown1() {

	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		abort <- struct{}{}
	}()
	go func() {
		ticker := time.Tick(1 * time.Second)
		for countdown := 10; countdown > 0; countdown-- {
			fmt.Println(countdown)
			<-ticker
		}
	}()
	select {
	case <-abort:
		fmt.Println("Aborted!")
		return

	}
	fmt.Println("Boom!")
}

var (
	x      int64
	wg     sync.WaitGroup
	lock   sync.Mutex
	rwlock sync.RWMutex
)

func write() {
	lock.Lock() // 加互斥锁
	//rwlock.Lock() // 加写锁
	x = x + 1
	time.Sleep(10 * time.Millisecond) // 假设读操作耗时10毫秒
	//rwlock.Unlock()                   // 解写锁
	lock.Unlock() // 解互斥锁
	wg.Done()
}

func read() {
	lock.Lock() // 加互斥锁
	//rwlock.RLock()               // 加读锁
	time.Sleep(time.Millisecond) // 假设读操作耗时1毫秒
	//rwlock.RUnlock()             // 解读锁
	lock.Unlock() // 解互斥锁
	wg.Done()
}

func testLock() {
	start := time.Now()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go write()
	}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go read()
	}

	wg.Wait()
	end := time.Now()
	fmt.Println(end.Sub(start))
}

func mm() {
	old := 10
	fmt.Println("left move", old, "->", old<<2)
	fmt.Println("right move", old, "->", old>>2)
}
func mm2() {
	old := 10
	fmt.Println("Original value:", old)

	// Left shift
	leftShift := old << 2
	fmt.Println("Left shift (old << 2):", old, "->", leftShift)
	fmt.Printf("Binary representation: %04b -> %04b\n", old, leftShift)
	fmt.Println("Explanation: Left shift by 2 is equivalent to multiplying by 2^2 = 4")

	// Right shift
	rightShift := old >> 2
	fmt.Println("Right shift (old >> 2):", old, "->", rightShift)
	fmt.Printf("Binary representation: %04b -> %04b\n", old, rightShift)
	fmt.Println("Explanation: Right shift by 2 is equivalent to integer division by 2^2 = 4")

	// Additional examples
	fmt.Println("\nAdditional examples:")

	// Left shift by 1
	fmt.Printf("15 << 1 = %d (Binary: %04b -> %04b)\n", 15<<1, 15, 15<<1)

	// Right shift by 1
	fmt.Printf("15 >> 1 = %d (Binary: %04b -> %04b)\n", 15>>1, 15, 15>>1)
}

func test() {
	x := 5
	fmt.Printf("%d", x)
	fmt.Println("test")
}

func DoTest() {
	testConsistent()
}
