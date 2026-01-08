package messaging

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// Helper function to generate test message content
func makeTestMessage(i int) string {
	return fmt.Sprintf("Message %d", i)
}

func TestPushPullBasic(t *testing.T) {
	producer := NewProducer(10)
	consumer := NewConsumer(producer.messages)

	// 生产者推送消息
	go func() {
		for i := 0; i < 5; i++ {
			msg := Message{Content: makeTestMessage(i)}
			producer.Push(msg)
		}
		close(producer.messages)
	}()

	// 消费者拉取消息
	count := 0
	for {
		msg, ok := <-consumer.messages
		if !ok {
			break
		}
		if msg.Content == "" {
			t.Error("received empty message")
		}
		count++
	}

	if count != 5 {
		t.Errorf("expected 5 messages, got %d", count)
	}
}

func TestProducerPush(t *testing.T) {
	producer := NewProducer(5)
	defer close(producer.messages)

	messages := []Message{
		{Content: "Message 1"},
		{Content: "Message 2"},
		{Content: "Message 3"},
	}

	// 推送消息
	for _, msg := range messages {
		producer.Push(msg)
	}

	// 验证消息在通道中
	for i, expected := range messages {
		select {
		case msg := <-producer.messages:
			if msg.Content != expected.Content {
				t.Errorf("message %d: expected '%s', got '%s'", i, expected.Content, msg.Content)
			}
		case <-time.After(100 * time.Millisecond):
			t.Errorf("timeout waiting for message %d", i)
		}
	}
}

func TestConsumerPull(t *testing.T) {
	messages := make(chan Message, 5)
	consumer := NewConsumer(messages)

	// 预先放入消息
	testMessages := []Message{
		{Content: "Message 1"},
		{Content: "Message 2"},
		{Content: "Message 3"},
	}

	for _, msg := range testMessages {
		messages <- msg
	}

	// 消费者拉取
	for i, expected := range testMessages {
		msg := consumer.Pull()
		if msg.Content != expected.Content {
			t.Errorf("message %d: expected '%s', got '%s'", i, expected.Content, msg.Content)
		}
	}
}

func TestProducerConsumerConcurrent(t *testing.T) {
	producer := NewProducer(100)
	consumer := NewConsumer(producer.messages)

	messageCount := 100
	received := make(map[string]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 生产者 goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < messageCount; i++ {
			msg := Message{Content: makeTestMessage(i)}
			producer.Push(msg)
		}
		close(producer.messages)
	}()

	// 消费者 goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range consumer.messages {
			mu.Lock()
			received[msg.Content] = true
			mu.Unlock()
		}
	}()

	wg.Wait()

	// 验证接收了消息
	if len(received) == 0 {
		t.Error("no messages received")
	}
}

func TestMultipleProducersSingleConsumer(t *testing.T) {
	messages := make(chan Message, 100)
	defer close(messages)

	numProducers := 3
	messagesPerProducer := 10
	totalMessages := numProducers * messagesPerProducer

	var wg sync.WaitGroup

	// 启动多个生产者
	for p := 0; p < numProducers; p++ {
		wg.Add(1)
		go func(producerID int) {
			defer wg.Done()
			producer := &Producer{messages: messages}
			for i := 0; i < messagesPerProducer; i++ {
				msg := Message{Content: fmt.Sprintf("Message from producer %d", producerID)}
				producer.Push(msg)
			}
		}(p)
	}

	// 消费者计数
	receivedCount := 0
	done := make(chan bool)
	go func() {
		consumer := NewConsumer(messages)
		for i := 0; i < totalMessages; i++ {
			consumer.Pull()
			receivedCount++
		}
		done <- true
	}()

	wg.Wait()
	<-done

	if receivedCount != totalMessages {
		t.Errorf("expected %d messages, got %d", totalMessages, receivedCount)
	}
}

func TestSingleProducerMultipleConsumers(t *testing.T) {
	messages := make(chan Message, 100)
	producer := NewProducer(100)
	producer.messages = messages

	numConsumers := 3
	totalMessages := 30
	receivedCount := int32(0)
	var mu sync.Mutex

	var wg sync.WaitGroup

	// 启动多个消费者
	for c := 0; c < numConsumers; c++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			consumer := NewConsumer(messages)
			for {
				select {
				case msg, ok := <-consumer.messages:
					if !ok {
						return
					}
					if msg.Content != "" {
						mu.Lock()
						receivedCount++
						mu.Unlock()
					}
				case <-time.After(200 * time.Millisecond):
					return
				}
			}
		}()
	}

	// 生产者发送消息
	for i := 0; i < totalMessages; i++ {
		producer.Push(Message{Content: makeTestMessage(i)})
	}

	time.Sleep(100 * time.Millisecond)
	close(messages)
	wg.Wait()

	if receivedCount != int32(totalMessages) {
		t.Errorf("expected %d messages received, got %d", totalMessages, receivedCount)
	}
}

func TestBufferOverflow(t *testing.T) {
	bufferSize := 5
	producer := NewProducer(bufferSize)

	// 填满缓冲区
	for i := 0; i < bufferSize; i++ {
		producer.Push(Message{Content: makeTestMessage(i)})
	}

	// 尝试推送更多消息（应该阻塞或失败）
	done := make(chan bool)
	go func() {
		producer.Push(Message{Content: "Overflow"})
		done <- true
	}()

	// 等待一小段时间，不应该完成
	select {
	case <-done:
		t.Error("push should have blocked")
	case <-time.After(50 * time.Millisecond):
		// 正确，仍然阻塞
	}

	// 消费一条消息腾出空间
	<-producer.messages

	// 现在推送应该成功
	select {
	case <-done:
		// 成功
	case <-time.After(100 * time.Millisecond):
		t.Error("push should have succeeded after making space")
	}
}

func TestPushPullExample(t *testing.T) {
	// 运行示例代码，确保不 panic
	ExamplePushPull()
}

func BenchmarkProducerPush(b *testing.B) {
	producer := NewProducer(1000)
	msg := Message{Content: "Benchmark message"}

	go func() {
		// 消费者，防止阻塞
		for range producer.messages {
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		producer.Push(msg)
	}
}

func BenchmarkConsumerPull(b *testing.B) {
	messages := make(chan Message, b.N)
	consumer := NewConsumer(messages)

	// 预填充消息
	msg := Message{Content: "Benchmark message"}
	for i := 0; i < b.N; i++ {
		messages <- msg
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		consumer.Pull()
	}
}

func BenchmarkProducerConsumerParallel(b *testing.B) {
	producer := NewProducer(1000)
	consumer := NewConsumer(producer.messages)
	msg := Message{Content: "Benchmark message"}

	go func() {
		for range consumer.messages {
		}
	}()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			producer.Push(msg)
		}
	})
}
