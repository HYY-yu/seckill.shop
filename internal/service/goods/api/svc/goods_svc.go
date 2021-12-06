package svc

import (
	"net/http"

	"github.com/HYY-yu/seckill/internal/pkg/cache"
	"github.com/HYY-yu/seckill/internal/pkg/core"
	"github.com/HYY-yu/seckill/internal/pkg/db"
	"github.com/HYY-yu/seckill/internal/service/goods/api/repo"
	"github.com/HYY-yu/seckill/internal/service/goods/model"
	"github.com/HYY-yu/seckill/pkg/page"
	"github.com/HYY-yu/seckill/pkg/response"
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

func (s *GoodsSvc) List(sctx core.SvcContext, pr *page.PageRequest) (*page.Page, error) {
	ctx := sctx.Context()

	limit, offset := pr.GetLimitAndOffset()
	pr.AddAllowSortField(model.GoodsColumns.CreateTime)
	sort, _ := pr.Sort()

	data, err := s.GoodsRepo.ListGoods(ctx, s.DB.GetDb(ctx), limit, offset, pr.Filter, sort)
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	count, err := s.GoodsRepo.CountGoods(ctx, s.DB.GetDb(ctx), pr.Filter)
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	var result = make([]model.GoodsListResp, len(data))
	for i := range data {
		r := model.GoodsListResp{
			ID:         data[i].ID,
			Name:       data[i].Name,
			Desc:       data[i].Desc,
			Count:      data[i].Count,
			CreateTime: data[i].CreateTime,
		}
		result[i] = r
	}

	return page.NewPage(
		count,
		result,
	), nil
}

func (s *GoodsSvc) AddGoods(sctx core.SvcContext) {

}
