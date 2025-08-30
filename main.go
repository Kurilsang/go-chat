package main

import (
	"fmt"
	"go_chat/config"
	"go_chat/global"
	"go_chat/model"
	"go_chat/router"

	"github.com/sirupsen/logrus"
)

func main() {
	// 初始化配置
	config.Init()

	// 演示使用全局数据库连接
	fmt.Println("\n=== 数据库连接状态 ===")

	// 检查MongoDB连接
	mongoClient := global.GetMongoDBClient()
	if mongoClient != nil {
		fmt.Println("MongoDB: 已连接")
	} else {
		fmt.Println("MongoDB: 未连接")
	}

	// 检查Redis连接
	redisClient := global.GetRedisClient()
	if redisClient != nil {
		fmt.Println("Redis: 已连接")
	} else {
		fmt.Println("Redis: 未连接")
	}

	// 检查MySQL连接和演示User模型
	db := model.GetDB()
	if db != nil {
		fmt.Println("MySQL: 已连接")
		demoUserModel()
	} else {
		fmt.Println("MySQL: 未连接")
		fmt.Println("\n=== MySQL数据库创建提示 ===")
		fmt.Println("请先创建MySQL数据库：")
		fmt.Println("1. 登录MySQL: mysql -u root -p")
		fmt.Println("2. 创建数据库: CREATE DATABASE go_chat;")
		fmt.Println("3. 重新运行程序")
	}

	// 启动HTTP服务器
	startServer()
}

// startServer 启动HTTP服务器
func startServer() {
	fmt.Println("\n=== 启动HTTP服务器 ===")

	// 初始化路由
	r := router.InitRouter()

	// 获取配置的端口
	port := config.HttpPort
	if port == "" {
		port = "8080"
	}

	fmt.Printf("服务器正在启动，端口: %s\n", port)
	fmt.Printf("API文档地址: http://localhost:%s/api/v1\n", port)
	fmt.Printf("注册接口: POST http://localhost:%s/api/v1/auth/register\n", port)
	fmt.Printf("登录接口: POST http://localhost:%s/api/v1/auth/login\n", port)
	fmt.Printf("用户信息: GET http://localhost:%s/api/v1/user/profile?id=1\n", port)

	// 启动服务器
	if err := r.Run(":" + port); err != nil {
		logrus.Fatal("服务器启动失败:", err)
	}
}

// demoUserModel 演示User模型的使用
func demoUserModel() {
	fmt.Println("\n=== User模型演示 ===")

	// 创建新用户
	user := &model.User{
		UserName: "testuser",
		Email:    "test@example.com",
		Phone:    "13800138000",
		Status:   model.Active,
		Avatar:   "https://example.com/avatar.png",
	}

	// 设置密码
	err := user.SetPassword("123456")
	if err != nil {
		logrus.Error("设置密码失败:", err)
		return
	}

	fmt.Printf("用户名: %s\n", user.UserName)
	fmt.Printf("邮箱: %s\n", user.Email)
	fmt.Printf("电话: %s\n", user.Phone)
	fmt.Printf("状态: %s\n", user.Status)
	fmt.Printf("头像: %s\n", user.AvatarURL())

	// 验证密码
	isValid := user.CheckPassword("123456")
	fmt.Printf("密码验证: %t\n", isValid)

	fmt.Println("User模型功能正常！")
}
