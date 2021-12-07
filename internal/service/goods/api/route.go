package api

func (s *Server) Route(c *Controllers) {

	v1Group := s.Engine.Group("/v1")
	{
		// v1Group.Use(core.WrapAuthHandler(s.Middles.Jwt))

		v1Group.GET("/list", c.goodsController.List)
		v1Group.PUT("/resource", c.goodsController.Add)
		v1Group.POST("/resource", c.goodsController.Update)
		v1Group.DELETE("/resource", c.goodsController.Delete)
	}
}
