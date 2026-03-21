// Package repository 提供数据访问层
package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
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
			&model.Content{},
			&model.ContentHistory{},
			&model.PublishRecord{},
			&model.TokenBlacklist{},
			&model.TokenUsage{},
			&model.LLMProvider{},
			&model.XHSConfig{},
			&model.PublishConfig{},
			&model.SystemDict{},
		)
		if err != nil {
			initErr = err
			return
		}

		// // 确保 token_usage 表有 input_tokens 和 output_tokens 列（兼容旧数据库）
		// if cfg.Type == "postgres" {
		// 	DB.Exec("ALTER TABLE token_usage ADD COLUMN IF NOT EXISTS input_tokens INTEGER DEFAULT 0")
		// 	DB.Exec("ALTER TABLE token_usage ADD COLUMN IF NOT EXISTS output_tokens INTEGER DEFAULT 0")
		// } else if cfg.Type == "sqlite" {
		// 	// SQLite 需要使用单独的 ALTER TABLE 语句
		// 	var colCount int64
		// 	DB.Raw("SELECT COUNT(*) FROM pragma_table_info('token_usage') WHERE name = 'input_tokens'").Scan(&colCount)
		// 	if colCount == 0 {
		// 		DB.Exec("ALTER TABLE token_usage ADD COLUMN input_tokens INTEGER DEFAULT 0")
		// 	}
		// 	DB.Raw("SELECT COUNT(*) FROM pragma_table_info('token_usage') WHERE name = 'output_tokens'").Scan(&colCount)
		// 	if colCount == 0 {
		// 		DB.Exec("ALTER TABLE token_usage ADD COLUMN output_tokens INTEGER DEFAULT 0")
		// 	}
		// }

		// 修复旧用户数据的 role_id（设置为默认的普通用户角色ID=3）
		var count int64
		DB.Model(&model.User{}).Where("role_id IS NULL OR role_id = 0").Count(&count)
		if count > 0 {
			log.Printf("Fixing %d users with invalid role_id...", count)
			DB.Model(&model.User{}).Where("role_id IS NULL OR role_id = 0").Update("role_id", 3)
		}

		// 自动添加表和字段注释
		if err := addTableComments(cfg.Type); err != nil {
			log.Printf("Warning: Failed to add table comments: %v", err)
		}

		log.Println("Database initialized successfully")

		// 初始化admin用户
		if err := initAdminUser(); err != nil {
			log.Printf("Failed to init admin user: %v", err)
		}

		// 初始化系统字典数据
		if err := initSystemDicts(); err != nil {
			log.Printf("Failed to init system dicts: %v", err)
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

		// 大模型配置模块
		{Name: "查看大模型配置", Code: "llm:view", Module: "llm", Description: "查看大模型服务商配置"},
		{Name: "创建大模型配置", Code: "llm:create", Module: "llm", Description: "创建新的大模型服务商配置"},
		{Name: "编辑大模型配置", Code: "llm:edit", Module: "llm", Description: "编辑大模型服务商配置"},
		{Name: "删除大模型配置", Code: "llm:delete", Module: "llm", Description: "删除大模型服务商配置"},

		// 小红书配置模块
		{Name: "查看小红书配置", Code: "xhs:view", Module: "xhs", Description: "查看小红书账号配置"},
		{Name: "创建小红书配置", Code: "xhs:create", Module: "xhs", Description: "创建新的小红书账号配置"},
		{Name: "编辑小红书配置", Code: "xhs:edit", Module: "xhs", Description: "编辑小红书账号配置"},
		{Name: "删除小红书配置", Code: "xhs:delete", Module: "xhs", Description: "删除小红书账号配置"},
		{Name: "验证小红书配置", Code: "xhs:verify", Module: "xhs", Description: "验证小红书Cookie有效性"},

		// 系统字典模块
		{Name: "查看系统字典", Code: "dict:view", Module: "dict", Description: "查看系统字典"},
		{Name: "管理系统字典", Code: "dict:manage", Module: "dict", Description: "管理系统字典（管理员）"},
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
			Email:    "admin@yy24365.com",
			Password: hashedPassword,
			Nickname: "超级管理员",
			RoleID:   adminRole.ID,
			Status:   1,
		}

		if err := DB.Create(adminUser).Error; err != nil {
			return err
		}

		log.Println("Admin user created created successfully")
	}

	return nil
}

// addTableComments 自动为所有表和字段添加注释
func addTableComments(dbType string) error {
	log.Println("Adding table and column comments...")

	// 字段注释映射 (表名 -> 字段名 -> 注释)
	fieldComments := map[string]map[string]string{
		"users": {
			"id":       "用户ID",
			"username": "用户名",
			"email":    "邮箱",
			"password": "密码",
			"nickname": "昵称",
			"avatar":   "头像URL",
			"role_id":  "角色ID",
			"status":   "状态 1:正常 0:禁用",
		},
		"token_blacklists": {
			"id":         "记录ID",
			"token":      "Token值",
			"user_id":    "用户ID",
			"expires_at": "过期时间",
		},
		"roles": {
			"id":          "角色ID",
			"name":        "角色名称",
			"code":        "角色代码",
			"description": "角色描述",
			"permissions": "权限列表JSON",
			"is_system":   "是否系统角色",
		},
		"permissions": {
			"id":          "权限ID",
			"name":        "权限名称",
			"code":        "权限代码",
			"module":      "所属模块",
			"description": "权限描述",
		},
		"contents": {
			"id":                   "内容ID",
			"user_id":              "用户ID",
			"title":                "标题",
			"title_options":        "备选标题JSON",
			"selected_title_index": "选中标题索引",
			"description":          "正文内容",
			"tags":                 "标签JSON数组",
			"images":               "图片URL数组JSON",
			"cover_suggestion":     "封面建议文案",
			"content_attributes":   "内容属性JSON",
			"render_attributes":    "渲染属性JSON",
			"status":               "状态 0:草稿 1:待发布 2:已发布 3:失败",
		},
		"content_histories": {
			"id":                   "历史记录ID",
			"content_id":           "内容ID",
			"user_id":              "用户ID",
			"type":                 "操作类型 create/edit/delete/publish",
			"title":                "标题",
			"title_options":        "备选标题JSON",
			"selected_title_index": "选中标题索引",
			"description":          "正文内容",
			"tags":                 "标签JSON数组",
			"images":               "图片URL数组JSON",
			"cover_suggestion":     "封面建议文案",
			"content_attributes":   "内容属性JSON",
			"render_attributes":    "渲染属性JSON",
			"change_reason":        "变更原因",
		},
		"user_configs": {
			"id":                   "配置ID",
			"user_id":              "用户ID",
			"llm_api_key":          "大模型API密钥",
			"llm_base_url":         "大模型API地址",
			"llm_model":            "大模型名称",
			"xiaohongshu_cookie":   "小红书Cookie",
			"xiaohongshu_user_id":  "小红书用户ID",
			"xiaohongshu_token":    "小红书Token",
			"default_publish_time": "默认发布时间 HH:mm",
			"auto_publish_enabled": "是否启用自动发布",
		},
		"publish_records": {
			"id":           "发布记录ID",
			"user_id":      "用户ID",
			"content_id":   "内容ID",
			"status":       "状态 0:待发布 1:发布中 2:成功 3:失败",
			"error_msg":    "错误信息",
			"scheduled_at": "计划发布时间",
			"published_at": "实际发布时间",
		},
		"token_usage": {
			"id":              "记录ID",
			"user_id":         "用户ID",
			"model":           "使用的模型",
			"provider":        "提供商",
			"input_tokens":    "输入tokens",
			"output_tokens":   "输出tokens",
			"total_tokens":    "总tokens",
			"cost":            "费用(美元)",
			"request_type":    "请求类型",
			"request_content": "请求内容摘要",
			"response_status": "响应状态",
			"error_message":   "错误信息",
			"ip_address":      "IP地址",
			"user_agent":      "用户代理",
		},
	}

	// 根据数据库类型执行不同的SQL
	if dbType == "mysql" {
		return addMySQLComments(fieldComments)
	} else {
		return addPostgreSQLComments(fieldComments)
	}
}

// addMySQLComments 为MySQL添加表和字段注释
func addMySQLComments(fieldComments map[string]map[string]string) error {
	for tableName, fields := range fieldComments {
		// 添加表注释
		if err := DB.Exec(fmt.Sprintf("ALTER TABLE `%s` COMMENT = ?", tableName), getTableComment(tableName)).Error; err != nil {
			log.Printf("Failed to add comment for table %s: %v", tableName, err)
			continue
		}

		// 添加字段注释
		for fieldName, fieldComment := range fields {
			sql := fmt.Sprintf("ALTER TABLE `%s` MODIFY COLUMN `%s` VARCHAR(255) COMMENT ?", tableName, fieldName)
			if err := DB.Exec(sql, fieldComment).Error; err != nil {
				log.Printf("Failed to add comment for column %s.%s: %v", tableName, fieldName, err)
			}
		}
		log.Printf("Added comments for table: %s", tableName)
	}
	return nil
}

// addPostgreSQLComments 为PostgreSQL添加表和字段注释
func addPostgreSQLComments(fieldComments map[string]map[string]string) error {
	for tableName, fields := range fieldComments {
		// 添加表注释（PostgreSQL 不支持参数化，需要直接拼接字符串）
		tableComment := getTableComment(tableName)
		sql := fmt.Sprintf("COMMENT ON TABLE %s IS '%s'", tableName, escapeSQLString(tableComment))
		if err := DB.Exec(sql).Error; err != nil {
			log.Printf("Failed to add comment for table %s: %v", tableName, err)
			continue
		}

		// 添加字段注释
		for fieldName, fieldComment := range fields {
			sql := fmt.Sprintf("COMMENT ON COLUMN %s.%s IS '%s'", tableName, fieldName, escapeSQLString(fieldComment))
			if err := DB.Exec(sql).Error; err != nil {
				log.Printf("Failed to add comment for column %s.%s: %v", tableName, fieldName, err)
			}
		}
		log.Printf("Added comments for table: %s", tableName)
	}
	return nil
}

// escapeSQLString 转义 SQL 字符串中的特殊字符
func escapeSQLString(s string) string {
	// 转义单引号（PostgreSQL 中用两个单引号表示一个单引号）
	s = strings.ReplaceAll(s, "'", "''")
	// 转义反斜杠
	s = strings.ReplaceAll(s, "\\", "\\\\")
	return s
}

// getTableComment 根据表名获取表注释
func getTableComment(tableName string) string {
	tableComments := map[string]string{
		"users":               "用户表",
		"token_blacklists":    "Token黑名单表",
		"roles":               "角色表",
		"permissions":         "权限表",
		"contents":            "小红书内容表",
		"content_histories":   "内容历史记录表",
		"user_configs":        "用户配置表",
		"publish_records":     "发布记录表",
		"token_usage":         "Token使用记录表",
		"llm_providers":       "大模型服务商配置表",
		"xiaohongshu_configs": "小红书账号配置表",
		"publish_configs":     "发布配置表",
		"system_dicts":        "系统字典表",
	}
	if comment, ok := tableComments[tableName]; ok {
		return comment
	}
	return tableName
}

// initSystemDicts 初始化系统字典数据
func initSystemDicts() error {
	// 检查是否已初始化
	var count int64
	DB.Model(&model.SystemDict{}).Where("category = ?", model.DictCategoryLLMProvider).Count(&count)
	if count > 0 {
		log.Printf("System dicts already initialized, skipping...")
		return nil
	}

	// 从 model 常量动态生成大模型服务商字典
	providers := generateLLMProviderDicts()

	// 小红书配置状态
	configStatuses := []model.SystemDict{
		{Category: "xhs_status", Code: "pending", Name: "待验证", Value: "pending", Description: "配置待验证状态", SortOrder: 1, Enabled: true},
		{Category: "xhs_status", Code: "active", Name: "正常", Value: "active", Description: "配置正常使用状态", SortOrder: 2, Enabled: true},
		{Category: "xhs_status", Code: "expired", Name: "已过期", Value: "expired", Description: "Cookie或Token已过期", SortOrder: 3, Enabled: true},
		{Category: "xhs_status", Code: "error", Name: "异常", Value: "error", Description: "配置异常状态", SortOrder: 4, Enabled: true},
	}

	// 插入数据
	for _, p := range providers {
		if err := DB.FirstOrCreate(&p, model.SystemDict{Category: p.Category, Code: p.Code}).Error; err != nil {
			log.Printf("Failed to create system dict %s: %v", p.Code, err)
		}
	}

	for _, s := range configStatuses {
		if err := DB.FirstOrCreate(&s, model.SystemDict{Category: s.Category, Code: s.Code}).Error; err != nil {
			log.Printf("Failed to create system dict %s: %v", s.Code, err)
		}
	}

	log.Printf("System dicts initialized successfully")
	return nil
}

// generateLLMProviderDicts 从 model.LLMProvider 常量生成系统字典
func generateLLMProviderDicts() []model.SystemDict {
	var dicts []model.SystemDict
	sortOrder := 1

	// 按优先级分组排序
	priorityOrder := []string{
		model.ProviderOpenAI,
		model.ProviderDeepSeek,
		model.ProviderAzure,
		model.ProviderClaude,
		model.ProviderGemini,
		model.ProviderQwen,
		model.ProviderGLM,
		model.ProviderQianfan,
		model.ProviderHunyuan,
		model.ProviderDoubao,
		model.ProviderSpark,
		model.ProviderBaichuan,
		model.ProviderMiniMax,
		model.ProviderMistral,
		model.ProviderCohere,
		model.ProviderGroq,
		model.ProviderReplicate,
		model.ProviderPerplexity,
		model.ProviderOllama,
		model.ProviderLMStudio,
		model.ProviderLocalAI,
		model.ProvidervLLM,
		model.ProviderCustom,
	}

	// 生成字典
	for _, provider := range priorityOrder {
		name := model.GetProviderDisplayName(provider)
		baseURL := model.GetProviderBaseURL(provider)
		description := model.GetProviderDescription(provider)
		models := model.GetProviderModels(provider)

		// 构建 Extra JSON
		extra := fmt.Sprintf(`{"models": %s}`, formatModelsJSON(models))

		dict := model.SystemDict{
			Category:    model.DictCategoryLLMProvider,
			Code:        provider,
			Name:        name,
			Value:       baseURL,
			Description: description,
			SortOrder:   sortOrder,
			Enabled:     true,
			Extra:       extra,
		}
		dicts = append(dicts, dict)
		sortOrder++
	}

	return dicts
}

// formatModelsJSON 格式化模型列表为 JSON 数组字符串
func formatModelsJSON(models []string) string {
	if len(models) == 0 {
		return "[]"
	}

	result := "["
	for i, m := range models {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf(`"%s"`, m)
	}
	result += "]"
	return result
}
