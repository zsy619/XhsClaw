// Package repository 提供数据访问层
package repository

import (
	"encoding/json"
	"log"
	"sync"
	"time"
	"xiaohongshu/internal/config"
	"xiaohongshu/internal/model"
	"xiaohongshu/internal/utils"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// DB 全局数据库连接实例
	DB *gorm.DB
	// once 确保数据库只初始化一次
	once sync.Once
)

// InitDatabase 初始化数据库连接
func InitDatabase(cfg *config.DatabaseConfig) error {
	var initErr error
	
	once.Do(func() {
		var err error
		
		// 根据运行模式配置日志级别
		logLevel := logger.Info
		if config.AppConfig.Server.Mode == "production" {
			logLevel = logger.Error
		}
		
		// 配置GORM日志
		gormConfig := &gorm.Config{
			Logger: logger.Default.LogMode(logLevel),
			// 启用准备语句缓存，提高性能
			PrepareStmt: true,
			// 禁用默认事务，提高性能（需要时手动开启）
			SkipDefaultTransaction: true,
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
			initErr = err
			return
		}

		// 获取底层数据库连接池
		sqlDB, err := DB.DB()
		if err != nil {
			initErr = err
			return
		}

		// 配置连接池
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
		sqlDB.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTime) * time.Second)

		// 验证连接
		if err := sqlDB.Ping(); err != nil {
			initErr = err
			return
		}

		log.Printf("Database connection pool configured: MaxOpenConns=%d, MaxIdleConns=%d",
			cfg.MaxOpenConns, cfg.MaxIdleConns)

		// 先迁移角色和权限表
		err = DB.AutoMigrate(
			&model.Role{},
			&model.Permission{},
		)
		if err != nil {
			initErr = err
			return
		}

		// 初始化角色和权限数据（在迁移用户表之前）
		if err := initPermissions(); err != nil {
			log.Printf("Failed to init permissions: %v", err)
		}
		if err := initRoles(); err != nil {
			log.Printf("Failed to init roles: %v", err)
		}

		// 临时禁用外键约束，处理旧数据
		if cfg.Type == "postgres" {
			DB.Exec("SET CONSTRAINTS ALL DEFERRED")
		}

		// 迁移用户表
		err = DB.AutoMigrate(
			&model.User{},
			&model.UserConfig{},
			&model.Content{},
			&model.ContentHistory{},
			&model.PublishRecord{},
			&model.TokenBlacklist{},
			&model.TokenUsage{},
		)
		if err != nil {
			initErr = err
			return
		}

		// 修复旧用户数据的 role_id（设置为默认的普通用户角色ID=3）
		var count int64
		DB.Model(&model.User{}).Where("role_id IS NULL OR role_id = 0").Count(&count)
		if count > 0 {
			log.Printf("Fixing %d users with invalid role_id...", count)
			DB.Model(&model.User{}).Where("role_id IS NULL OR role_id = 0").Update("role_id", 3)
		}

		log.Println("Database initialized successfully")
		
		// 初始化admin用户
		if err := initAdminUser(); err != nil {
			log.Printf("Failed to init admin user: %v", err)
		}
	})
	
	return initErr
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// InitData 初始化数据（保留用于向后兼容）
func InitData() error {
	// 初始化权限
	if err := initPermissions(); err != nil {
		return err
	}
	
	// 初始化角色
	if err := initRoles(); err != nil {
		return err
	}
	
	// 初始化admin用户
	if err := initAdminUser(); err != nil {
		return err
	}
	
	return nil
}

// initPermissions 初始化权限
func initPermissions() error {
	permissions := []model.Permission{
		// 数据概览模块
		{Name: "查看数据概览", Code: "dashboard:view", Module: "dashboard", Description: "查看系统数据统计和概览"},
		
		// 创作中心模块
		{Name: "内容生成", Code: "creation:generate", Module: "creation", Description: "使用AI生成内容"},
		{Name: "内容编辑", Code: "creation:edit", Module: "creation", Description: "编辑生成的内容"},
		{Name: "内容保存", Code: "creation:save", Module: "creation", Description: "保存内容"},
		
		// 内容管理模块
		{Name: "查看我的笔记", Code: "content:view", Module: "content", Description: "查看自己的内容列表"},
		{Name: "编辑内容", Code: "content:edit", Module: "content", Description: "编辑内容"},
		{Name: "删除内容", Code: "content:delete", Module: "content", Description: "删除内容"},
		{Name: "查看历史记录", Code: "content:history", Module: "content", Description: "查看内容历史版本"},
		{Name: "恢复历史版本", Code: "content:restore", Module: "content", Description: "恢复内容历史版本"},
		
		// 发布管理模块
		{Name: "发布测试", Code: "publish:test", Module: "publish", Description: "测试发布功能"},
		{Name: "立即发布", Code: "publish:now", Module: "publish", Description: "立即发布内容"},
		{Name: "定时发布", Code: "publish:schedule", Module: "publish", Description: "定时发布内容"},
		{Name: "查看发布历史", Code: "publish:history", Module: "publish", Description: "查看发布记录"},
		{Name: "取消发布", Code: "publish:cancel", Module: "publish", Description: "取消定时发布"},
		
		// 用户管理模块
		{Name: "查看用户列表", Code: "user:view", Module: "user", Description: "查看所有用户"},
		{Name: "编辑用户", Code: "user:edit", Module: "user", Description: "编辑用户信息"},
		{Name: "删除用户", Code: "user:delete", Module: "user", Description: "删除用户"},
		{Name: "设置用户角色", Code: "user:role", Module: "user", Description: "设置用户角色"},
		{Name: "启用/禁用用户", Code: "user:status", Module: "user", Description: "启用或禁用用户"},
		
		// 权限设置模块
		{Name: "查看角色列表", Code: "role:view", Module: "role", Description: "查看所有角色"},
		{Name: "创建角色", Code: "role:create", Module: "role", Description: "创建新角色"},
		{Name: "编辑角色", Code: "role:edit", Module: "role", Description: "编辑角色信息"},
		{Name: "删除角色", Code: "role:delete", Module: "role", Description: "删除角色"},
		{Name: "配置角色权限", Code: "role:permission", Module: "role", Description: "配置角色权限"},
		
		// 系统设置模块
		{Name: "查看系统设置", Code: "settings:view", Module: "settings", Description: "查看系统配置"},
		{Name: "修改系统设置", Code: "settings:edit", Module: "settings", Description: "修改系统配置"},
		
		// Token使用统计模块
		{Name: "查看Token使用统计", Code: "token:view", Module: "token", Description: "查看Token使用统计"},
		{Name: "查看全局Token统计", Code: "token:view_global", Module: "token", Description: "查看全局Token使用统计（管理员）"},
	}
	
	for _, p := range permissions {
		var count int64
		DB.Model(&model.Permission{}).Where("code = ?", p.Code).Count(&count)
		if count == 0 {
			if err := DB.Create(&p).Error; err != nil {
				return err
			}
		}
	}
	
	return nil
}

// initRoles 初始化角色
func initRoles() error {
	// 获取所有权限
	var allPermissions []model.Permission
	if err := DB.Find(&allPermissions).Error; err != nil {
		return err
	}
	
	allPermissionCodes := make([]string, len(allPermissions))
	for i, p := range allPermissions {
		allPermissionCodes[i] = p.Code
	}
	allPermissionsJSON, _ := json.Marshal(allPermissionCodes)
	
	// 普通用户权限
	userPermissions := []string{
		"dashboard:view",
		"creation:generate",
		"creation:edit",
		"creation:save",
		"content:view",
		"content:edit",
		"content:delete",
		"content:history",
		"content:restore",
		"token:view",
	}
	userPermissionsJSON, _ := json.Marshal(userPermissions)
	
	// 内容管理员权限
	contentManagerPermissions := append(userPermissions,
		"publish:test",
		"publish:now",
		"publish:schedule",
		"publish:history",
		"publish:cancel",
		"token:view",
	)
	contentManagerPermissionsJSON, _ := json.Marshal(contentManagerPermissions)
	
	roles := []model.Role{
		{
			Name:        "超级管理员",
			Code:        "super_admin",
			Description: "拥有系统所有权限",
			Permissions: string(allPermissionsJSON),
			IsSystem:    true,
		},
		{
			Name:        "内容管理员",
			Code:        "content_manager",
			Description: "负责内容的创作、编辑和发布",
			Permissions: string(contentManagerPermissionsJSON),
			IsSystem:    true,
		},
		{
			Name:        "普通用户",
			Code:        "user",
			Description: "普通用户权限，只能查看和创建自己的内容",
			Permissions: string(userPermissionsJSON),
			IsSystem:    true,
		},
	}
	
	for _, r := range roles {
		var count int64
		DB.Model(&model.Role{}).Where("code = ?", r.Code).Count(&count)
		if count == 0 {
			if err := DB.Create(&r).Error; err != nil {
				return err
			}
		}
	}
	
	return nil
}

// initAdminUser 初始化admin用户
func initAdminUser() error {
	var count int64
	DB.Model(&model.User{}).Where("username = ?", "admin").Count(&count)
	
	if count == 0 {
		// 获取超级管理员角色
		var adminRole model.Role
		if err := DB.Where("code = ?", "super_admin").First(&adminRole).Error; err != nil {
			return err
		}
		
		// 加密密码
		hashedPassword, err := utils.HashPassword("admin123")
		if err != nil {
			return err
		}
		
		// 创建admin用户
		adminUser := &model.User{
			Username: "admin",
			Email:    "admin@example.com",
			Password: hashedPassword,
			Nickname: "超级管理员",
			RoleID:   adminRole.ID,
			Status:   1,
		}
		
		if err := DB.Create(adminUser).Error; err != nil {
			return err
		}
		
		log.Println("Admin user created successfully")
	}
	
	return nil
}
