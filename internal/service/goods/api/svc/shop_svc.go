package svc

import (
	"github.com/HYY-yu/seckill/internal/pkg/cache"
	"github.com/HYY-yu/seckill/internal/pkg/db"
	"github.com/HYY-yu/seckill/internal/service/goods/api/repo"
	"go.uber.org/zap"
)

type ShopSvc struct {
	Logger *zap.Logger
	DB     db.Repo
	Cache  cache.Repo

	ShopRepo *repo.ShopRepo
}

func NewShopSvc(logger *zap.Logger, db db.Repo, ca cache.Repo, shopRepo *repo.ShopRepo) *ShopSvc {
	return &ShopSvc{
		Logger:   logger,
		DB:       db,
		Cache:    ca,
		ShopRepo: shopRepo,
	}
}

func (s *ShopSvc) ListShop() {

}
