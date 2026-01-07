package synchronization

// Monitor Pattern 监视器模式 并发设计模式
// 将共享资源的访问封装在一个对象中 通过互斥锁和条件变量来控制对共享资源的访问
// 适用于多个协程需要访问共享资源且需要同步访问的场景

import (
	"sync"
)

type Monitor struct {
	mu      sync.Mutex
	cond    *sync.Cond
	balance int
}

func NewMonitor() *Monitor {
	m := &Monitor{}
	m.cond = sync.NewCond(&m.mu)
	return m
}

func (m *Monitor) Deposit(amount int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.balance += amount
	m.cond.Signal()
}

func (m *Monitor) Withdraw(amount int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for m.balance < amount {
		m.cond.Wait()
	}
	m.balance -= amount
}

// 优点: 提供了对共享资源的安全访问机制 避免数据竞争和不一致性
// 缺点: 增加了代码复杂度 可能导致死锁和性能瓶颈
// 适用场景: 多个协程需要访问共享资源且需要同步访问的场景 如银行账户管理 生产者-消费者问题等
// 对比互斥锁: 互斥锁只能保护临界区 而监视器模式封装了共享资源的访问逻辑
// 对比读写锁: 读写锁允许多个读者或一个写者访问 而监视器模式通过条件变量实现更复杂的同步逻辑
// 举例: 银行账户管理 使用监视器模式可以确保多个协程对账户余额的安全访问和更新
// Web领域例子: 会话管理 使用监视器模式可以确保多个请求对用户会话数据的安全访问
// 数据库领域例子: 连接池管理 使用监视器模式可以确保多个协程对数据库连接池的安全访问和管理
