package router

import (
	"go_chat/controller"
	"go_chat/websocket"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.Default()

	// 设置中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 静态文件服务
	r.Static("/static", "./static")
	r.LoadHTMLGlob("static/*.html")

	// 创建WebSocket Hub
	hub := websocket.NewHub()
	go hub.Run()

	// 创建WebSocket处理器
	wsHandler := websocket.NewHandler(hub)

	// WebSocket路由
	r.GET("/ws", wsHandler.HandleWebSocket)

	// API v1 路由分组
	v1 := r.Group("/api/v1")
	{
		// 认证相关路由
		auth := v1.Group("/auth")
		{
			auth.POST("/register", controller.Register)
			auth.POST("/login", controller.Login)
		}

		// 用户相关路由 (需要认证)
		user := v1.Group("/user")
		{
			// user.Use(middleware.AuthRequired()) // 后续可以添加认证中间件
			user.GET("/profile", controller.GetProfile)
		}

		// WebSocket相关路由
		ws := v1.Group("/ws")
		{
			ws.GET("/online-users", wsHandler.GetOnlineUsers)
		}
	}

	// 测试页面路由
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "Go Chat 测试页面",
		})
	})

	return r
}
