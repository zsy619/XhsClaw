// Package utils 提供安全相关工具函数
package utils

import (
	"html"
	"regexp"
	"strings"
)

var (
	// 危险的HTML标签正则
	dangerousTagsRegex = regexp.MustCompile(`(?i)<(script|iframe|object|embed|form|input|button|link|meta|style|svg|img)[^>]*>`)
	// 危险的JavaScript事件正则
	dangerousEventsRegex = regexp.MustCompile(`(?i)\s(on\w+)\s*=`)
	// 危险的协议正则
	dangerousProtocolsRegex = regexp.MustCompile(`(?i)(javascript|vbscript|data):`)
)

// SanitizeHTML 清理HTML，防止XSS攻击
func SanitizeHTML(input string) string {
	if input == "" {
		return input
	}

	// 1. 移除危险标签
	result := dangerousTagsRegex.ReplaceAllString(input, "")

	// 2. 移除危险事件处理器
	result = dangerousEventsRegex.ReplaceAllString(result, " ")

	// 3. 移除危险协议
	result = dangerousProtocolsRegex.ReplaceAllString(result, "#")

	// 4. HTML实体编码剩余的特殊字符
	result = html.EscapeString(result)

	return result
}

// SanitizeString 清理普通字符串
func SanitizeString(input string) string {
	if input == "" {
		return input
	}

	// 移除首尾空白
	input = strings.TrimSpace(input)

	// HTML实体编码
	return html.EscapeString(input)
}

// ValidateInputLength 验证输入长度
func ValidateInputLength(input string, min, max int) bool {
	length := len(input)
	return length >= min && length <= max
}

// ValidateUsername 验证用户名格式
func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	// 只允许字母、数字、下划线
	match, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username)
	return match
}

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) bool {
	match, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
	return match
}

// ValidatePasswordStrength 验证密码强度
func ValidatePasswordStrength(password string) (bool, string) {
	if len(password) < 8 {
		return false, "密码长度至少8位"
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)

	if !hasUpper {
		return false, "密码需要包含大写字母"
	}
	if !hasLower {
		return false, "密码需要包含小写字母"
	}
	if !hasNumber {
		return false, "密码需要包含数字"
	}
	if !hasSpecial {
		return false, "密码需要包含特殊字符"
	}

	return true, ""
}
