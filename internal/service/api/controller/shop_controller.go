package controller

import (
	"github.com/HYY-yu/seckill/internal/pkg/core"
	"github.com/HYY-yu/seckill/internal/service/api/svc"
)

type ShopController struct {
	shopSvc *svc.ShopSvc
}

func NewShopController(shopSvc *svc.ShopSvc) *ShopController {
	return &ShopController{
		shopSvc: shopSvc,
	}
}

func (s *ShopController) ListShop(c core.Context) {

}
