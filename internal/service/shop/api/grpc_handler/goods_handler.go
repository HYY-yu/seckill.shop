package grpc_handler

import (
	"context"

	"go.uber.org/zap"

	"github.com/HYY-yu/seckill.shop/internal/service/shop/api/svc"
	"github.com/HYY-yu/seckill.shop/proto"
)

type GoodsHandler struct {
	proto.UnimplementedShopServer
	logger   *zap.Logger
	goodsSvc *svc.GoodsSvc
}

func NewGoodsHandler(logger *zap.Logger, goodsSvc *svc.GoodsSvc) *GoodsHandler {
	return &GoodsHandler{
		logger:   logger,
		goodsSvc: goodsSvc,
	}
}

func (g *GoodsHandler) List(ctx context.Context, req *proto.ListReq) (*proto.ListResp, error) {
	data, err := g.goodsSvc.GrpcList(ctx, req)
	return &proto.ListResp{
		Data: data,
	}, err
}

func (g *GoodsHandler) Incr(ctx context.Context, req *proto.IncrReq) (*proto.IncrResp, error) {
	err := g.goodsSvc.IncrCount(ctx, req)
	incrResp := &proto.IncrResp{
		OK: err == nil,
	}
	return incrResp, err
}
