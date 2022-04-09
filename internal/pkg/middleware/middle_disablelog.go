package middleware

import (
	"github.com/HYY-yu/seckill.pkg/core"
)

func (m *middleware) DisableLog() core.HandlerFunc {
	return func(c core.Context) {
		c.DisableLog(true)
	}
}
