package handler

import (
	"Vnollx/cmd/api/rpc"
	"Vnollx/kitex_gen/base"
	"Vnollx/kitex_gen/user"
	"Vnollx/pkg/middlerware"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// 定义 WebSocket 连接的结构体
type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	userid int64
}

// 定义 WebSocket 服务器
type WebSocketServer struct {
	clients    map[int64]*Client // 连接的客户端
	broadcast  chan []byte       // 广播消息的通道
	register   chan *Client      // 新客户端的通道
	unregister chan *Client      // 断开连接的客户端通道
	mu         sync.Mutex        // 锁，保护 clients 数据
}

var server *WebSocketServer

func init() {
	server = NewWebSocketServer()
	go HandleMessages()
}

// 创建新的 WebSocket 服务器
func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		clients:    make(map[int64]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}
func findClientByUserID(userID int64) *Client {
	server.mu.Lock()
	defer server.mu.Unlock()

	return server.clients[userID] // 直接根据 userID 查找目标客户端
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketConnections(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		logger.Error("Token缺失")
		ctx.JSON(http.StatusBadRequest, base.BaseResponse{
			Code: http.StatusBadRequest,
			Msg:  "Token缺失",
		})
		return
	}

	mc, err := middlerware.ParseToken(token)
	if err != nil {
		logger.Error("token解析失败：", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, base.BaseResponse{
			Code: http.StatusBadRequest,
			Msg:  "Token解析失败",
		})
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}

	client := &Client{
		conn:   conn,
		send:   make(chan []byte),
		userid: mc.UserId,
	}

	server.register <- client

	go client.ReadMessages()
	go client.SendMessages()
	select {}
	// 这里可以考虑添加逻辑来处理连接关闭
	// 比如使用 context 或者监听客户端断开事件
}

// 读取客户端的消息
func (client *Client) ReadMessages() {
	for {
		messageType, message, err := client.conn.ReadMessage()
		messageStr := string(message)
		//log.Println(messageStr)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Println("Connection closed normally")
			} else {
				log.Println("Error reading message:", err)
			}
			break
		}
		if messageType == websocket.TextMessage {
			// 假设前端传过来的消息是 JSON 格式
			// 示例：{"token": "user_token", "touserid": "123", "message": "Hello!"}
			var msgData struct {
				Type     string
				Token    string
				Touserid string
				Content  string
			}

			// 解析消息
			err = json.Unmarshal([]byte(messageStr), &msgData)
			if err != nil {
				log.Println("Error unmarshalling message:", err)
				continue
			}
			touser_id, _ := strconv.ParseInt(msgData.Touserid, 10, 64)

			// 调用 RPC 服务将消息存入数据库
			mc, err := middlerware.ParseToken(msgData.Token)
			if err != nil {
				log.Println("无效的token:", err)
				continue
			}
			req := &user.SendMessageRequest{
				UserId:   mc.UserId,
				TouserId: touser_id,
				Content:  msgData.Content,
			}
			//log.Println(req.Content)
			_, _ = rpc.SendMessage(context.Background(), req)
			targetClient := findClientByUserID(touser_id)
			if targetClient != nil {
				targetClient.send <- message
			}
		}
	}
}

// 向客户端发送消息
func (client *Client) SendMessages() {
	for message := range client.send {
		err := client.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Error sending message:", err)
			break
		}
	}
}

func HandleMessages() {
	for {
		select {
		case client := <-server.register:
			server.mu.Lock()
			server.clients[client.userid] = client // 使用 userID 作为 key 注册
			server.mu.Unlock()
			log.Println("New client connected")
		case client := <-server.unregister:
			server.mu.Lock()
			delete(server.clients, client.userid) // 通过 userID 注销
			close(client.send)
			server.mu.Unlock()
			log.Println("Client disconnected")
		}
		// 不再处理 broadcast 消息，改为在 ReadMessages 中直接发送给目标客户端
	}
}
