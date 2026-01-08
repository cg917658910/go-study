package messaging

import "fmt"

// Push & Pull 模式
// Push模式: 生产者主动将消息发送给消费者
// Pull模式: 消费者主动从生产者获取消息
// 适用于需要解耦生产者和消费者的场景

type Message struct {
	Content string
}

type Producer struct {
	messages chan Message
}

func NewProducer(bufferSize int) *Producer {
	return &Producer{
		messages: make(chan Message, bufferSize),
	}
}

func (p *Producer) Push(msg Message) {
	p.messages <- msg
}

type Consumer struct {
	messages chan Message
}

func NewConsumer(messages chan Message) *Consumer {
	return &Consumer{
		messages: messages,
	}
}

func (c *Consumer) Pull() Message {
	return <-c.messages
}

func ExamplePushPull() {
	producer := NewProducer(10)
	consumer := NewConsumer(producer.messages)

	// 生产者推送消息
	go func() {
		for i := 0; i < 5; i++ {
			msg := Message{Content: fmt.Sprintf("Message %d", i+1)}
			producer.Push(msg)
		}
		close(producer.messages)
	}()

	// 消费者拉取消息
	for {
		msg, ok := <-consumer.messages
		if !ok {
			break
		}
		println("Consumed:", msg.Content)
	}
}

// WebSocket 示例代码参考 (简化版)
// func (s *Server) handleConnection(conn *websocket.Conn) {
// 	return
// }

// func generateID() string {
// 	return "client_" + time.Now().Format("20060102150405")
// }
