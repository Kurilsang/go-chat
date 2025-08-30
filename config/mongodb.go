package config

import (
	"context"
	"fmt"
	"go_chat/global"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/ini.v1"
)

var (
	MongoDBName string
	MongoDBAddr string
	MongoDBPwd  string
	MongoDBPort string
)

// InitMongoDB 初始化MongoDB连接
func InitMongoDB() {
	// 设置mongoDB客户端连接信息
	clientOptions := options.Client().ApplyURI("mongodb://" + MongoDBAddr + ":" + MongoDBPort)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logrus.Error("MongoDB连接失败:", err)
		logrus.Warn("MongoDB服务不可用，相关功能将被禁用")
		return // 不panic，允许程序继续运行
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		logrus.Error("MongoDB ping失败:", err)
		logrus.Warn("MongoDB服务不可用，相关功能将被禁用")
		return // 不panic，允许程序继续运行
	}

	// 设置全局MongoDB客户端
	global.SetMongoDBClient(client)
	logrus.Info("MongoDB连接成功")
}

// LoadMongoDB 加载MongoDB配置数据
func LoadMongoDB(file *ini.File) {
	MongoDBName = file.Section("MongoDB").Key("MongoDBName").String()
	MongoDBAddr = file.Section("MongoDB").Key("MongoDBAddr").String()
	MongoDBPwd = file.Section("MongoDB").Key("MongoDBPwd").String()
	MongoDBPort = file.Section("MongoDB").Key("MongoDBPort").String()
}

// PrintMongoDBConfig 打印MongoDB配置
func PrintMongoDBConfig() {
	fmt.Println("\n=== MongoDB配置 ===")
	fmt.Printf("数据库名: %s\n", MongoDBName)
	fmt.Printf("地址: %s\n", MongoDBAddr)
	fmt.Printf("密码: %s\n", MongoDBPwd)
	fmt.Printf("端口: %s\n", MongoDBPort)
}
