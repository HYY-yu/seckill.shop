package api

func (s *Server) Route(c *Controllers) {

	v1Group := s.Engine.Group("/v1/shop")
	{
		v1Group.GET("/list", c.shopController.ListShop)
	}
}
