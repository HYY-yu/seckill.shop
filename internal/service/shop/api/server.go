package api

import (
	"errors"

	"github.com/HYY-yu/seckill.pkg/cache"
	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/db"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"

	"github.com/HYY-yu/seckill.pkg/pkg/metrics"

	"github.com/HYY-yu/seckill.shop/internal/pkg/middleware"
	"github.com/HYY-yu/seckill.shop/internal/service/shop/api/handler"
	"github.com/HYY-yu/seckill.shop/internal/service/shop/config"

	"github.com/HYY-yu/seckill.pkg/pkg/jaeger"
)

type Handlers struct {
	goodsHandler *handler.GoodsHandler
}

func NewHandlers(
	goodsHandler *handler.GoodsHandler,
) *Handlers {
	return &Handlers{
		goodsHandler: goodsHandler,
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
	tp, err := jaeger.InitJaeger(config.Get().Server.ServerName, config.Get().Jaeger.UdpEndpoint)
	if err != nil {
		logger.Error("jaeger error", zap.Error(err))
	}
	s.Trace = tp

	// Metrics
	metrics.InitMetrics(config.Get().Server.ServerName, "api")

	opts := make([]core.Option, 0)
	opts = append(opts, core.WithEnableCors())
	opts = append(opts, core.WithRecordMetrics(metrics.RecordMetrics))
	if !config.Get().Server.Pprof {
		opts = append(opts, core.WithDisablePProf())
	}

	engine, err := core.New(cfg.Server.ServerName, logger, opts...)
	if err != nil {
		panic(err)
	}
	s.Engine = engine

	s.Middles = middleware.New(logger)

	// Init Repo Svc Handler
	c, err := initHandlers(s.DB, s.Cache)
	if err != nil {
		panic(err)
	}

	s.Route(c)
	return s, nil
}
