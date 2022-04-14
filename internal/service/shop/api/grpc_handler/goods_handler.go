package grpc_handler

import (
	"go.uber.org/zap"

	"github.com/HYY-yu/seckill.shop/internal/service/shop/api/grpc_handler/proto"
	"github.com/HYY-yu/seckill.shop/internal/service/shop/api/svc"
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

func (g GoodsHandler) List(req *proto.ListReq, server proto.Shop_ListServer) error {

	return nil
}
