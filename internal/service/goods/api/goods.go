package api

import (
	"github.com/HYY-yu/seckill/internal/pkg/middleware"
	"github.com/HYY-yu/seckill/internal/service/goods/api/controller"
	"github.com/HYY-yu/seckill/internal/service/goods/config"
	"github.com/HYY-yu/seckill/pkg/metrics"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"

	"github.com/HYY-yu/seckill/internal/pkg/cache"
	"github.com/HYY-yu/seckill/internal/pkg/core"
	"github.com/HYY-yu/seckill/internal/pkg/db"
	"github.com/HYY-yu/seckill/pkg/jaeger"
	"github.com/HYY-yu/seckill/pkg/werror"
)

type Controllers struct {
	goodsController *controller.GoodsController
}

func NewControllers(
	goodsController *controller.GoodsController,
) *Controllers {
	return &Controllers{
		goodsController: goodsController,
	}
}

type Server struct {
	Logger  *zap.Logger
	Engine  core.Engine
	DB      db.Repo
	Cache   cache.Repo
	Trace   *trace.TracerProvider
	Middles middleware.Middleware
}

func NewApiServer(logger *zap.Logger) (*Server, error) {
	if logger == nil {
		return nil, werror.New("logger required")
	}
	s := &Server{}
	s.Logger = logger

	dbRepo, err := db.New()
	if err != nil {
		logger.Fatal("new db err", zap.Error(err))
	}
	s.DB = dbRepo

	cacheRepo, err := cache.New()
	if err != nil {
		logger.Fatal("new cache err", zap.Error(err))
	}
	s.Cache = cacheRepo

	// Jaeger
	tp, err := jaeger.InitJaeger(config.Get().Server.ServerName, config.Get().Jaeger.UdpEndpoint)
	if err != nil {
		logger.Error("jaeger error", zap.Error(err))
	}
	s.Trace = tp

	// Metrics
	metrics.InitMetrics(config.Get().Server.ServerName, "api")

	engine, err := core.New(logger,
		core.WithEnableCors(),
		core.WithRecordMetrics(metrics.RecordMetrics),
	)
	if err != nil {
		panic(err)
	}
	s.Engine = engine

	s.Middles = middleware.New(logger)

	// Init Repo Svc Controller
	c, err := initControllers(logger, s.DB, s.Cache)
	if err != nil {
		panic(err)
	}

	s.Route(c)
	return s, nil
}
