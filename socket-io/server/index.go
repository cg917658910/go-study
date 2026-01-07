package server

import (
	"log"
	"net/http"
	"sync"
	"time"

	socketio "github.com/googollee/go-socket.io"
)

type Client struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	JoinedAt time.Time `json:"joined_at"`
}

type SocketServer struct {
	server  *socketio.Server
	Clients map[string]*Client
	mu      sync.RWMutex
}

func NewSocketIOServer() *SocketServer {
	// 创建Socket.IO服务器
	server := socketio.NewServer(nil)

	s := &SocketServer{
		server:  server,
		Clients: make(map[string]*Client),
	}

	s.setupHandlers()
	return s
}

func (s *SocketServer) setupHandlers() {
	// 连接事件
	s.server.OnConnect("/", func(conn socketio.Conn) error {
		clientID := conn.ID()

		s.mu.Lock()
		s.Clients[clientID] = &Client{
			ID:       clientID,
			Username: "User_" + clientID[:6],
			JoinedAt: time.Now(),
		}
		s.mu.Unlock()

		log.Printf("✅ Client connected: %s", clientID)

		// 发送欢迎消息
		conn.Emit("welcome", map[string]interface{}{
			"message":   "Welcome to Go Socket.IO Server!",
			"id":        clientID,
			"timestamp": time.Now().Unix(),
		})

		// 广播用户加入消息
		s.Broadcast("user_joined", map[string]interface{}{
			"id":       clientID,
			"username": s.Clients[clientID].Username,
			"time":     time.Now().Format("15:04:05"),
		})

		// 发送当前在线用户列表
		onlineUsers := s.GetOnlineUsers()
		conn.Emit("users_online", onlineUsers)

		return nil
	})

	// 处理聊天消息
	s.server.OnEvent("/", "chat", func(conn socketio.Conn, msg string) {
		clientID := conn.ID()
		username := "unknown"

		s.mu.RLock()
		if client, exists := s.Clients[clientID]; exists {
			username = client.Username
		}
		s.mu.RUnlock()

		log.Printf("💬 Chat message from %s (%s): %s", username, clientID, msg)

		// 广播消息给所有客户端
		s.Broadcast("chat", map[string]interface{}{
			"from":      username,
			"message":   msg,
			"id":        clientID,
			"timestamp": time.Now().Unix(),
		})
	})

	// 处理设置用户名
	s.server.OnEvent("/", "set_username", func(conn socketio.Conn, username string) {
		clientID := conn.ID()
		oldUsername := ""

		s.mu.Lock()
		if client, exists := s.Clients[clientID]; exists {
			oldUsername = client.Username
			client.Username = username
		}
		s.mu.Unlock()

		log.Printf("👤 Username changed: %s -> %s", oldUsername, username)

		conn.Emit("username_updated", map[string]interface{}{
			"old": oldUsername,
			"new": username,
		})

		// 广播用户名变更
		s.Broadcast("username_changed", map[string]interface{}{
			"id":  clientID,
			"old": oldUsername,
			"new": username,
		})
	})

	// 处理断开连接
	s.server.OnDisconnect("/", func(conn socketio.Conn, reason string) {
		clientID := conn.ID()
		username := "unknown"

		s.mu.Lock()
		if client, exists := s.Clients[clientID]; exists {
			username = client.Username
			delete(s.Clients, clientID)
		}
		s.mu.Unlock()

		log.Printf("❌ Client disconnected: %s (%s), reason: %s", username, clientID, reason)

		// 广播用户离开
		s.Broadcast("user_left", map[string]interface{}{
			"id":       clientID,
			"username": username,
			"time":     time.Now().Format("15:04:05"),
		})
	})

	// 错误处理
	s.server.OnError("/", func(conn socketio.Conn, err error) {
		log.Printf("⚠️ Socket error: %v", err)
	})

	// 心跳检测
	s.server.OnEvent("/", "ping", func(conn socketio.Conn) {
		conn.Emit("pong", map[string]interface{}{
			"timestamp": time.Now().Unix(),
			"server":    "Go Socket.IO Server",
		})
	})
}

// 广播消息给所有客户端
func (s *SocketServer) Broadcast(event string, data interface{}) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 广播给所有连接的客户端
	s.server.BroadcastToNamespace("/", event, data)
}

// 获取在线用户列表
func (s *SocketServer) GetOnlineUsers() []Client {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]Client, 0, len(s.Clients))
	for _, client := range s.Clients {
		users = append(users, *client)
	}

	return users
}

// HTTP处理器
func (s *SocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("📨 %s %s", r.Method, r.URL.Path)
	s.server.ServeHTTP(w, r)
}
