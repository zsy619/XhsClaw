// Package config 提供应用程序配置
package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config 应用程序配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	DeepSeek DeepSeekConfig `mapstructure:"deepseek"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Security SecurityConfig `mapstructure:"security"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int    `mapstructure:"port"`
	Mode         string `mapstructure:"mode"`
	ReadTimeout  int    `mapstructure:"read_timeout"`  // 读取超时（秒）
	WriteTimeout int    `mapstructure:"write_timeout"` // 写入超时（秒）
	IdleTimeout  int    `mapstructure:"idle_timeout"`  // 空闲超时（秒）
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type            string `mapstructure:"type"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	SSLMode         string `mapstructure:"sslmode"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`    // 最大打开连接数
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`    // 最大空闲连接数
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"` // 连接最大生命周期（秒）
	ConnMaxIdleTime int    `mapstructure:"conn_max_idle_time"` // 连接最大空闲时间（秒）
}

// DeepSeekConfig DeepSeek API配置
type DeepSeekConfig struct {
	APIKey     string `mapstructure:"api_key"`
	Model      string `mapstructure:"model"`
	BaseURL    string `mapstructure:"base_url"`
	Timeout    int    `mapstructure:"timeout"`     // 请求超时（秒）
	MaxRetries int    `mapstructure:"max_retries"` // 最大重试次数
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret           string `mapstructure:"secret"`
	Expire           int    `mapstructure:"expire"`            // 过期时间（小时）
	RefreshExpire    int    `mapstructure:"refresh_expire"`    // Refresh Token过期时间（天）
	Issuer           string `mapstructure:"issuer"`            // 签发者
	Audience         string `mapstructure:"audience"`          // 受众
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	CORS            CORSConfig            `mapstructure:"cors"`
	RateLimit       RateLimitConfig       `mapstructure:"rate_limit"`
	RequestSecurity RequestSecurityConfig `mapstructure:"request_security"`
}

// CORSConfig 跨域配置
type CORSConfig struct {
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	ExposeHeaders    []string `mapstructure:"expose_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"` // 预检请求缓存时间（秒）
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled bool    `mapstructure:"enabled"`
	RPS     float64 `mapstructure:"rps"`      // 每秒请求数
	Burst   int     `mapstructure:"burst"`     // 令牌桶容量
}

// RequestSecurityConfig 请求安全配置
type RequestSecurityConfig struct {
	MaxBodySize      int64  `mapstructure:"max_body_size"`       // 最大请求体大小（字节）
	EnableCSRF       bool   `mapstructure:"enable_csrf"`          // 启用CSRF保护
	EnableXSS        bool   `mapstructure:"enable_xss"`           // 启用XSS防护
	SecurityHeaders  bool   `mapstructure:"security_headers"`     // 启用安全头
	ContentSecurityPolicy string `mapstructure:"content_security_policy"` // 内容安全策略
}

var AppConfig *Config

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()
	
	// 设置配置文件类型
	v.SetConfigType("yaml")
	
	// 设置环境变量前缀
	v.SetEnvPrefix("XHS")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	
	// 设置默认值
	setDefaults(v)
	
	// 读取环境变量
	v.AutomaticEnv()
	
	// 尝试读取配置文件
	if configPath != "" {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			// 如果配置文件不存在，只使用环境变量和默认值
			fmt.Printf("Warning: Config file not found at %s, using defaults and environment variables\n", configPath)
		}
	} else {
		// 尝试从默认位置读取配置文件
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.AddConfigPath("/etc/xiaohongshu")
		v.SetConfigName("config")
		_ = v.ReadInConfig() // 忽略错误，继续使用默认值
	}
	
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	// 验证配置
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	
	AppConfig = &config
	return &config, nil
}

// setDefaults 设置默认配置值
func setDefaults(v *viper.Viper) {
	// 服务器配置
	v.SetDefault("server.port", 8000)
	v.SetDefault("server.mode", "debug")
	v.SetDefault("server.read_timeout", 60)
	v.SetDefault("server.write_timeout", 60)
	v.SetDefault("server.idle_timeout", 120)
	
	// 数据库配置
	v.SetDefault("database.type", "postgres")
	v.SetDefault("database.host", "127.0.0.1")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "123456")
	v.SetDefault("database.dbname", "xiaohongshu")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.max_open_conns", 100)
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.conn_max_lifetime", 3600)
	v.SetDefault("database.conn_max_idle_time", 600)
	
	// DeepSeek配置
	v.SetDefault("deepseek.model", "deepseek-chat")
	v.SetDefault("deepseek.base_url", "https://api.deepseek.com")
	v.SetDefault("deepseek.timeout", 120)
	v.SetDefault("deepseek.max_retries", 3)
	
	// JWT配置
	v.SetDefault("jwt.secret", "xiaohongshu-secret-key-change-in-production")
	v.SetDefault("jwt.expire", 24)
	v.SetDefault("jwt.refresh_expire", 7)
	v.SetDefault("jwt.issuer", "xiaohongshu-backend")
	v.SetDefault("jwt.audience", "xiaohongshu-frontend")
	
	// 安全配置 - CORS
	v.SetDefault("security.cors.allow_origins", []string{"http://localhost:5173", "http://localhost:3000"})
	v.SetDefault("security.cors.allow_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	v.SetDefault("security.cors.allow_headers", []string{"Origin", "Content-Type", "Authorization", "X-Requested-With"})
	v.SetDefault("security.cors.expose_headers", []string{"Content-Length", "Content-Type"})
	v.SetDefault("security.cors.allow_credentials", true)
	v.SetDefault("security.cors.max_age", 86400)
	
	// 安全配置 - 限流
	v.SetDefault("security.rate_limit.enabled", true)
	v.SetDefault("security.rate_limit.rps", 100.0)
	v.SetDefault("security.rate_limit.burst", 200)
	
	// 安全配置 - 请求安全
	v.SetDefault("security.request_security.max_body_size", 10485760) // 10MB
	v.SetDefault("security.request_security.enable_csrf", false) // 前端使用Bearer Token，不需要CSRF
	v.SetDefault("security.request_security.enable_xss", true)
	v.SetDefault("security.request_security.security_headers", true)
	v.SetDefault("security.request_security.content_security_policy", "default-src 'self'")
}

// validateConfig 验证配置的有效性
func validateConfig(config *Config) error {
	// 验证服务器端口
	if config.Server.Port < 1 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}
	
	// 验证数据库配置
	if config.Database.Type != "mysql" && config.Database.Type != "postgres" {
		return fmt.Errorf("unsupported database type: %s", config.Database.Type)
	}
	
	// 验证JWT密钥
	if config.JWT.Secret == "" || config.JWT.Secret == "xiaohongshu-secret-key-change-in-production" {
		if config.Server.Mode == "production" {
			return fmt.Errorf("JWT secret must be changed in production mode")
		}
	}
	
	return nil
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
	// 默认返回 PostgreSQL 格式，添加UTF8编码支持
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s client_encoding=utf8",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.DBName,
		c.SSLMode,
	)
}
