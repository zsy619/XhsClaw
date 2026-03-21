// Package repository 提供数据访问层
package repository

import (
	"log"
	"time"

	"gorm.io/gorm"

	"xiaohongshu/internal/model"
)

// MigrateUserConfigsToSeparateTables 将 user_configs 表数据迁移到独立表中
// 迁移策略：
// 1. user_configs.llm_api_key/base_url/model -> llm_providers 表
// 2. user_configs.xiaohongshu_cookie/user_id/token -> xiaohongshu_configs 表
// 3. user_configs.default_publish_time/auto_publish_enabled -> publish_configs 表
// 4. 删除 user_configs 表
func MigrateUserConfigsToSeparateTables() error {
	if DB == nil {
		log.Println("Database not initialized, skipping migration")
		return nil
	}

	// 检查 user_configs 表是否存在
	var tableCount int64
	if err := DB.Raw(`
		SELECT COUNT(*) FROM information_schema.tables 
		WHERE table_schema = DATABASE() AND table_name = 'user_configs'
	`).Scan(&tableCount).Error; err != nil {
		// 如果查询失败，尝试使用原生 SQL 查询
		log.Printf("Cannot check user_configs table existence: %v", err)
		return nil
	}

	if tableCount == 0 {
		log.Println("user_configs table does not exist, migration not needed")
		return nil
	}

	// 检查是否已有迁移记录（防止重复执行）
	var migrationLog model.SystemDict
	err := DB.Where("category = ? AND code = ?", "migration", "user_configs_migrated").First(&migrationLog).Error
	if err == nil {
		log.Println("user_configs migration already completed, skipping")
		return nil
	}

	log.Println("Starting user_configs migration...")

	// 开始事务
	tx := DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 确保在错误时回滚
	success := false
	defer func() {
		if !success {
			tx.Rollback()
		}
	}()

	// 1. 迁移 LLM 配置到 llm_providers 表
	if err := migrateLLMConfig(tx); err != nil {
		log.Printf("Failed to migrate LLM config: %v", err)
		return err
	}

	// 2. 迁移小红书配置到 xiaohongshu_configs 表
	if err := migrateXiaohongshuConfig(tx); err != nil {
		log.Printf("Failed to migrate Xiaohongshu config: %v", err)
		return err
	}

	// 3. 迁移发布配置到 publish_configs 表
	if err := migratePublishConfig(tx); err != nil {
		log.Printf("Failed to migrate Publish config: %v", err)
		return err
	}

	// 4. 删除 user_configs 表
	if err := dropUserConfigsTable(tx); err != nil {
		log.Printf("Failed to drop user_configs table: %v", err)
		return err
	}

	// 5. 记录迁移完成
	migrationRecord := model.SystemDict{
		Category:    "migration",
		Code:       "user_configs_migrated",
		Name:       "用户配置迁移",
		Value:      "true",
		Description: "user_configs 表已迁移到独立表",
		SortOrder:  1,
		Enabled:    true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if err := tx.Create(&migrationRecord).Error; err != nil {
		log.Printf("Failed to record migration: %v", err)
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return err
	}

	success = true
	log.Println("user_configs migration completed successfully")
	return nil
}

// migrateLLMConfig 迁移 LLM 配置
func migrateLLMConfig(tx *gorm.DB) error {
	// 查询所有有 LLM 配置的 user_configs 记录
	type UserConfig struct {
		ID         uint
		UserID     uint
		LLMAPIKey  string
		LLMBaseURL string
		LLMModel   string
	}

	var configs []UserConfig
	if err := tx.Table("user_configs").
		Where("llm_api_key IS NOT NULL AND llm_api_key != ''").
		Find(&configs).Error; err != nil {
		return err
	}

	log.Printf("Found %d user configs with LLM data to migrate", len(configs))

	for _, config := range configs {
		// 检查是否已存在对应的 LLM Provider
		var existingCount int64
		tx.Model(&model.LLMProvider{}).
			Where("user_id = ? AND is_default = ?", config.UserID, true).
			Count(&existingCount)

		// 如果已存在默认的 LLM Provider，跳过或更新
		if existingCount > 0 {
			// 更新现有的默认 LLM Provider
			err := tx.Model(&model.LLMProvider{}).
				Where("user_id = ? AND is_default = ?", config.UserID, true).
				Updates(map[string]interface{}{
					"api_key":     config.LLMAPIKey,
					"base_url":    config.LLMBaseURL,
					"model_name":  config.LLMModel,
					"updated_at": time.Now(),
				}).Error
			if err != nil {
				log.Printf("Failed to update LLM provider for user %d: %v", config.UserID, err)
				continue
			}
		} else {
			// 创建新的 LLM Provider
			provider := model.LLMProvider{
				UserID:     config.UserID,
				Name:       "默认配置",
				Provider:   model.ProviderCustom, // 使用 custom 类型，因为无法确定原始提供商
				APIKey:     config.LLMAPIKey,
				BaseURL:    config.LLMBaseURL,
				ModelName:  config.LLMModel,
				IsDefault:  true,
				IsEnabled:  true,
				Timeout:    60,
				RetryCount: 3,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}
			if err := tx.Create(&provider).Error; err != nil {
				log.Printf("Failed to create LLM provider for user %d: %v", config.UserID, err)
				continue
			}
		}

		log.Printf("Migrated LLM config for user %d", config.UserID)
	}

	return nil
}

// migrateXiaohongshuConfig 迁移小红书配置
func migrateXiaohongshuConfig(tx *gorm.DB) error {
	// 查询所有有小红书配置的 user_configs 记录
	type UserConfig struct {
		ID                  uint
		UserID              uint
		XiaohongshuCookie   string
		XiaohongshuUserId   string
		XiaohongshuToken    string
	}

	var configs []UserConfig
	if err := tx.Table("user_configs").
		Where("xiaohongshu_cookie IS NOT NULL AND xiaohongshu_cookie != ''").
		Find(&configs).Error; err != nil {
		return err
	}

	log.Printf("Found %d user configs with Xiaohongshu data to migrate", len(configs))

	for _, config := range configs {
		// 检查是否已存在对应的小红书配置
		var existingCount int64
		tx.Model(&model.XHSConfig{}).
			Where("user_id = ? AND is_default = ?", config.UserID, true).
			Count(&existingCount)

		if existingCount > 0 {
			// 更新现有的默认配置
			err := tx.Model(&model.XHSConfig{}).
				Where("user_id = ? AND is_default = ?", config.UserID, true).
				Updates(map[string]interface{}{
					"cookie":     config.XiaohongshuCookie,
					"xhs_user_id": config.XiaohongshuUserId,
					"token":      config.XiaohongshuToken,
					"updated_at": time.Now(),
				}).Error
			if err != nil {
				log.Printf("Failed to update Xiaohongshu config for user %d: %v", config.UserID, err)
				continue
			}
		} else {
			// 创建新的小红书配置
			xhsConfig := model.XHSConfig{
				UserID:      config.UserID,
				Name:        "默认账号",
				Cookie:      config.XiaohongshuCookie,
				XHSUserID:   config.XiaohongshuUserId,
				Token:       config.XiaohongshuToken,
				IsDefault:   true,
				IsEnabled:   true,
				Status:      "pending", // 需要重新验证
				Description: "从 user_configs 迁移",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			if err := tx.Create(&xhsConfig).Error; err != nil {
				log.Printf("Failed to create Xiaohongshu config for user %d: %v", config.UserID, err)
				continue
			}
		}

		log.Printf("Migrated Xiaohongshu config for user %d", config.UserID)
	}

	return nil
}

// migratePublishConfig 迁移发布配置
func migratePublishConfig(tx *gorm.DB) error {
	// 查询所有有发布配置的 user_configs 记录
	type UserConfig struct {
		ID                   uint
		UserID               uint
		DefaultPublishTime   string
		AutoPublishEnabled   bool
	}

	var configs []UserConfig
	if err := tx.Table("user_configs").
		Where("default_publish_time IS NOT NULL OR auto_publish_enabled = ?", true).
		Find(&configs).Error; err != nil {
		return err
	}

	log.Printf("Found %d user configs with Publish data to migrate", len(configs))

	for _, config := range configs {
		// 检查是否已存在对应的发布配置
		var existingCount int64
		tx.Model(&model.PublishConfig{}).
			Where("user_id = ?", config.UserID).
			Count(&existingCount)

		if existingCount > 0 {
			// 更新现有配置
			err := tx.Model(&model.PublishConfig{}).
				Where("user_id = ?", config.UserID).
				Updates(map[string]interface{}{
					"default_publish_time": config.DefaultPublishTime,
					"auto_publish_enabled": config.AutoPublishEnabled,
					"updated_at": time.Now(),
				}).Error
			if err != nil {
				log.Printf("Failed to update Publish config for user %d: %v", config.UserID, err)
				continue
			}
		} else {
			// 创建新的发布配置
			publishConfig := model.PublishConfig{
				UserID:               config.UserID,
				DefaultPublishTime:   config.DefaultPublishTime,
				AutoPublishEnabled:   config.AutoPublishEnabled,
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			}
			if err := tx.Create(&publishConfig).Error; err != nil {
				log.Printf("Failed to create Publish config for user %d: %v", config.UserID, err)
				continue
			}
		}

		log.Printf("Migrated Publish config for user %d", config.UserID)
	}

	return nil
}

// dropUserConfigsTable 删除 user_configs 表
func dropUserConfigsTable(tx *gorm.DB) error {
	// 执行删除表操作
	sql := "DROP TABLE IF EXISTS user_configs"
	if err := tx.Exec(sql).Error; err != nil {
		log.Printf("Failed to drop user_configs table: %v", err)
		return err
	}

	log.Println("Successfully dropped user_configs table")
	return nil
}
