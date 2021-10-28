package response

import (
	"encoding/json"
	"github.com/pkg/errors"
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

func NewError(httpCode, businessCode int, msg string) Error {
	return &err{
		HttpCode:     httpCode,
		BusinessCode: businessCode,
		Message:      msg,
	}
}

func NewErrorWithStatusOk(businessCode int, msg string) Error {
	return &err{
		HttpCode:     http.StatusOK,
		BusinessCode: businessCode,
		Message:      msg,
	}
}

func NewErrorWithStatusOkAutoMsg(businessCode int) Error {
	return &err{
		HttpCode:     http.StatusOK,
		BusinessCode: businessCode,
		Message:      Text(businessCode),
	}
}

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

func (e *err) WithErr(err error) Error {
	e.Err = errors.WithStack(err)
	return e
}

func (e *err) GetHttpCode() int {
	return e.HttpCode
}

func (e *err) GetBusinessCode() int {
	return e.BusinessCode
}

func (e *err) GetMsg() string {
	return e.Message
}

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
