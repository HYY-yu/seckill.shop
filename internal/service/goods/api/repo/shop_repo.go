package repo

import (
	"github.com/HYY-yu/seckill/internal/pkg/db"
	"go.uber.org/zap"
)

type ShopRepo struct {
	Logger *zap.Logger
	DB     db.Repo
}

func NewShopRepo(logger *zap.Logger, db db.Repo) *ShopRepo {
	return &ShopRepo{
		Logger: logger,
		DB:     db,
	}
}
