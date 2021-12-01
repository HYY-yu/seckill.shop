package svc

import (
	"github.com/HYY-yu/seckill/internal/pkg/cache"
	"github.com/HYY-yu/seckill/internal/pkg/core"
	"github.com/HYY-yu/seckill/internal/pkg/db"
	"github.com/HYY-yu/seckill/internal/service/goods/api/repo"
	"github.com/HYY-yu/seckill/pkg/page"
)

type GoodsSvc struct {
	DB    db.Repo
	Cache cache.Repo

	GoodsRepo *repo.GoodsRepo
}

func NewGoodsSvc(db db.Repo, ca cache.Repo, goodsRepo *repo.GoodsRepo) *GoodsSvc {
	return &GoodsSvc{
		DB:        db,
		Cache:     ca,
		GoodsRepo: goodsRepo,
	}
}

func (s *GoodsSvc) List(ctx core.SvcContext, pr *page.PageRequest) (*page.Page, error) {

	return page.NewPage(
		0,
		nil,
	), nil
}
