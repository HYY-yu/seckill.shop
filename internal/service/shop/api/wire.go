//go:build wireinject
// +build wireinject

//go:generate wire gen .
package api

import (
	"github.com/google/wire"
	"go.uber.org/zap"

	"github.com/HYY-yu/seckill.pkg/cache"
	"github.com/HYY-yu/seckill.pkg/db"

	"github.com/HYY-yu/seckill.shop/internal/service/shop/api/grpc_handler"
	"github.com/HYY-yu/seckill.shop/internal/service/shop/api/handler"
	"github.com/HYY-yu/seckill.shop/internal/service/shop/api/repo"
	"github.com/HYY-yu/seckill.shop/internal/service/shop/api/svc"
)

// initHandlers init Handlers.
func initHandlers(l *zap.Logger, d db.Repo, c cache.Repo) (*Handlers, error) {
	panic(wire.Build(repo.NewGoodsRepo, svc.NewGoodsSvc, handler.NewGoodsHandler, grpc_handler.NewGoodsHandler, NewHandlers))
}
