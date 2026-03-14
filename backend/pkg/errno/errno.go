// Package errno 定义错误码
package errno

// ErrNo 错误码结构
type ErrNo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// 通用错误码
var (
	Success          = NewErrNo(0, "成功")
	InvalidParams    = NewErrNo(10001, "参数错误")
	Unauthorized     = NewErrNo(10002, "未授权")
	Forbidden        = NewErrNo(10003, "禁止访问")
	NotFound         = NewErrNo(10004, "资源不存在")
	InternalError    = NewErrNo(10005, "服务器内部错误")
	ServiceUnavailable = NewErrNo(10006, "服务暂不可用")
)

// 用户相关错误码
var (
	UserAlreadyExists = NewErrNo(20001, "用户已存在")
	UserNotFound      = NewErrNo(20002, "用户不存在")
	WrongPassword     = NewErrNo(20003, "密码错误")
	UserDisabled      = NewErrNo(20004, "用户已被禁用")
)

// 内容相关错误码
var (
	ContentNotFound = NewErrNo(30001, "内容不存在")
	GenerateFailed  = NewErrNo(30002, "内容生成失败")
	PublishFailed   = NewErrNo(30003, "发布失败")
)

// 角色和权限相关错误码
var (
	RoleNotFound         = NewErrNo(40001, "角色不存在")
	RoleAlreadyExists    = NewErrNo(40002, "角色已存在")
	RoleInUse            = NewErrNo(40003, "角色正在被使用，无法删除")
	CannotDeleteSystemRole = NewErrNo(40004, "系统角色不可删除")
)

// NewErrNo 创建新的错误码
func NewErrNo(code int, message string) *ErrNo {
	return &ErrNo{
		Code:    code,
		Message: message,
	}
}

// Error 实现error接口
func (e *ErrNo) Error() string {
	return e.Message
}

// WithMessage 自定义错误信息
func (e *ErrNo) WithMessage(message string) *ErrNo {
	return &ErrNo{
		Code:    e.Code,
		Message: message,
	}
}
