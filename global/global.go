package global

import (
	"sync"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// 数据库连接实例
	MongoDBClient *mongo.Client
	RedisClient   *redis.Client
	MySQLClient   *gorm.DB

	// 初始化控制
	mongoOnce sync.Once
	redisOnce sync.Once
	mysqlOnce sync.Once
)

// SetMongoDBClient 设置MongoDB客户端 (只初始化一次)
func SetMongoDBClient(client *mongo.Client) {
	mongoOnce.Do(func() {
		MongoDBClient = client
	})
}

// SetRedisClient 设置Redis客户端 (只初始化一次)
func SetRedisClient(client *redis.Client) {
	redisOnce.Do(func() {
		RedisClient = client
	})
}

// SetMySQLClient 设置MySQL客户端 (只初始化一次)
func SetMySQLClient(client *gorm.DB) {
	mysqlOnce.Do(func() {
		MySQLClient = client
	})
}

// GetMongoDBClient 获取MongoDB客户端
func GetMongoDBClient() *mongo.Client {
	return MongoDBClient
}

// GetRedisClient 获取Redis客户端
func GetRedisClient() *redis.Client {
	return RedisClient
}

// GetMySQLClient 获取MySQL客户端
func GetMySQLClient() *gorm.DB {
	return MySQLClient
}
