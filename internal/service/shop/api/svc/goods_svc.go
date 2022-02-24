package svc

import (
	"net/http"
	"time"

	"github.com/HYY-yu/seckill.pkg/pkg/mysqlerr_helper"
	"github.com/HYY-yu/seckill.pkg/pkg/page"
	"github.com/HYY-yu/seckill.pkg/pkg/response"

	"github.com/HYY-yu/seckill.shop/internal/pkg/cache"
	"github.com/HYY-yu/seckill.shop/internal/pkg/core"
	"github.com/HYY-yu/seckill.shop/internal/pkg/db"
	"github.com/HYY-yu/seckill.shop/internal/service/shop/api/repo"
	"github.com/HYY-yu/seckill.shop/internal/service/shop/model"
)

type GoodsSvc struct {
	DB    db.Repo
	Cache cache.Repo

	GoodsRepo repo.GoodsRepo
}

func NewGoodsSvc(db db.Repo, ca cache.Repo, goodsRepo repo.GoodsRepo) *GoodsSvc {
	return &GoodsSvc{
		DB:        db,
		Cache:     ca,
		GoodsRepo: goodsRepo,
	}
}

func (s *GoodsSvc) List(sctx core.SvcContext, pr *page.PageRequest) (*page.Page, error) {
	ctx := sctx.Context()
	mgr := s.GoodsRepo.Mgr(ctx, s.DB.GetDb(ctx))

	limit, offset := pr.GetLimitAndOffset()
	pr.AddAllowSortField(model.GoodsColumns.CreateTime)
	sort, _ := pr.Sort()

	data, err := mgr.ListGoods(limit, offset, pr.Filter, sort)
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	count, err := mgr.CountGoods(pr.Filter)
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	var result = make([]model.GoodsListResp, len(data))
	for i, e := range data {
		r := model.GoodsListResp{
			ID:         e.ID,
			Name:       e.Name,
			Desc:       e.Desc,
			Count:      e.Count,
			CreateTime: e.CreateTime,
		}
		result[i] = r
	}

	return page.NewPage(
		count,
		result,
	), nil
}

func (s *GoodsSvc) AddGoods(sctx core.SvcContext, param *model.GoodsAdd) error {
	ctx := sctx.Context()
	mgr := s.GoodsRepo.Mgr(ctx, s.DB.GetDb(ctx))
	now := time.Now().Unix()

	bean := &model.Goods{
		Name:       param.Name,
		Desc:       param.Desc,
		Count:      int(param.Count),
		CreateTime: int(now),
	}

	err := mgr.CreateGoods(bean)
	if err != nil {
		if mysqlerr_helper.IsMysqlDupEntryError(err) {
			return response.NewErrorWithStatusOk(
				response.ParamBindError,
				"商品名重复",
			)
		}
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	return nil
}

func (s *GoodsSvc) UpdateGoods(sctx core.SvcContext, param *model.GoodsUpdate) error {
	ctx := sctx.Context()
	mgr := s.GoodsRepo.Mgr(ctx, s.DB.GetDb(ctx))

	bean := &model.Goods{
		ID: param.Id,
	}
	updateColumns := make([]string, 0)

	if param.Name != nil {
		bean.Name = *param.Name
		updateColumns = append(updateColumns, model.GoodsColumns.Name)
	}
	if param.Desc != nil {
		bean.Desc = *param.Desc
		updateColumns = append(updateColumns, model.GoodsColumns.Desc)
	}
	if param.Count != nil {
		bean.Count = int(*param.Count)
		updateColumns = append(updateColumns, model.GoodsColumns.Count)
	}

	err := mgr.WithSelects(model.GoodsColumns.ID, updateColumns...).UpdateGoods(bean)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	return nil
}

func (s *GoodsSvc) DeleteGoods(sctx core.SvcContext, goodsId int) error {
	// 软删除
	ctx := sctx.Context()
	mgr := s.GoodsRepo.Mgr(ctx, s.DB.GetDb(ctx))
	now := time.Now().Unix()

	bean := &model.Goods{
		ID:         goodsId,
		DeleteTime: int(now),
	}
	err := mgr.UpdateGoods(bean)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	return nil
}
