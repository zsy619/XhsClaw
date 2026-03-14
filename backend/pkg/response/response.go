// Package response 统一响应处理
package response

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"xiaohongshu/pkg/errno"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *app.RequestContext, data interface{}) {
	c.JSON(consts.StatusOK, Response{
		Code:    errno.Success.Code,
		Message: errno.Success.Message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *app.RequestContext, err *errno.ErrNo) {
	c.JSON(consts.StatusOK, Response{
		Code:    err.Code,
		Message: err.Message,
	})
}

// ErrorWithMessage 带自定义消息的错误响应
func ErrorWithMessage(c *app.RequestContext, err *errno.ErrNo, message string) {
	c.JSON(consts.StatusOK, Response{
		Code:    err.Code,
		Message: message,
	})
}

// ParamError 参数错误响应
func ParamError(c *app.RequestContext, message string) {
	ErrorWithMessage(c, errno.InvalidParams, message)
}
