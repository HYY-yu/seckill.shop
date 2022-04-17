package client

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/HYY-yu/seckill.shop/proto"
)

func Connect(host string) (proto.ShopClient, error) {
	// 尝试连接GRPC
	var optsGrpc []grpc.DialOption
	optsGrpc = append(optsGrpc,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithChainStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)

	clientConn, err := grpc.Dial(host, optsGrpc...)
	if err != nil {
		return nil, err
	}

	shopClient := proto.NewShopClient(clientConn)
	return shopClient, nil
}
