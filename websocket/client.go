package websocket

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	// 时间配置
	writeWait      = 10 * time.Second    // 写入超时
	pongWait       = 60 * time.Second    // Pong等待时间
	pingPeriod     = (pongWait * 9) / 10 // Ping发送间隔
	maxMessageSize = 512                 // 最大消息大小
)

// Client 代表一个WebSocket客户端连接
type Client struct {
	// 基础信息
	ID       string `json:"id"`       // 客户端唯一ID
	UserID   uint   `json:"user_id"`  // 用户ID
	UserName string `json:"username"` // 用户名
	Avatar   string `json:"avatar"`   // 头像

	// 连接管理
	Hub  *Hub            `json:"-"` // 连接池引用
	Conn *websocket.Conn `json:"-"` // WebSocket连接
	Send chan []byte     `json:"-"` // 发送消息通道

	// 状态管理
	IsAlive     bool      `json:"is_alive"`     // 连接状态
	LastSeen    time.Time `json:"last_seen"`    // 最后活跃时间
	ConnectedAt time.Time `json:"connected_at"` // 连接时间

	// 并发控制
	mu sync.RWMutex `json:"-"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 在生产环境中应该检查Origin
		return true
	},
}

// NewClient 创建新的客户端连接
func NewClient(hub *Hub, conn *websocket.Conn, userID uint, userName string) *Client {
	return &Client{
		ID:          generateClientID(),
		UserID:      userID,
		UserName:    userName,
		Hub:         hub,
		Conn:        conn,
		Send:        make(chan []byte, 256),
		IsAlive:     true,
		LastSeen:    time.Now(),
		ConnectedAt: time.Now(),
	}
}

// ReadPump 处理从客户端接收的消息
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	// 设置连接参数
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		c.updateLastSeen()
		return nil
	})

	for {
		_, messageData, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Errorf("WebSocket错误: %v", err)
			}
			break
		}

		// 解析消息
		message, err := FromJSON(messageData)
		if err != nil {
			logrus.Errorf("消息解析失败: %v", err)
			c.SendError(400, "消息格式错误")
			continue
		}

		// 设置发送者信息
		message.FromUserID = c.UserID
		message.Timestamp = time.Now()

		// 更新最后活跃时间
		c.updateLastSeen()

		// 处理不同类型的消息
		c.handleMessage(message)
	}
}

// WritePump 处理向客户端发送的消息
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 发送队列中的其他消息
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SendMessage 发送消息给客户端
func (c *Client) SendMessage(message *Message) {
	data, err := message.ToJSON()
	if err != nil {
		logrus.Errorf("消息序列化失败: %v", err)
		return
	}

	select {
	case c.Send <- data:
	default:
		// 发送失败，客户端可能已断开连接
		logrus.Warnf("向用户 %d 发送消息失败，连接可能已断开", c.UserID)
		go func() {
			c.Hub.Unregister <- c
		}()
	}
}

// SendError 发送错误消息
func (c *Client) SendError(code int, message string) {
	errorMsg := NewErrorMessage(code, message)
	c.SendMessage(errorMsg)
}

// handleMessage 处理接收到的消息
func (c *Client) handleMessage(message *Message) {
	logrus.Infof("用户 %d 发送消息: 类型=%s, 目标用户=%d, 内容=%s",
		c.UserID, message.Type, message.ToUserID, message.Content)

	switch message.Type {
	case MessageTypePrivate:
		c.Hub.HandlePrivateMessage(c, message)
	case MessageTypeHeartbeat:
		c.handleHeartbeat(message)
	case MessageTypeTyping:
		c.Hub.HandleTypingMessage(c, message)
	case MessageTypeRead:
		c.Hub.HandleReadMessage(c, message)
	default:
		logrus.Warnf("未知消息类型: %s", message.Type)
		c.SendError(400, "未知消息类型")
	}
}

// handleHeartbeat 处理心跳消息
func (c *Client) handleHeartbeat(message *Message) {
	response := NewMessage(MessageTypeHeartbeat, 0, 0, "pong")
	c.SendMessage(response)
}

// updateLastSeen 更新最后活跃时间
func (c *Client) updateLastSeen() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.LastSeen = time.Now()
}

// Close 关闭客户端连接
func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.IsAlive {
		c.IsAlive = false
		close(c.Send)
		c.Conn.Close()
	}
}

// ToOnlineUser 转换为在线用户信息
func (c *Client) ToOnlineUser() OnlineUser {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return OnlineUser{
		UserID:   c.UserID,
		UserName: c.UserName,
		Avatar:   c.Avatar,
		Status:   "online",
		LastSeen: c.LastSeen,
	}
}
