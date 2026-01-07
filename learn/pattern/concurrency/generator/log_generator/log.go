package gen_ext

import (
	"fmt"
	"strings"
	"time"
)

// LogEntry 表示一条日志记录
type LogEntry struct {
	Timestamp time.Time
	Message   string
	Level     string
}

// LogGenerator 模拟日志生成器
func LogGenerator(done <-chan struct{}) <-chan LogEntry {
	out := make(chan LogEntry)
	go func() {
		defer close(out)
		levels := []string{"INFO", "WARN", "ERROR"}
		for i := 0; ; i++ {
			select {
			case <-done:
				return
			case out <- LogEntry{
				Timestamp: time.Now(),
				Message:   fmt.Sprintf("Log message %d", i),
				Level:     levels[i%len(levels)],
			}:
				time.Sleep(500 * time.Millisecond) // 模拟日志生成间隔
			}
		}
	}()
	return out
}

// FilterLogs 过滤日志，只保留指定级别的日志
func FilterLogs(done <-chan struct{}, in <-chan LogEntry, level string) <-chan LogEntry {
	out := make(chan LogEntry)
	go func() {
		defer close(out)
		for log := range in {
			if strings.EqualFold(log.Level, level) {
				select {
				case <-done:
					return
				case out <- log:
				}
			}
		}
	}()
	return out
}

// FormatLogs 格式化日志为字符串
func FormatLogs(done <-chan struct{}, in <-chan LogEntry) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for log := range in {
			formatted := fmt.Sprintf("[%s] %s: %s", log.Timestamp.Format(time.RFC3339), log.Level, log.Message)
			select {
			case <-done:
				return
			case out <- formatted:
			}
		}
	}()
	return out
}

func runLog() {
	done := make(chan struct{})
	defer close(done)

	// 生成日志数据
	logs := LogGenerator(done)

	// 过滤日志，只保留 ERROR 级别
	filteredLogs := FilterLogs(done, logs, "ERROR")

	// 格式化日志
	formattedLogs := FormatLogs(done, filteredLogs)

	// 消费日志
	for i := 0; i < 5; i++ {
		fmt.Println(<-formattedLogs)
	}
}

// 优点: 实现数据的惰性生成和处理 提高内存效率 适用于流式数据处理
// 缺点: 需要处理协程和通道的同步问题 增加代码复杂度
// 适用场景: 处理大量日志数据 流式数据处理 数据管道等
// 对比迭代器模式: 迭代器模式通过对象封装遍历逻辑 而生成器模式通过协程和通道实现数据生成和处理
// 对比发布-订阅模式: 发布-订阅模式关注消息的分发和订阅 而生成器模式关注数据的生成和消费
// 举例: 实时日志处理 使用生成器模式可以按需生成和处理日志条目 避免一次性加载大量日志数据
// Web领域例子: 处理实时用户活动流 使用生成器模式可以按需生成用户活动事件 并通过通道传递给消费者
// 数据库领域例子: 处理大规模查询结果 使用生成器模式可以按需生成查询结果 避免一次性加载大量数据到内存
