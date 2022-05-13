package repo

import (
	"context"

	"github.com/spf13/cast"
	"gorm.io/gorm"

	"github.com/HYY-yu/seckill.pkg/pkg/util"

	"github.com/HYY-yu/seckill.shop/internal/service/shop/model"
)

type GoodsRepo interface {
	Mgr(ctx context.Context, db *gorm.DB) *_GoodsMgr
}

// goodsRepo 薄薄的一层，用来封装_xxMgr
// Repo 中不要出现字段，否则容易出现并发安全问题。
type goodsRepo struct {
}

func NewGoodsRepo() GoodsRepo {
	return &goodsRepo{}
}

func (*goodsRepo) Mgr(ctx context.Context, db *gorm.DB) *_GoodsMgr {
	goodsMgr := GoodsMgr(ctx, db)
	return goodsMgr
}

// ------- 自定义方法 -------

func (obj *_GoodsMgr) ListGoods(
	limit, offset int,
	filter map[string]interface{},
	sort string,
) (result []model.Goods, err error) {
	err = obj.
		addWhere(filter[model.GoodsColumns.Name], util.IsNotZero, func(db *gorm.DB, v interface{}) *gorm.DB {
			return db.Where(model.GoodsColumns.Name+" LIKE ?", "%"+cast.ToString(v)+"%")
		}).
		addWhere(filter[model.GoodsColumns.ID], util.IsNotZero, func(db *gorm.DB, v interface{}) *gorm.DB {
			return db.Where(model.GoodsColumns.ID+" = ?", v)
		}).
		addWhere(filter["ids"], util.IsNotZero, func(db *gorm.DB, v interface{}) *gorm.DB {
			return db.Where(model.GoodsColumns.ID+" IN (?)", v)
		}).
		sort(sort, model.GoodsColumns.ID+" desc").
		Where(model.GoodsColumns.DeleteTime + " = 0").
		Limit(limit).
		Offset(offset).
		Find(&result).Error
	return
}

func (obj *_GoodsMgr) CountGoods(
	filter map[string]interface{},
) (count int64, err error) {
	err = obj.
		addWhere(filter[model.GoodsColumns.Name], util.IsNotZero, func(db *gorm.DB, v interface{}) *gorm.DB {
			return db.Where(model.GoodsColumns.Name+" LIKE ?", "%"+cast.ToString(v)+"%")
		}).
		addWhere(filter[model.GoodsColumns.ID], util.IsNotZero, func(db *gorm.DB, i interface{}) *gorm.DB {
			return db.Where(model.GoodsColumns.ID+" = ?", i)
		}).
		Where(model.GoodsColumns.DeleteTime + " = 0").
		Count(&count).Error
	return
}

func (obj *_GoodsMgr) IncrCount(id int, count int) error {
	db := obj.DB.WithContext(obj.ctx)

	err := db.Where(model.GoodsColumns.ID+" = ?", id).UpdateColumn(model.GoodsColumns.Count, gorm.Expr(model.GoodsColumns.Count+"+ ?", count)).Error
	return err
}
