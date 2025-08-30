package config

import (
	"fmt"
	"go_chat/global"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

var (
	Db         string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassWord string
	DbName     string
)

// InitMySQL 初始化MySQL连接
func InitMySQL() {
	connString := GetMySQLConnectionString()
	db, err := gorm.Open("mysql", connString)
	if err != nil {
		logrus.Error("MySQL数据库连接失败:", err)
		logrus.Warn("MySQL服务不可用，相关功能将被禁用")
		return // 不panic，允许程序继续运行
	}

	// 设置日志模式
	db.LogMode(true)
	// 在生产环境中可以关闭日志: db.LogMode(false)

	// 数据库设置
	db.SingularTable(true)                       // 默认不加复数s
	db.DB().SetMaxIdleConns(20)                  // 设置连接池，空闲
	db.DB().SetMaxOpenConns(100)                 // 设置打开最大连接
	db.DB().SetConnMaxLifetime(time.Second * 30) // 连接最大生命周期

	// 设置全局MySQL客户端
	global.SetMySQLClient(db)
	logrus.Info("MySQL数据库连接成功")
}

// LoadMysqlData 加载MySQL配置数据
func LoadMysqlData(file *ini.File) {
	Db = file.Section("mysql").Key("Db").String()
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassWord = file.Section("mysql").Key("DbPassWord").String()
	DbName = file.Section("mysql").Key("DbName").String()
}

// GetMySQLConnectionString 获取MySQL连接字符串
func GetMySQLConnectionString() string {
	return strings.Join([]string{DbUser, ":", DbPassWord, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8&parseTime=true"}, "")
}

// PrintMySQLConfig 打印MySQL配置
func PrintMySQLConfig() {
	fmt.Println("\n=== MySQL配置 ===")
	fmt.Printf("数据库类型: %s\n", Db)
	fmt.Printf("主机地址: %s\n", DbHost)
	fmt.Printf("端口: %s\n", DbPort)
	fmt.Printf("用户名: %s\n", DbUser)
	fmt.Printf("密码: %s\n", DbPassWord)
	fmt.Printf("数据库名: %s\n", DbName)
	fmt.Printf("连接字符串: %s\n", GetMySQLConnectionString())
}
