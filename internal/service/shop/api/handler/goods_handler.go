package handler

import (
	"context"
	"net/http"

	"github.com/gogf/gf/v2/util/gvalid"

	"github.com/HYY-yu/seckill.pkg/pkg/page"
	"github.com/HYY-yu/seckill.pkg/pkg/response"

	"github.com/HYY-yu/seckill.shop/internal/pkg/core"
	"github.com/HYY-yu/seckill.shop/internal/service/shop/api/svc"
	"github.com/HYY-yu/seckill.shop/internal/service/shop/model"
)

type GoodsHandler struct {
	goodsSvc *svc.GoodsSvc
}

func NewGoodsHandler(goodsSvc *svc.GoodsSvc) *GoodsHandler {
	return &GoodsHandler{
		goodsSvc: goodsSvc,
	}
}

func (s *GoodsHandler) List(c core.Context) {
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

func (s *GoodsHandler) Add(c core.Context) {
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

func (s *GoodsHandler) Update(c core.Context) {
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

func (s *GoodsHandler) Delete(c core.Context) {
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
