package controller

import (
	"context"
	"net/http"

	"github.com/gogf/gf/v2/util/gvalid"

	"github.com/HYY-yu/seckill.pkg/pkg/page"
	"github.com/HYY-yu/seckill.pkg/pkg/response"
	"github.com/HYY-yu/seckill/internal/pkg/core"
	"github.com/HYY-yu/seckill/internal/service/goods/api/svc"
	"github.com/HYY-yu/seckill/internal/service/goods/model"
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
		return
	}
	pageRequest := page.NewPageFromRequest(c.RequestContext().Request.Form)

	data, err := s.goodsSvc.List(c.SvcContext(), pageRequest)
	c.AbortWithError(err)
	c.Payload(data)
}

func (s *GoodsController) Add(c core.Context) {
	params := &model.GoodsAdd{}
	err := c.ShouldBindJSON(params)
	if err != nil {
		c.AbortWithError(response.NewErrorAutoMsg(
			http.StatusBadRequest,
			response.ParamBindError,
		).WithErr(err))
		return
	}

	validErr := gvalid.CheckStruct(context.Background(), params, nil)
	if validErr != nil {
		c.AbortWithError(response.NewError(
			http.StatusBadRequest,
			response.ParamBindError,
			validErr.Error(),
		))
		return
	}

	err = s.goodsSvc.AddGoods(c.SvcContext(), params)
	c.AbortWithError(err)
	c.Payload(nil)
}

func (s *GoodsController) Update(c core.Context) {
	params := &model.GoodsUpdate{}
	err := c.ShouldBindJSON(params)
	if err != nil {
		c.AbortWithError(response.NewErrorAutoMsg(
			http.StatusBadRequest,
			response.ParamBindError,
		).WithErr(err))
		return
	}

	validErr := gvalid.CheckStruct(context.Background(), params, nil)
	if validErr != nil {
		c.AbortWithError(response.NewError(
			http.StatusBadRequest,
			response.ParamBindError,
			validErr.Error(),
		))
		return
	}

	err = s.goodsSvc.UpdateGoods(c.SvcContext(), params)
	c.AbortWithError(err)
	c.Payload(nil)
}

func (s *GoodsController) Delete(c core.Context) {
	type DeleteParam struct {
		Id int `form:"id" v:"required"`
	}
	param := &DeleteParam{}
	err := c.ShouldBindForm(param)
	if err != nil {
		c.AbortWithError(response.NewErrorAutoMsg(
			http.StatusBadRequest,
			response.ParamBindError,
		).WithErr(err))
		return
	}

	validErr := gvalid.CheckStruct(context.Background(), param, nil)
	if validErr != nil {
		c.AbortWithError(response.NewError(
			http.StatusBadRequest,
			response.ParamBindError,
			validErr.Error(),
		))
		return
	}

	err = s.goodsSvc.DeleteGoods(c.SvcContext(), param.Id)
	c.AbortWithError(err)
	c.Payload(nil)
}
