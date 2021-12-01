package core

import (
	"net/http"
	"time"

	"github.com/HYY-yu/seckill/internal/service/goods/config"
	"github.com/HYY-yu/seckill/pkg/response"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// RecordMetrics 记录prometheus指标用
type RecordMetrics func(method, uri string, httpCode, businessCode int, costSeconds float64, traceId string)

// OpenTelemetry
// logger \ metrics \ trace 三者归一
type OpenTelemetry struct {
	recordMetrics RecordMetrics
}

func NewOpenTelemetry(rm RecordMetrics) *OpenTelemetry {
	return &OpenTelemetry{
		recordMetrics: rm,
	}
}

func (o *OpenTelemetry) Telemetry(c Context) {
	ctx := c.RequestContext()
	logger := c.Logger()
	ts := time.Now()

	tr := otel.GetTracerProvider().Tracer(config.Get().Server.ServerName + ".HTTP")
	jCtx := otel.GetTextMapPropagator().Extract(ctx.Request.Context(), propagation.HeaderCarrier(ctx.Request.Header))
	jCtx, span := tr.Start(jCtx, ctx.Request.URL.String(), trace.WithSpanKind(trace.SpanKindServer), trace.WithNewRoot())
	defer span.End()
	traceId := span.SpanContext().SpanID().String()

	// 设置到 请求上下文 中
	ctx.Request.WithContext(jCtx)

	// 设置到logger中
	logger = logger.With(zap.String("trace_id", traceId))
	c.setLogger(logger)

	ctx.Next()

	if ctx.Writer.Status() == http.StatusNotFound {
		return
	}

	// 获取返回信息
	resp := c.getResponse()
	if resp == nil {
		return
	}
	jsonResp := resp.(*response.JsonResponse)

	decodedURL := c.URI()
	telemetry := &RequestTelemetry{
		Method:       ctx.Request.Method,
		Path:         decodedURL,
		HttpCode:     ctx.Writer.Status(),
		BusinessCode: jsonResp.Code,
		CostSeconds:  time.Since(ts).Seconds(),
	}

	//  Jaeger追踪
	span.SetAttributes(
		attribute.String("http.method", telemetry.Method),
		attribute.String("http.path", telemetry.Path),
		attribute.Int("http.http_code", telemetry.HttpCode),
		attribute.Int("http.business_code", telemetry.BusinessCode),
		attribute.Float64("http.cost_seconds", telemetry.CostSeconds),
	)
	if !ctx.IsAborted() && ctx.Writer.Status() == http.StatusOK {
		span.SetStatus(codes.Ok, "")
	} else {
		span.SetStatus(codes.Error, jsonResp.Message)
	}

	// Metrics
	if o.recordMetrics != nil {
		// metrics output
		o.recordMetrics(
			telemetry.Method,
			telemetry.Path,
			telemetry.HttpCode,
			telemetry.BusinessCode,
			telemetry.CostSeconds,
			traceId,
		)
	}
	if c.getDisableLog() {
		return
	}

	// logger output
	logger.Info("core-interceptor",
		zap.Any("method", telemetry.Method),
		zap.Any("path", telemetry.Path),
		zap.Any("http_code", telemetry.HttpCode),
		zap.Any("business_code", telemetry.BusinessCode),
		zap.Any("cost_seconds", telemetry.CostSeconds),
		zap.Any("trace_id", traceId),
	)
}

type RequestTelemetry struct {
	Method       string  `json:"method"`
	Path         string  `json:"path"`
	HttpCode     int     `json:"http_code"`
	BusinessCode int     `json:"business_code"`
	CostSeconds  float64 `json:"cost_seconds"`
}
