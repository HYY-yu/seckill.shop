package api

func (s *Server) Route(c *Controllers) {

	v1Group := s.Engine.Group("/v1")
	{
		// v1Group.Use(core.WrapAuthHandler(s.Middles.Jwt))

		v1Group.GET("/list", c.goodsController.List)
	}
}
