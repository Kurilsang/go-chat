package config

import (
	"fmt"
	"go_chat/model"

	"gopkg.in/ini.v1"
)

var (
	AppMode  string
	HttpPort string
)

func Init() {
	// 从本地读取环境变量
	file, err := ini.Load("./config/config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径:", err)
		return
	}

	// 加载各模块配置
	LoadServer(file)
	LoadMysqlData(file)
	LoadMongoDB(file)
	LoadRedisData(file)

	fmt.Println("配置文件加载完成!")

	// 初始化数据库连接
	InitMySQL()
	InitMongoDB()
	InitRedis()

	// 执行MySQL数据库迁移
	model.Migration()

	// 打印配置信息
	PrintConfig()
}

func LoadServer(file *ini.File) {
	AppMode = file.Section("service").Key("AppMode").String()
	HttpPort = file.Section("service").Key("HttpPort").String()
}

func PrintConfig() {
	fmt.Println("=== 服务器配置 ===")
	fmt.Printf("应用模式: %s\n", AppMode)
	fmt.Printf("HTTP端口: %s\n", HttpPort)

	PrintMySQLConfig()
	PrintMongoDBConfig()
	PrintRedisConfig()
}
