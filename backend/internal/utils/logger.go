// Package utils 提供日志工具
package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// LogLevel 日志级别
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// 日志级别字符串
var logLevelStrings = []string{
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
	"FATAL",
}

// Logger 日志记录器
type Logger struct {
	mu         sync.Mutex
	file       *os.File
	consoleLog *log.Logger
	fileLog    *log.Logger
	level      LogLevel
}

var (
	instance *Logger
	once     sync.Once
)

// GetLogger 获取单例日志记录器
func GetLogger() *Logger {
	once.Do(func() {
		instance = initLogger()
	})
	return instance
}

// initLogger 初始化日志记录器
func initLogger() *Logger {
	logger := &Logger{
		level: INFO,
	}

	// 控制台输出
	logger.consoleLog = log.New(os.Stdout, "", log.LstdFlags)

	// 创建日志目录
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		logger.consoleLog.Printf("Failed to create log directory: %v", err)
		return logger
	}

	// 打开日志文件
	logFile := filepath.Join(logDir, fmt.Sprintf("app_%s.log", time.Now().Format("2006-01-02")))
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logger.consoleLog.Printf("Failed to open log file: %v", err)
		return logger
	}

	logger.file = file
	logger.fileLog = log.New(file, "", log.LstdFlags)

	// 每天自动切换日志文件
	go logger.rotateDaily()

	return logger
}

// SetLogLevel 设置日志级别
func (l *Logger) SetLogLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// log 内部日志记录方法
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	levelStr := logLevelStrings[level]
	message := fmt.Sprintf(format, args...)
	logMessage := fmt.Sprintf("[%s] %s", levelStr, message)

	// 输出到控制台
	l.consoleLog.Println(logMessage)

	// 输出到文件
	if l.fileLog != nil {
		l.fileLog.Println(logMessage)
	}
}

// Debug 调试日志
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info 信息日志
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn 警告日志
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error 错误日志
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Fatal 致命错误日志
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
	os.Exit(1)
}

// rotateDaily 每日日志轮转
func (l *Logger) rotateDaily() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	lastDate := time.Now().Format("2006-01-02")

	for range ticker.C {
		currentDate := time.Now().Format("2006-01-02")
		if currentDate != lastDate {
			l.mu.Lock()
			if l.file != nil {
				l.file.Close()
			}

			logFile := filepath.Join("logs", fmt.Sprintf("app_%s.log", currentDate))
			file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				l.consoleLog.Printf("Failed to rotate log file: %v", err)
			} else {
				l.file = file
				l.fileLog = log.New(file, "", log.LstdFlags)
			}

			lastDate = currentDate
			l.mu.Unlock()
		}
	}
}

// Close 关闭日志文件
func (l *Logger) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		l.file.Close()
		l.file = nil
	}
}

// 便捷方法
func Debug(format string, args ...interface{}) {
	GetLogger().Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	GetLogger().Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	GetLogger().Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	GetLogger().Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	GetLogger().Fatal(format, args...)
}
