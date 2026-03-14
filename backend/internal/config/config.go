// Package config 提供应用程序配置
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config 应用程序配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	DeepSeek DeepSeekConfig `mapstructure:"deepseek"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type     string `mapstructure:"type"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// DeepSeekConfig DeepSeek API配置
type DeepSeekConfig struct {
	APIKey string `mapstructure:"api_key"`
	Model  string `mapstructure:"model"`
	BaseURL string `mapstructure:"base_url"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"` // 过期时间（小时）
}

var AppConfig *Config

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()
	
	// 设置默认值
	v.SetDefault("server.port", 8000)
	v.SetDefault("server.mode", "debug")
	v.SetDefault("database.type", "postgres")
	v.SetDefault("database.host", "127.0.0.1")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "123456")
	v.SetDefault("database.dbname", "xiaohongshu")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("deepseek.model", "deepseek-chat")
	v.SetDefault("deepseek.base_url", "https://api.deepseek.com")
	v.SetDefault("jwt.secret", "xiaohongshu-secret-key")
	v.SetDefault("jwt.expire", 24)
	
	// 读取环境变量
	v.AutomaticEnv()
	
	// 尝试读取配置文件
	if configPath != "" {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}
	
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	AppConfig = &config
	return &config, nil
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	if c.Type == "mysql" {
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.User,
			c.Password,
			c.Host,
			c.Port,
			c.DBName,
		)
	}
	// 默认返回 PostgreSQL 格式
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.DBName,
		c.SSLMode,
	)
}
