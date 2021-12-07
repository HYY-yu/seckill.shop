package controller

import (
	"context"
	"net/http"

	"github.com/gogf/gf/v2/util/gvalid"

	"github.com/HYY-yu/seckill/internal/pkg/core"
	"github.com/HYY-yu/seckill/internal/service/goods/api/svc"
	"github.com/HYY-yu/seckill/internal/service/goods/model"
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
	err := c.RequestContext().Request.ParseForm()
	if err != nil {
		c.AbortWithError(response.NewErrorAutoMsg(
			http.StatusBadRequest,
			response.ParamBindError,
		).WithErr(err))
	}
	pageRequest := page.NewPageFromRequest(c.RequestContext().Request.Form)
	if err != nil {
		c.AbortWithError(response.NewErrorAutoMsg(
			http.StatusBadRequest,
			response.ParamBindError,
		).WithErr(err))
	}
	data, err := s.goodsSvc.List(c.SvcContext(), pageRequest)
	c.AbortWithError(err)
	c.Payload(data)
}

func (s *GoodsController) Add(c core.Context) {
	params := &model.GoodsAdd{}
	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.AbortWithError(response.NewErrorAutoMsg(
			http.StatusBadRequest,
			response.ParamBindError,
		).WithErr(err))
	}

	validErr := gvalid.CheckStruct(context.Background(), &params, nil)
	if validErr != nil {
		c.AbortWithError(response.NewError(
			http.StatusBadRequest,
			response.ParamBindError,
			validErr.Error(),
		))
	}

	err = s.goodsSvc.AddGoods(c.SvcContext(), params)
	c.AbortWithError(err)
	c.Payload(nil)
}
