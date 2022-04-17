package api

import (
	"errors"
	"net/http"

	"github.com/HYY-yu/seckill.pkg/cache"
	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/core/middleware"
	"github.com/HYY-yu/seckill.pkg/db"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/HYY-yu/seckill.pkg/pkg/metrics"

	"github.com/HYY-yu/seckill.pkg/pkg/jaeger"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/HYY-yu/seckill.shop/internal/service/shop/api/grpc_handler"
	"github.com/HYY-yu/seckill.shop/internal/service/shop/api/handler"
	"github.com/HYY-yu/seckill.shop/internal/service/shop/config"
	"github.com/HYY-yu/seckill.shop/proto"
)

type Handlers struct {
	goodsHandler     *handler.GoodsHandler
	grpcGoodsHandler *grpc_handler.GoodsHandler
}

func NewHandlers(
	goodsHandler *handler.GoodsHandler,
	grpcGoodsHandler *grpc_handler.GoodsHandler,
) *Handlers {
	return &Handlers{
		goodsHandler:     goodsHandler,
		grpcGoodsHandler: grpcGoodsHandler,
	}
}

type Server struct {
	Logger      *zap.Logger
	HttpServer  *http.Server
	GrpcServer  *grpc.Server
	DB          db.Repo
	Cache       cache.Repo
	Trace       *trace.TracerProvider
	HTTPMiddles middleware.Middleware
}

func NewApiServer(logger *zap.Logger) (*Server, error) {
	if logger == nil {
		return nil, errors.New("logger required")
	}
	s := &Server{}
	s.Logger = logger
	cfg := config.Get()

	dbRepo, err := db.New(&db.DBConfig{
		User:            cfg.MySQL.Base.User,
		Pass:            cfg.MySQL.Base.Pass,
		Addr:            cfg.MySQL.Base.Addr,
		Name:            cfg.MySQL.Base.Name,
		MaxOpenConn:     cfg.MySQL.Base.MaxOpenConn,
		MaxIdleConn:     cfg.MySQL.Base.MaxIdleConn,
		ConnMaxLifeTime: cfg.MySQL.Base.ConnMaxLifeTime,
		ServerName:      cfg.Server.ServerName,
	})
	if err != nil {
		logger.Fatal("new db err", zap.Error(err))
	}
	s.DB = dbRepo

	cacheRepo, err := cache.New(cfg.Server.ServerName, &cache.RedisConf{
		Addr:         cfg.Redis.Addr,
		Pass:         cfg.Redis.Pass,
		Db:           cfg.Redis.Db,
		MaxRetries:   cfg.Redis.MaxRetries,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
	})
	if err != nil {
		logger.Fatal("new cache err", zap.Error(err))
	}
	s.Cache = cacheRepo

	// Jaeger
	var tp *trace.TracerProvider
	if cfg.Jaeger.StdOut {
		tp, err = jaeger.InitStdOutForDevelopment(cfg.Server.ServerName, cfg.Jaeger.UdpEndpoint)
	} else {
		tp, err = jaeger.InitJaeger(cfg.Server.ServerName, cfg.Jaeger.UdpEndpoint)
	}
	if err != nil {
		logger.Error("jaeger error", zap.Error(err))
	}
	s.Trace = tp

	// Metrics
	metrics.InitMetrics(cfg.Server.ServerName, "api")
	err = metrics.InitGrpcMetrics()
	if err != nil {
		panic(err)
	}

	// Repo Svc Handler
	c, err := initHandlers(logger, s.DB, s.Cache)
	if err != nil {
		panic(err)
	}

	// HTTP Server
	opts := make([]core.Option, 0)
	opts = append(opts, core.WithEnableCors())
	opts = append(opts, core.WithRecordMetrics(metrics.RecordMetrics))
	if !cfg.Server.Pprof {
		opts = append(opts, core.WithDisablePProf())
	}
	engine, err := core.New(cfg.Server.ServerName, logger, opts...)
	if err != nil {
		panic(err)
	}
	// Init HTTP Middles
	s.HTTPMiddles = middleware.New(logger, cfg.JWT.Secret)

	// Route
	s.Route(c, engine)
	server := &http.Server{
		Handler: engine,
	}
	s.HttpServer = server

	// GRPC Server
	var optsGrpc []grpc.ServerOption
	optsGrpc = append(optsGrpc, grpc.ChainUnaryInterceptor(
		otelgrpc.UnaryServerInterceptor(),
		grpc_recovery.UnaryServerInterceptor(),
		grpc_zap.UnaryServerInterceptor(logger),
		metrics.GRPCMetrics.UnaryServerInterceptor(),
	))
	optsGrpc = append(optsGrpc, grpc.ChainStreamInterceptor(
		otelgrpc.StreamServerInterceptor(),
		grpc_recovery.StreamServerInterceptor(),
		grpc_zap.StreamServerInterceptor(logger),
		metrics.GRPCMetrics.StreamServerInterceptor(),
	))
	grpcServer := grpc.NewServer(optsGrpc...)
	proto.RegisterShopServer(grpcServer, c.grpcGoodsHandler)
	s.GrpcServer = grpcServer

	return s, nil
}
