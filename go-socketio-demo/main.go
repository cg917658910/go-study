package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
	id   string
}

type Server struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (s *Server) run() {
	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			s.clients[client] = true
			s.mu.Unlock()
			log.Printf("✅ Client connected: %s", client.id)

			// 发送欢迎消息
			welcome, _ := json.Marshal(map[string]interface{}{
				"type": "welcome",
				"id":   client.id,
				"time": time.Now().Unix(),
			})
			client.send <- welcome

		case client := <-s.unregister:
			s.mu.Lock()
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.send)
			}
			s.mu.Unlock()
			log.Printf("❌ Client disconnected: %s", client.id)

		case message := <-s.broadcast:
			s.mu.RLock()
			for client := range s.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(s.clients, client)
				}
			}
			s.mu.RUnlock()
		}
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
		id:   generateID(),
	}

	s.register <- client

	// 读取协程
	go func() {
		defer func() {
			s.unregister <- client
			conn.Close()
		}()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}

			// 处理消息
			var msg map[string]interface{}
			if err := json.Unmarshal(message, &msg); err == nil {
				log.Printf("📨 Received: %v", msg)

				// 广播消息
				broadcastMsg, _ := json.Marshal(map[string]interface{}{
					"type":    "message",
					"from":    client.id,
					"content": msg["content"],
					"time":    time.Now().Unix(),
				})
				s.broadcast <- broadcastMsg
			}
		}
	}()

	// 写入协程
	go func() {
		defer conn.Close()
		for {
			select {
			case message, ok := <-client.send:
				if !ok {
					return
				}
				conn.WriteMessage(websocket.TextMessage, message)
			}
		}
	}()
}

func generateID() string {
	return "client_" + time.Now().Format("20060102150405")
}

func main() {
	server := NewServer()
	go server.run()

	r := gin.Default()

	// 静态文件
	r.Static("/static", "./static")

	// 主页
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// WebSocket 端点
	r.GET("/ws", gin.WrapH(server))

	log.Println("🚀 Server starting on :8080")
	r.Run(":8090")
}
