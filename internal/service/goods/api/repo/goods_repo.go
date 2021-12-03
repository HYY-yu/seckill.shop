package repo

import (
	"context"

	"github.com/HYY-yu/seckill/internal/service/goods/model"
	"gorm.io/gorm"
)

type GoodsRepo struct {
}

func NewGoodsRepo() *GoodsRepo {
	return &GoodsRepo{}
}

func (*GoodsRepo) ListGoods(ctx context.Context, db *gorm.DB) {
	goodMgr := model.GoodsMgr(db)

	_ = goodMgr
}
