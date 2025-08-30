package router

import (
	"go_chat/controller"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.Default()

	// 设置中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

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
	}

	return r
}
