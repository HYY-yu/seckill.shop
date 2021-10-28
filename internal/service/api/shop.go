package api

import (
	"go.uber.org/zap"

	"github.com/HYY-yu/seckill/internal/pkg/cache"
	"github.com/HYY-yu/seckill/internal/pkg/core"
	"github.com/HYY-yu/seckill/internal/pkg/db"
	"github.com/HYY-yu/seckill/internal/pkg/metrics"
	"github.com/HYY-yu/seckill/internal/service/api/router"
	"github.com/HYY-yu/seckill/internal/service/api/router/middleware"
	"github.com/HYY-yu/seckill/pkg/errors"
)

type Server struct {
	Logger  *zap.Logger
	Engine  core.Engine
	DB      db.Repo
	Cache   cache.Repo
	Middles middleware.Middleware
}

func NewApiServer(logger *zap.Logger) (*Server, error) {
	if logger == nil {
		return nil, errors.New("logger required")
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

	engine, err := core.New(logger,
		core.WithEnableCors(),
		core.WithEnableRate(),
		core.WithRecordMetrics(metrics.RecordMetrics),
	)
	if err != nil {
		panic(err)
	}
	s.Engine = engine

	s.Middles = middleware.New(logger)

	router.SetRouter(s.Engine)
	return s, nil
}
