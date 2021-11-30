package svc

import (
	"errors"

	"github.com/HYY-yu/seckill/internal/pkg/cache"
	"github.com/HYY-yu/seckill/internal/pkg/db"
	"github.com/HYY-yu/seckill/internal/service/goods/api/repo"
	"github.com/HYY-yu/seckill/pkg/page"
	"github.com/HYY-yu/seckill/pkg/response"
	"github.com/HYY-yu/seckill/pkg/werror"
	"go.uber.org/zap"
)

type GoodsSvc struct {
	Logger *zap.Logger
	DB     db.Repo
	Cache  cache.Repo

	GoodsRepo *repo.GoodsRepo
}

func NewGoodsSvc(logger *zap.Logger, db db.Repo, ca cache.Repo, goodsRepo *repo.GoodsRepo) *GoodsSvc {
	return &GoodsSvc{
		Logger:    logger,
		DB:        db,
		Cache:     ca,
		GoodsRepo: goodsRepo,
	}
}

func (s *GoodsSvc) List(pr *page.PageRequest) error {
	return response.NewErrorWithStatusOkAutoMsg(response.ServerError).
		WithErr(werror.WithStack(errors.New("haha error")))
}
