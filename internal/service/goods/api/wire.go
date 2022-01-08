//go:build wireinject
// +build wireinject

//go:generate wire gen .
package api

import (
	"github.com/google/wire"

	"github.com/HYY-yu/seckill.shop/internal/pkg/cache"
	"github.com/HYY-yu/seckill.shop/internal/pkg/db"
	"github.com/HYY-yu/seckill.shop/internal/service/goods/api/controller"
	"github.com/HYY-yu/seckill.shop/internal/service/goods/api/repo"
	"github.com/HYY-yu/seckill.shop/internal/service/goods/api/svc"
)

// initControllers init Controllers.
func initControllers(d db.Repo, c cache.Repo) (*Controllers, error) {
	panic(wire.Build(repo.NewGoodsRepo, svc.NewGoodsSvc, controller.NewGoodsController, NewControllers))
}
