package websocket

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// MessageType 消息类型枚举
type MessageType string

const (
	// 系统消息
	MessageTypeJoin     MessageType = "join"      // 用户加入
	MessageTypeLeave    MessageType = "leave"     // 用户离开
	MessageTypeUserList MessageType = "user_list" // 在线用户列表
	MessageTypeError    MessageType = "error"     // 错误消息

	// 聊天消息
	MessageTypePrivate   MessageType = "private"   // 私聊消息
	MessageTypeHeartbeat MessageType = "heartbeat" // 心跳检测
	MessageTypeTyping    MessageType = "typing"    // 正在输入
	MessageTypeRead      MessageType = "read"      // 消息已读
)

// Message WebSocket消息结构
type Message struct {
	ID         string      `json:"id"`             // 消息唯一ID
	Type       MessageType `json:"type"`           // 消息类型
	FromUserID uint        `json:"from_user_id"`   // 发送者ID
	ToUserID   uint        `json:"to_user_id"`     // 接收者ID (私聊时使用)
	Content    string      `json:"content"`        // 消息内容
	Timestamp  time.Time   `json:"timestamp"`      // 时间戳
	Data       interface{} `json:"data,omitempty"` // 附加数据
}

// PrivateMessageData 私聊消息附加数据
type PrivateMessageData struct {
	MessageID string `json:"message_id"` // MongoDB中的消息ID
	SessionID string `json:"session_id"` // 会话ID
	IsOffline bool   `json:"is_offline"` // 是否离线消息
}

// UserListData 用户列表附加数据
type UserListData struct {
	Users []OnlineUser `json:"users"`
}

// OnlineUser 在线用户信息
type OnlineUser struct {
	UserID   uint      `json:"user_id"`
	UserName string    `json:"username"`
	Avatar   string    `json:"avatar"`
	Status   string    `json:"status"`
	LastSeen time.Time `json:"last_seen"`
}

// ErrorData 错误消息附加数据
type ErrorData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewMessage 创建新消息
func NewMessage(msgType MessageType, fromUserID, toUserID uint, content string) *Message {
	return &Message{
		ID:         generateMessageID(),
		Type:       msgType,
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Content:    content,
		Timestamp:  time.Now(),
	}
}

// NewSystemMessage 创建系统消息
func NewSystemMessage(msgType MessageType, data interface{}) *Message {
	return &Message{
		ID:        generateMessageID(),
		Type:      msgType,
		Timestamp: time.Now(),
		Data:      data,
	}
}

// NewErrorMessage 创建错误消息
func NewErrorMessage(code int, message string) *Message {
	return &Message{
		ID:        generateMessageID(),
		Type:      MessageTypeError,
		Timestamp: time.Now(),
		Data: ErrorData{
			Code:    code,
			Message: message,
		},
	}
}

// ToJSON 消息序列化为JSON
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// FromJSON 从JSON反序列化消息
func FromJSON(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return &msg, err
}

// generateMessageID 生成消息ID
func generateMessageID() string {
	return fmt.Sprintf("%d_%s", time.Now().UnixNano(),
		generateRandomString(8))
}

// generateRandomString 生成随机字符串
func generateRandomString(length int) string {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)[:length]
}

// generateClientID 生成客户端ID
func generateClientID() string {
	return fmt.Sprintf("client_%d_%s", time.Now().UnixNano(),
		generateRandomString(6))
}
