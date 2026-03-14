// Package utils 提供工具函数
package utils

import "encoding/json"

// ToJSON 将对象转为JSON字符串
func ToJSON(v interface{}) string {
	bytes, _ := json.Marshal(v)
	return string(bytes)
}

// FromJSON 将JSON字符串转为对象
func FromJSON(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}
