package config

import (
	"fmt"
	"go_chat/global"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

var (
	RedisDb     string
	RedisAddr   string
	RedisPw     string
	RedisDbName string
)

// InitRedis 初始化Redis连接
func InitRedis() {
	db, _ := strconv.ParseUint(RedisDbName, 10, 64)
	client := redis.NewClient(&redis.Options{
		Addr:     RedisAddr,
		Password: RedisPw, // 如果没有密码，RedisPw为空字符串
		DB:       int(db),
	})

	// 测试连接
	_, err := client.Ping().Result()
	if err != nil {
		logrus.Error("Redis连接失败:", err)
		logrus.Warn("Redis服务不可用，相关功能将被禁用")
		return // 不panic，允许程序继续运行
	}

	// 设置全局Redis客户端
	global.SetRedisClient(client)
	logrus.Info("Redis连接成功")
}

// LoadRedisData 加载Redis配置数据
func LoadRedisData(file *ini.File) {
	RedisDb = file.Section("redis").Key("RedisDb").String()
	RedisAddr = file.Section("redis").Key("RedisAddr").String()
	RedisPw = file.Section("redis").Key("RedisPw").String()
	RedisDbName = file.Section("redis").Key("RedisDbName").String()
}

// PrintRedisConfig 打印Redis配置
func PrintRedisConfig() {
	fmt.Println("\n=== Redis配置 ===")
	fmt.Printf("Redis类型: %s\n", RedisDb)
	fmt.Printf("Redis地址: %s\n", RedisAddr)
	fmt.Printf("Redis密码: %s\n", RedisPw)
	fmt.Printf("Redis数据库: %s\n", RedisDbName)
}
