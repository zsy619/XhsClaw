// Package repository 提供数据访问层
package repository

import (
	"log"
	"xiaohongshu/internal/config"
	"xiaohongshu/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase(cfg *config.DatabaseConfig) error {
	var err error
	
	// 配置GORM日志
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// 根据数据库类型选择驱动
	var dialector gorm.Dialector
	if cfg.Type == "mysql" {
		dialector = mysql.Open(cfg.GetDSN())
	} else {
		dialector = postgres.Open(cfg.GetDSN())
	}

	// 连接数据库
	DB, err = gorm.Open(dialector, gormConfig)
	if err != nil {
		return err
	}

	// 自动迁移表结构
	err = DB.AutoMigrate(
		&model.User{},
		&model.UserConfig{},
		&model.Content{},
		&model.PublishRecord{},
		&model.TokenBlacklist{},
	)
	if err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
