package api

func (s *Server) Route(c *Handlers) {

	v1Group := s.Engine.Group("/v1")
	{
		// v1Group.Use(core.WrapAuthHandler(s.Middles.Jwt))

		v1Group.GET("/list", c.goodsHandler.List)
		v1Group.PUT("/resource", c.goodsHandler.Add)
		v1Group.POST("/resource", c.goodsHandler.Update)
		v1Group.DELETE("/resource", c.goodsHandler.Delete)
	}
}
