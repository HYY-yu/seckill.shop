package middleware

import (
	"github.com/HYY-yu/seckill.pkg/core"
	"go.uber.org/zap"

	"github.com/HYY-yu/seckill.pkg/pkg/response"
)

var _ Middleware = (*middleware)(nil)

type Middleware interface {
	// i 为了避免被其他包实现
	i()

	// Jwt 中间件
	Jwt(ctx core.Context) (userId int64, userName string, err response.Error)

	// DisableLog 不记录日志
	DisableLog() core.HandlerFunc
}

type middleware struct {
	logger *zap.Logger
}

func New(logger *zap.Logger) Middleware {
	return &middleware{
		logger: logger,
	}
}

func (m *middleware) i() {}
