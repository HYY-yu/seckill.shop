package controller

import (
	"net/http"

	"github.com/HYY-yu/seckill/internal/pkg/core"
	"github.com/HYY-yu/seckill/internal/service/goods/api/svc"
	"github.com/HYY-yu/seckill/pkg/page"
	"github.com/HYY-yu/seckill/pkg/response"
)

type GoodsController struct {
	goodsSvc *svc.GoodsSvc
}

func NewGoodsController(goodsSvc *svc.GoodsSvc) *GoodsController {
	return &GoodsController{
		goodsSvc: goodsSvc,
	}
}

func (s *GoodsController) List(c core.Context) {
	var pageRequest page.PageRequest
	err := c.ShouldBindForm(&pageRequest)
	if err != nil {
		c.AbortWithError(response.NewErrorAutoMsg(
			http.StatusBadRequest,
			response.ParamBindError,
		).WithErr(err))
	}
	data, err := s.goodsSvc.List(c.SvcContext(), &pageRequest)
	c.AbortWithError(err)
	c.Payload(data)
}
