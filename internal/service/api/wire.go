//go:build wireinject
// +build wireinject

//go:generate wire gen .
package api

import (
	"github.com/HYY-yu/seckill/internal/pkg/cache"
	"github.com/HYY-yu/seckill/internal/pkg/db"
	"github.com/HYY-yu/seckill/internal/service/api/controller"
	"github.com/HYY-yu/seckill/internal/service/api/repo"
	"github.com/HYY-yu/seckill/internal/service/api/svc"
	"github.com/google/wire"
	"go.uber.org/zap"
)

// initControllers init Controllers.
func initControllers(l *zap.Logger, d db.Repo, c cache.Repo) (*Controllers, error) {
	panic(wire.Build(repo.NewShopRepo, svc.NewShopSvc, controller.NewShopController, NewControllers))
}
