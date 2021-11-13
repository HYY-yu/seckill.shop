package repo

import (
	"github.com/HYY-yu/seckill/internal/pkg/db"
	"go.uber.org/zap"
)

type GoodsRepo struct {
	Logger *zap.Logger
	DB     db.Repo
}

func NewGoodsRepo(logger *zap.Logger, db db.Repo) *GoodsRepo {
	return &GoodsRepo{
		Logger: logger,
		DB:     db,
	}
}
