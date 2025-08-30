package websocket

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Hub 维护活跃客户端集合并向客户端广播消息
type Hub struct {
	// 客户端管理
	Clients     map[*Client]bool // 已注册的客户端
	UserClients map[uint]*Client // 用户ID到客户端的映射

	// 消息通道
	Register   chan *Client  // 客户端注册通道
	Unregister chan *Client  // 客户端注销通道
	Broadcast  chan *Message // 广播消息通道

	// 私聊处理
	PrivateMessage chan *PrivateMessageRequest // 私聊消息通道

	// 并发控制
	mu sync.RWMutex
}

// PrivateMessageRequest 私聊消息请求
type PrivateMessageRequest struct {
	From    *Client
	Message *Message
}

// NewHub 创建新的Hub实例
func NewHub() *Hub {
	return &Hub{
		Clients:        make(map[*Client]bool),
		UserClients:    make(map[uint]*Client),
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		Broadcast:      make(chan *Message),
		PrivateMessage: make(chan *PrivateMessageRequest),
	}
}

// Run 运行Hub主循环
func (h *Hub) Run() {
	logrus.Info("🚀 Hub主循环开始运行...")
	// 启动清理协程
	go h.startCleanupRoutine()

	for {
		select {
		case client := <-h.Register:
			logrus.Infof("📝 Hub收到注册请求: 用户 %s (ID: %d)", client.UserName, client.UserID)
			h.registerClient(client)

		case client := <-h.Unregister:
			h.unregisterClient(client)

		case message := <-h.Broadcast:
			h.broadcastMessage(message)

		case privateMsg := <-h.PrivateMessage:
			h.handlePrivateMessage(privateMsg)
		}
	}
}

// registerClient 注册新客户端
func (h *Hub) registerClient(client *Client) {
	logrus.Infof("Hub开始注册客户端: 用户 %s (ID: %d)", client.UserName, client.UserID)

	h.mu.Lock()
	defer h.mu.Unlock()

	// 检查是否已有同用户的连接
	if existingClient, exists := h.UserClients[client.UserID]; exists {
		logrus.Infof("用户 %d 重新连接，关闭旧连接", client.UserID)
		existingClient.Close()
		delete(h.Clients, existingClient)
	}

	// 注册新客户端
	h.Clients[client] = true
	h.UserClients[client.UserID] = client

	logrus.Infof("用户 %s (ID: %d) 已连接，当前在线用户: %d",
		client.UserName, client.UserID, len(h.Clients))

	// 发送用户列表给新用户
	logrus.Infof("发送用户列表给新用户 %d", client.UserID)
	h.sendUserListToClient(client)

	// 通知其他用户有新用户加入
	logrus.Infof("广播用户加入消息: %d", client.UserID)
	h.broadcastUserJoin(client)

	logrus.Infof("用户 %d 注册完成", client.UserID)
}

// unregisterClient 注销客户端
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.Clients[client]; ok {
		delete(h.Clients, client)
		delete(h.UserClients, client.UserID)
		client.Close()

		logrus.Infof("用户 %s (ID: %d) 已断开连接，当前在线用户: %d",
			client.UserName, client.UserID, len(h.Clients))

		// 通知其他用户有用户离开
		h.broadcastUserLeave(client)
	}
}

// broadcastMessage 广播消息给所有客户端
func (h *Hub) broadcastMessage(message *Message) {
	h.mu.RLock()
	clients := make([]*Client, 0, len(h.Clients))
	for client := range h.Clients {
		clients = append(clients, client)
	}
	h.mu.RUnlock()

	// 发送给所有客户端
	for _, client := range clients {
		client.SendMessage(message)
	}
}

// handlePrivateMessage 处理私聊消息
func (h *Hub) handlePrivateMessage(req *PrivateMessageRequest) {
	h.mu.RLock()
	targetClient, exists := h.UserClients[req.Message.ToUserID]
	h.mu.RUnlock()

	if !exists {
		// 目标用户不在线，发送离线消息提示
		req.From.SendError(404, "目标用户不在线")
		logrus.Infof("私聊消息发送失败：用户 %d 不在线", req.Message.ToUserID)
		return
	}

	// 发送消息给目标用户
	targetClient.SendMessage(req.Message)

	// 发送确认消息给发送者
	ackMessage := NewMessage(MessageTypeRead, req.Message.ToUserID, req.Message.FromUserID, "")
	ackMessage.Data = map[string]interface{}{
		"original_message_id": req.Message.ID,
		"status":              "delivered",
	}
	req.From.SendMessage(ackMessage)

	logrus.Infof("私聊消息已发送：%d -> %d", req.Message.FromUserID, req.Message.ToUserID)
}

// HandlePrivateMessage 处理来自客户端的私聊消息
func (h *Hub) HandlePrivateMessage(from *Client, message *Message) {
	// 验证消息
	if message.ToUserID == 0 {
		from.SendError(400, "目标用户ID不能为空")
		return
	}

	if message.ToUserID == from.UserID {
		from.SendError(400, "不能向自己发送消息")
		return
	}

	// 发送到私聊处理通道
	h.PrivateMessage <- &PrivateMessageRequest{
		From:    from,
		Message: message,
	}
}

// HandleTypingMessage 处理正在输入消息
func (h *Hub) HandleTypingMessage(from *Client, message *Message) {
	h.mu.RLock()
	targetClient, exists := h.UserClients[message.ToUserID]
	h.mu.RUnlock()

	if exists {
		targetClient.SendMessage(message)
	}
}

// HandleReadMessage 处理消息已读
func (h *Hub) HandleReadMessage(from *Client, message *Message) {
	// 这里可以添加消息已读的处理逻辑
	// 比如更新数据库中的消息状态
	logrus.Infof("用户 %d 已读消息", from.UserID)
}

// sendUserListToClient 发送在线用户列表给指定客户端
func (h *Hub) sendUserListToClient(client *Client) {
	users := h.getOnlineUsers()
	userListMessage := NewSystemMessage(MessageTypeUserList, UserListData{Users: users})
	client.SendMessage(userListMessage)
}

// broadcastUserJoin 广播用户加入消息
func (h *Hub) broadcastUserJoin(client *Client) {
	joinMessage := NewSystemMessage(MessageTypeJoin, client.ToOnlineUser())
	h.Broadcast <- joinMessage
}

// broadcastUserLeave 广播用户离开消息
func (h *Hub) broadcastUserLeave(client *Client) {
	leaveMessage := NewSystemMessage(MessageTypeLeave, client.ToOnlineUser())
	h.Broadcast <- leaveMessage
}

// getOnlineUsers 获取在线用户列表
func (h *Hub) getOnlineUsers() []OnlineUser {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := make([]OnlineUser, 0, len(h.Clients))
	for client := range h.Clients {
		if client.IsAlive {
			users = append(users, client.ToOnlineUser())
		}
	}
	return users
}

// GetOnlineUserCount 获取在线用户数量
func (h *Hub) GetOnlineUserCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.Clients)
}

// GetUserClient 根据用户ID获取客户端连接
func (h *Hub) GetUserClient(userID uint) (*Client, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	client, exists := h.UserClients[userID]
	return client, exists
}

// startCleanupRoutine 启动清理协程，定期清理断开的连接
func (h *Hub) startCleanupRoutine() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		h.cleanupDeadConnections()
	}
}

// cleanupDeadConnections 清理死连接
func (h *Hub) cleanupDeadConnections() {
	h.mu.Lock()
	defer h.mu.Unlock()

	now := time.Now()
	var toRemove []*Client

	for client := range h.Clients {
		// 检查连接是否超时（超过2分钟无活动）
		if now.Sub(client.LastSeen) > 2*time.Minute {
			toRemove = append(toRemove, client)
		}
	}

	for _, client := range toRemove {
		logrus.Warnf("清理超时连接：用户 %d", client.UserID)
		delete(h.Clients, client)
		delete(h.UserClients, client.UserID)
		client.Close()
	}
}
