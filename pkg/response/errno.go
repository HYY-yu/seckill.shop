package response

import (
	"encoding/json"
	"net/http"
)

var _ Error = (*err)(nil)

// Error 封装返回错误（HTTP状态码+自定状态码+返回错误信息堆栈）
type Error interface {
	error
	// WithErr 设置错误信息
	WithErr(err error) Error
	// GetBusinessCode 获取 Business Code
	GetBusinessCode() int
	// GetHttpCode 获取 HTTP Code
	GetHttpCode() int
	// GetMsg 获取 Msg
	GetMsg() string
	// GetErr 获取错误信息
	GetErr() error
	// ToString 返回 JSON 格式的错误详情
	ToString() string
}

type err struct {
	HttpCode     int    // HTTP Code
	BusinessCode int    // Business Code
	Message      string // 描述信息
	Err          error  // 错误信息
}

// NewError 新建一个 Error
func NewError(httpCode, businessCode int, msg string) Error {
	return &err{
		HttpCode:     httpCode,
		BusinessCode: businessCode,
		Message:      msg,
	}
}

// NewErrorWithStatusOk 新建 Error，httpCode 默认为 http.StatusOK
func NewErrorWithStatusOk(businessCode int, msg string) Error {
	return &err{
		HttpCode:     http.StatusOK,
		BusinessCode: businessCode,
		Message:      msg,
	}
}

// NewErrorWithStatusOkAutoMsg 新建Error，httpCode 默认为 http.StatusOK，msg自动从错误码注册表中获取
func NewErrorWithStatusOkAutoMsg(businessCode int) Error {
	return &err{
		HttpCode:     http.StatusOK,
		BusinessCode: businessCode,
		Message:      Text(businessCode),
	}
}

// NewErrorAutoMsg 新建Error，msg自动从错误码注册表中获取
func NewErrorAutoMsg(httpCode, businessCode int) Error {
	return &err{
		HttpCode:     httpCode,
		BusinessCode: businessCode,
		Message:      Text(businessCode),
	}
}

func (e *err) Error() string {
	return e.ToString()
}

// WithErr 封装真实 err
func (e *err) WithErr(err error) Error {
	e.Err = err
	return e
}

// GetHttpCode 获取 HttpCode
func (e *err) GetHttpCode() int {
	return e.HttpCode
}

// GetBusinessCode 获取 BusinessCode
func (e *err) GetBusinessCode() int {
	return e.BusinessCode
}

// GetMsg 获取 Message
func (e *err) GetMsg() string {
	return e.Message
}

// GetErr 获取 Err
func (e *err) GetErr() error {
	return e.Err
}

// ToString 返回 JSON 格式的错误详情
func (e *err) ToString() string {
	err := &struct {
		HttpCode     int    `json:"http_code"`
		BusinessCode int    `json:"business_code"`
		Message      string `json:"message"`
	}{
		HttpCode:     e.HttpCode,
		BusinessCode: e.BusinessCode,
		Message:      e.Message,
	}

	raw, _ := json.Marshal(err)
	return string(raw)
}
