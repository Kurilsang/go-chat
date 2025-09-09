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

	logrus.Info("开始执行数据库迁移...")

	// 自动迁移User表
	if err := db.AutoMigrate(&User{}).Error; err != nil {
		logrus.Errorf("User表迁移失败: %v", err)
		return
	}

	logrus.Info("✅ User表迁移完成")
	logrus.Info("🎉 数据库迁移全部完成")
}

// GetDB 获取数据库连接 (保持向后兼容性)
func GetDB() *gorm.DB {
	return global.GetMySQLClient()
}
