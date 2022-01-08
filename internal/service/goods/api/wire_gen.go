// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package api

import (
	"github.com/HYY-yu/seckill.shop/internal/pkg/cache"
	"github.com/HYY-yu/seckill.shop/internal/pkg/db"
	"github.com/HYY-yu/seckill.shop/internal/service/goods/api/controller"
	"github.com/HYY-yu/seckill.shop/internal/service/goods/api/repo"
	"github.com/HYY-yu/seckill.shop/internal/service/goods/api/svc"
)

// Injectors from wire.go:

// initControllers init Controllers.
func initControllers(d db.Repo, c cache.Repo) (*Controllers, error) {
	goodsRepo := repo.NewGoodsRepo()
	goodsSvc := svc.NewGoodsSvc(d, c, goodsRepo)
	goodsController := controller.NewGoodsController(goodsSvc)
	controllers := NewControllers(goodsController)
	return controllers, nil
}
