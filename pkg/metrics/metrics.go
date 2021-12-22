package metrics

// package metrics 用于注册和记录 HTTP 服务一般所需要的几个指标
// 用于 prometheus 的分析

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cast"
)

var metricsRequestsTotal *prometheus.CounterVec

var metricsRequestsCost *prometheus.HistogramVec

// InitMetrics 主动注册 Metric 指标
func InitMetrics(namespace string, subsystem string) {
	// metricsRequestsTotal metrics for request total 计数器（Counter）
	metricsRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "requests_total",
			Help:      "request(ms) total",
		},
		[]string{"method", "path"},
	)

	// metricsRequestsCost metrics for requests cost 累积直方图（Histogram）
	metricsRequestsCost = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "requests_cost",
			Help:      "request(ms) cost milliseconds",
		},
		[]string{"method", "path", "http_code", "business_code", "cost_milliseconds", "trace_id"},
	)

	prometheus.MustRegister(metricsRequestsTotal, metricsRequestsCost)
}

// RecordMetrics 记录指标
// 请注意需要先调用 InitMetrics
func RecordMetrics(method, uri string, httpCode, businessCode int, costSeconds float64, traceId string) {
	metricsRequestsTotal.With(prometheus.Labels{
		"method": method,
		"path":   uri,
	}).Inc()

	metricsRequestsCost.With(prometheus.Labels{
		"method":            method,
		"path":              uri,
		"http_code":         cast.ToString(httpCode),
		"business_code":     cast.ToString(businessCode),
		"cost_milliseconds": cast.ToString(costSeconds * 1000),
		"trace_id":          traceId,
	}).Observe(costSeconds)
}
