// Package repository 提供数据访问层
package repository

import (
	"fmt"
	"log"
	"time"
)

// BackupUserConfigs 备份 user_configs 表数据
// 返回备份记录的 JSON 格式字符串，以便在需要时恢复
func BackupUserConfigs() (string, error) {
	if DB == nil {
		return "", fmt.Errorf("database not initialized")
	}

	// 查询所有 user_configs 数据
	type UserConfigBackup struct {
		ID                   uint      `json:"id"`
		UserID               uint      `json:"user_id"`
		LLMAPIKey            string    `json:"llm_api_key"`
		LLMBaseURL           string    `json:"llm_base_url"`
		LLMModel             string    `json:"llm_model"`
		XiaohongshuCookie     string    `json:"xiaohongshu_cookie"`
		XiaohongshuUserId     string    `json:"xiaohongshu_user_id"`
		XiaohongshuToken      string    `json:"xiaohongshu_token"`
		DefaultPublishTime    string    `json:"default_publish_time"`
		AutoPublishEnabled    bool      `json:"auto_publish_enabled"`
		CreatedAt             time.Time `json:"created_at"`
		UpdatedAt             time.Time `json:"updated_at"`
	}

	var configs []UserConfigBackup
	if err := DB.Table("user_configs").Find(&configs).Error; err != nil {
		return "", fmt.Errorf("failed to query user_configs: %w", err)
	}

	if len(configs) == 0 {
		log.Println("No user_configs data to backup")
		return "[]", nil
	}

	log.Printf("Found %d user_configs records to backup", len(configs))

	// 创建备份表
	backupTableName := fmt.Sprintf("user_configs_backup_%d", time.Now().Unix())
	
	// 创建备份表并复制数据
	if err := DB.Exec(fmt.Sprintf(`
		CREATE TABLE %s AS SELECT * FROM user_configs
	`, backupTableName)).Error; err != nil {
		return "", fmt.Errorf("failed to create backup table: %w", err)
	}

	log.Printf("Backup table %s created successfully", backupTableName)

	// 为了返回备份数据，我们将其序列化为 JSON
	// 实际上在生产环境中，应该保存到文件或专门的备份存储中
	// 这里我们记录到日志中
	for _, config := range configs {
		log.Printf("Backup record: UserID=%d, HasLLM=%v, HasXHS=%v, HasPublish=%v",
			config.UserID,
			config.LLMAPIKey != "",
			config.XiaohongshuCookie != "",
			config.DefaultPublishTime != "" || config.AutoPublishEnabled,
		)
	}

	// 返回备份表名
	return backupTableName, nil
}

// RestoreUserConfigsFromBackup 从备份表恢复数据
func RestoreUserConfigsFromBackup(backupTableName string) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	// 检查备份表是否存在
	var count int64
	if err := DB.Raw(fmt.Sprintf(`
		SELECT COUNT(*) FROM information_schema.tables 
		WHERE table_schema = DATABASE() AND table_name = '%s'
	`, backupTableName)).Scan(&count).Error; err != nil {
		return fmt.Errorf("failed to check backup table: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("backup table %s does not exist", backupTableName)
	}

	// 重新创建 user_configs 表
	if err := DB.Exec("DROP TABLE IF EXISTS user_configs").Error; err != nil {
		return fmt.Errorf("failed to drop user_configs table: %w", err)
	}

	// 从备份表恢复
	if err := DB.Exec(fmt.Sprintf(`
		CREATE TABLE user_configs AS SELECT * FROM %s
	`, backupTableName)).Error; err != nil {
		return fmt.Errorf("failed to restore from backup: %w", err)
	}

	log.Printf("Successfully restored user_configs from backup table %s", backupTableName)
	return nil
}

// CleanBackupTables 清理过期的备份表
func CleanBackupTables() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	// 查找所有备份表
	var tables []string
	if err := DB.Raw(`
		SELECT table_name FROM information_schema.tables 
		WHERE table_schema = DATABASE() AND table_name LIKE 'user_configs_backup_%'
	`).Pluck("table_name", &tables).Error; err != nil {
		return fmt.Errorf("failed to query backup tables: %w", err)
	}

	// 删除超过 7 天的备份表
	cutoffDate := time.Now().AddDate(0, 0, -7)
	for _, table := range tables {
		// 从表名中提取时间戳
		var timestamp int64
		_, err := fmt.Sscanf(table, "user_configs_backup_%d", &timestamp)
		if err != nil {
			log.Printf("Failed to parse timestamp from table %s: %v", table, err)
			continue
		}

		tableDate := time.Unix(timestamp, 0)
		if tableDate.Before(cutoffDate) {
			if err := DB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table)).Error; err != nil {
				log.Printf("Failed to drop backup table %s: %v", table, err)
				continue
			}
			log.Printf("Dropped old backup table: %s", table)
		}
	}

	return nil
}
