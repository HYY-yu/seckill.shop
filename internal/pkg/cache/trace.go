package cache

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/HYY-yu/seckill.shop/internal/service/shop/config"
)

type Redis struct {
	Timestamp   string  `json:"timestamp"`     // 时间，格式：2006-01-02 15:04:05
	Handle      string  `json:"handle"`        // 操作，SET/GET 等
	Key         string  `json:"key"`           // Key
	TTL         float64 `json:"ttl,omitempty"` // 超时时长(单位分)
	CostSeconds float64 `json:"cost_seconds"`  // 执行时间(单位秒)
	Err         error   `json:"err"`
}

func addTracing(ctx context.Context, r *Redis) {
	if r == nil {
		return
	}

	tr := otel.Tracer(config.Get().Server.ServerName + ".Redis")
	_, span := tr.Start(ctx, "Redis."+r.Handle, trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	if r.Err != nil {
		if r.Err != redis.Nil {
			span.SetStatus(codes.Error, fmt.Sprintf("%+v", r.Err))
		}
	}

	span.SetAttributes(attribute.String("timestamp", r.Timestamp))
	span.SetAttributes(attribute.String("handle", r.Handle))
	span.SetAttributes(attribute.String("key", r.Key))
	span.SetAttributes(attribute.Float64("ttl", r.TTL))
	span.SetAttributes(attribute.Float64("cost_seconds", r.CostSeconds))
}
