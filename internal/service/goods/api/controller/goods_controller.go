package controller

import (
	"github.com/HYY-yu/seckill/internal/pkg/core"
	"github.com/HYY-yu/seckill/internal/service/goods/api/svc"
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

}
