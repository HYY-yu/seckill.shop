package api

import (
	"github.com/HYY-yu/seckill.pkg/core"
)

func (s *Server) Route(c *Handlers, engine core.Engine) {

	v1Group := engine.Group("/v1")
	{
		// v1Group.Use(core.WrapAuthHandler(s.HTTPMiddles.Jwt))

		v1Group.GET("/list", c.goodsHandler.List)
		v1Group.PUT("/resource", c.goodsHandler.Add)
		v1Group.POST("/resource", c.goodsHandler.Update)
		v1Group.DELETE("/resource", c.goodsHandler.Delete)
	}
}
