package websocket

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Handler WebSocket处理器
type Handler struct {
	Hub *Hub
}

// NewHandler 创建新的WebSocket处理器
func NewHandler(hub *Hub) *Handler {
	return &Handler{
		Hub: hub,
	}
}

// HandleWebSocket 处理WebSocket连接升级
func (h *Handler) HandleWebSocket(c *gin.Context) {
	logrus.Infof("收到WebSocket连接请求: %s", c.Request.URL.String())

	// 从查询参数获取用户信息（暂时用于测试）
	userIDStr := c.Query("user_id")
	userName := c.Query("username")

	logrus.Infof("连接参数: userID=%s, userName=%s", userIDStr, userName)

	if userIDStr == "" || userName == "" {
		logrus.Warnf("连接参数不完整: userID=%s, userName=%s", userIDStr, userName)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id和username参数不能为空",
		})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		logrus.Errorf("userID解析失败: %s, 错误: %v", userIDStr, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id格式错误",
		})
		return
	}

	logrus.Infof("准备升级WebSocket连接...")

	// 升级HTTP连接为WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.Errorf("WebSocket升级失败: %v", err)
		return
	}

	logrus.Infof("WebSocket升级成功，创建客户端连接...")

	// 创建客户端连接
	client := NewClient(h.Hub, conn, uint(userID), userName)

	logrus.Infof("客户端创建完成，准备注册到Hub...")

	// 注册客户端
	h.Hub.Register <- client

	logrus.Infof("客户端已发送到注册通道，启动协程...")

	// 启动客户端协程
	go client.WritePump()
	go client.ReadPump()

	logrus.Infof("WebSocket连接处理完成：用户 %s (ID: %d)", userName, userID)
}

// GetOnlineUsers 获取在线用户列表API
func (h *Handler) GetOnlineUsers(c *gin.Context) {
	users := h.Hub.getOnlineUsers()
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取在线用户列表成功",
		"data": gin.H{
			"users": users,
			"count": len(users),
		},
	})
}
