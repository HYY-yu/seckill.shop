//go:build wireinject
// +build wireinject

//go:generate wire gen .
package api

import (
	"github.com/HYY-yu/seckill/internal/pkg/cache"
	"github.com/HYY-yu/seckill/internal/pkg/db"
	"github.com/HYY-yu/seckill/internal/service/goods/api/controller"
	"github.com/HYY-yu/seckill/internal/service/goods/api/repo"
	"github.com/HYY-yu/seckill/internal/service/goods/api/svc"
	"github.com/google/wire"
)

// initControllers init Controllers.
func initControllers(d db.Repo, c cache.Repo) (*Controllers, error) {
	panic(wire.Build(repo.NewGoodsRepo, svc.NewGoodsSvc, controller.NewGoodsController, NewControllers))
}
