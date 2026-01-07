package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cg917658910/go-study/socket-io/server"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 创建Socket.IO服务器
	socketServer := server.NewSocketIOServer()

	// 创建Gin路由
	router := gin.Default()

	// 配置CORS（允许跨域）
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 设置路由
	setupRoutes(router, socketServer)

	// 启动服务器
	port := getPort()
	log.Printf("🚀 Server starting on :%s", port)
	log.Printf("📡 WebSocket endpoint: ws://localhost:%s/socket.io/", port)
	log.Printf("🌐 Web interface: http://localhost:%s/", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}

func setupRoutes(router *gin.Engine, socketServer *server.SocketServer) {
	// 主页
	router.GET("/", func(c *gin.Context) {
		html := `
<!DOCTYPE html>
<html>
<head>
    <title>Go Socket.IO Demo</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background: #f0f2f5;
        }
        .container {
            background: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 { color: #333; }
        .status { 
            padding: 10px; 
            border-radius: 5px;
            margin: 10px 0;
        }
        .connected { background: #d4edda; color: #155724; }
        .disconnected { background: #f8d7da; color: #721c24; }
        #messages {
            height: 300px;
            overflow-y: auto;
            border: 1px solid #ddd;
            padding: 10px;
            margin: 10px 0;
            background: #fafafa;
        }
        input, button {
            padding: 10px;
            margin: 5px;
            border: 1px solid #ddd;
            border-radius: 5px;
        }
        button {
            background: #007bff;
            color: white;
            border: none;
            cursor: pointer;
        }
        button:hover { background: #0056b3; }
        button:disabled { background: #ccc; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 Go Socket.IO Demo</h1>
        <p>Open browser console to see connection details</p>
        <div id="status" class="status disconnected">Disconnected</div>
        <div id="messages"></div>
        <input type="text" id="message" placeholder="Type a message...">
        <button onclick="sendMessage()" disabled>Send</button>
        <button onclick="connectSocket()">Connect</button>
        <button onclick="disconnectSocket()">Disconnect</button>
    </div>

    <script src="https://cdn.socket.io/4.5.4/socket.io.min.js"></script>
    <script>
        let socket = null;
        
        function connectSocket() {
            if (socket?.connected) {
                alert("Already connected!");
                return;
            }
            
            // 连接到Socket.IO服务器
            socket = io(window.location.origin, {
                path: '/socket.io/',
                transports: ['websocket', 'polling'],
                reconnection: true,
                reconnectionAttempts: 5,
                reconnectionDelay: 1000
            });
            
            socket.on('connect', () => {
                console.log('✅ Connected to server:', socket.id);
                updateStatus(true);
                document.querySelector('button').disabled = false;
                
                // 添加欢迎消息
                const messages = document.getElementById('messages');
                messages.innerHTML += '<p>✅ Connected to server!</p>';
            });
            
            socket.on('disconnect', (reason) => {
                console.log('❌ Disconnected:', reason);
                updateStatus(false);
                document.querySelector('button').disabled = true;
                
                const messages = document.getElementById('messages');
                messages.innerHTML += '<p>❌ Disconnected from server</p>';
            });
            
            socket.on('connect_error', (error) => {
                console.error('❌ Connection error:', error);
                const messages = document.getElementById('messages');
                messages.innerHTML += '<p>❌ Connection error: ' + error.message + '</p>';
            });
            
            socket.on('welcome', (data) => {
                console.log('📩 Welcome message:', data);
                const messages = document.getElementById('messages');
                messages.innerHTML += '<p>📩 Server says: ' + data.message + '</p>';
                messages.scrollTop = messages.scrollHeight;
            });
            
            socket.on('chat', (data) => {
                console.log('💬 Chat message:', data);
                const messages = document.getElementById('messages');
                messages.innerHTML += '<p><strong>' + data.from + ':</strong> ' + data.message + '</p>';
                messages.scrollTop = messages.scrollHeight;
            });
            
            socket.on('user_joined', (data) => {
                console.log('👋 User joined:', data);
                const messages = document.getElementById('messages');
                messages.innerHTML += '<p>👋 ' + data.username + ' joined the chat</p>';
                messages.scrollTop = messages.scrollHeight;
            });
            
            socket.on('user_left', (data) => {
                console.log('👋 User left:', data);
                const messages = document.getElementById('messages');
                messages.innerHTML += '<p>👋 ' + data.username + ' left the chat</p>';
                messages.scrollTop = messages.scrollHeight;
            });
        }
        
        function disconnectSocket() {
            if (socket) {
                socket.disconnect();
            }
        }
        
        function sendMessage() {
            const input = document.getElementById('message');
            const message = input.value.trim();
            
            if (message && socket?.connected) {
                socket.emit('chat', message);
                const messages = document.getElementById('messages');
                messages.innerHTML += '<p><strong>You:</strong> ' + message + '</p>';
                input.value = '';
                messages.scrollTop = messages.scrollHeight;
            }
        }
        
        function updateStatus(connected) {
            const status = document.getElementById('status');
            status.textContent = connected ? '✅ Connected' : '❌ Disconnected';
            status.className = 'status ' + (connected ? 'connected' : 'disconnected');
        }
        
        // 页面加载时自动连接
        window.addEventListener('load', connectSocket);
    </script>
</body>
</html>`
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, html)
	})

	// Socket.IO处理器 - 必须放在Gin包装中
	router.GET("/socket.io/*any", func(c *gin.Context) {
		socketServer.ServeHTTP(c.Writer, c.Request)
	})

	router.POST("/socket.io/*any", func(c *gin.Context) {
		socketServer.ServeHTTP(c.Writer, c.Request)
	})

	// 健康检查端点
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"clients": len(socketServer.Clients),
			"uptime":  "running",
		})
	})

	// 获取在线用户
	router.GET("/users", func(c *gin.Context) {
		users := socketServer.GetOnlineUsers()
		c.JSON(200, gin.H{
			"count": len(users),
			"users": users,
		})
	})
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
