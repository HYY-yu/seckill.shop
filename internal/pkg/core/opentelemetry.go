package core

import (
	"net/http"
	"net/url"
	"time"

	"github.com/HYY-yu/seckill/internal/service/config"
	"github.com/HYY-yu/seckill/pkg/response"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// RecordMetrics 记录prometheus指标用
type RecordMetrics func(method, uri string, success bool, httpCode, businessCode int, costSeconds float64, traceId string)

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

	tr := otel.GetTracerProvider().Tracer(config.Get().Server.ServerName + "/HTTPServer")
	jCtx := otel.GetTextMapPropagator().Extract(ctx.Request.Context(), propagation.HeaderCarrier(ctx.Request.Header))
	jCtx, span := tr.Start(jCtx, ctx.Request.URL.String(), trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()
	traceId := span.SpanContext().SpanID().String()

	// 设置到 全局上下文中
	ctx.Request.WithContext(jCtx)

	ctx.Next()

	if ctx.Writer.Status() == http.StatusNotFound {
		return
	}

	// TODO skip url

	// 获取返回信息
	resp := c.getResponse()
	if resp == nil {
		return
	}
	jsonResp := resp.(*response.JsonResponse)

	// TODO Jaeger追踪

	if o.recordMetrics != nil {
		uri := c.URI()

		// metrics output
		o.recordMetrics(
			c.RequestContext().Request.Method,
			uri,
			!ctx.IsAborted() && ctx.Writer.Status() == http.StatusOK,
			ctx.Writer.Status(),
			jsonResp.Code,
			time.Since(ts).Seconds(),
			traceId,
		)
	}
	if c.getDisableLog() {
		return
	}
	decodedURL, _ := url.QueryUnescape(ctx.Request.URL.RequestURI())

	// logger output
	logger.Info("core-interceptor",
		zap.Any("method", ctx.Request.Method),
		zap.Any("path", decodedURL),
		zap.Any("http_code", ctx.Writer.Status()),
		zap.Any("business_code", jsonResp.Code),
		zap.Any("success", !ctx.IsAborted() && ctx.Writer.Status() == http.StatusOK),
		zap.Any("cost_seconds", time.Since(ts).Seconds()),
	)
}
