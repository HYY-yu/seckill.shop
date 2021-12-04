package repo

import (
	"context"

	"github.com/HYY-yu/seckill/internal/service/goods/model"
	"github.com/HYY-yu/seckill/pkg/util"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type GoodsRepo struct {
}

func NewGoodsRepo() *GoodsRepo {
	return &GoodsRepo{}
}

func (*GoodsRepo) ListGoods(ctx context.Context, db *gorm.DB,
	limit, offset int,
	filter map[string]interface{},
	sort string,
) (result []model.Goods, err error) {
	goodMgr := model.GoodsMgr(db)

	err = goodMgr.
		AddWhere(filter[model.GoodsColumns.Name], util.IsZero, func(db *gorm.DB, v interface{}) *gorm.DB {
			return db.Where(model.GoodsColumns.Name+" LIKE ?", "%"+cast.ToString(v)+"%")
		}).
		AddWhere(filter[model.GoodsColumns.ID], util.IsZero, func(db *gorm.DB, i interface{}) *gorm.DB {
			return db.Where(model.GoodsColumns.ID+" = ?", i)
		}).
		Sort(sort, "id desc").
		Limit(limit).
		Offset(offset).
		Find(&result).Error
	return
}

func (*GoodsRepo) CountGoods(ctx context.Context, db *gorm.DB,
	filter map[string]interface{},
) (count int64, err error) {
	goodMgr := model.GoodsMgr(db)

	err = goodMgr.
		AddWhere(filter[model.GoodsColumns.Name], util.IsZero, func(db *gorm.DB, v interface{}) *gorm.DB {
			return db.Where(model.GoodsColumns.Name+" LIKE ?", "%"+cast.ToString(v)+"%")
		}).
		AddWhere(filter[model.GoodsColumns.ID], util.IsZero, func(db *gorm.DB, i interface{}) *gorm.DB {
			return db.Where(model.GoodsColumns.ID+" = ?", i)
		}).
		Count(&count).Error
	return
}
