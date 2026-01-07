package synchronization

import (
	"sync"
	"testing"
)

func TestMonitor(t *testing.T) {
	monitor := NewMonitor()
	var wg sync.WaitGroup

	// 测试存款和取款操作
	wg.Add(2)

	go func() {
		defer wg.Done()
		monitor.Deposit(100) // 存款 100
	}()

	go func() {
		defer wg.Done()
		monitor.Withdraw(50) // 取款 50
	}()

	wg.Wait()

	// 验证最终余额是否正确
	expectedBalance := 50
	if monitor.balance != expectedBalance {
		t.Errorf("expected balance %d, got %d", expectedBalance, monitor.balance)
	}
}

func TestMonitor_WithdrawWait(t *testing.T) {
	monitor := NewMonitor()
	var wg sync.WaitGroup

	// 测试取款等待逻辑
	wg.Add(2)

	go func() {
		defer wg.Done()
		monitor.Withdraw(50) // 取款 50，应该等待存款
	}()

	go func() {
		defer wg.Done()
		monitor.Deposit(100) // 存款 100
	}()

	wg.Wait()

	// 验证最终余额是否正确
	expectedBalance := 50
	if monitor.balance != expectedBalance {
		t.Errorf("expected balance %d, got %d", expectedBalance, monitor.balance)
	}
}
