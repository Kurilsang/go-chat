package model

import (
	"go_chat/global"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// Migration 执行数据库迁移
func Migration() {
	db := global.GetMySQLClient()
	if db == nil {
		logrus.Error("MySQL客户端未初始化，无法执行数据库迁移")
		return
	}

	// 自动迁移模式
	err := db.AutoMigrate(&User{})
	if err != nil {
		logrus.Error("数据库迁移失败:", err)
		return
	}
	logrus.Info("数据库表迁移完成")
}

// GetDB 获取数据库连接 (保持向后兼容性)
func GetDB() *gorm.DB {
	return global.GetMySQLClient()
}
