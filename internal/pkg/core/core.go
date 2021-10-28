package core

import (
	"errors"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	cors "github.com/rs/cors/wrapper/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"github.com/HYY-yu/seckill/internal/service/config"
	"github.com/HYY-yu/seckill/pkg/response"
)

const _MaxBurstSize = 100000

type Option func(*option)

type option struct {
	disablePProf      bool
	disableSwagger    bool
	disablePrometheus bool
	panicNotify       OnPanicNotify
	recordMetrics     RecordMetrics
	enableCors        bool
	enableRate        bool
	enableOpenBrowser string
}

// OnPanicNotify 发生panic时通知用
type OnPanicNotify func(ctx Context, err interface{}, stackInfo string)

// RecordMetrics 记录prometheus指标用
type RecordMetrics func(method, uri string, success bool, httpCode, businessCode int, costSeconds float64, traceId string)

// WithDisablePProf 禁用 pprof
func WithDisablePProf() Option {
	return func(opt *option) {
		opt.disablePProf = true
	}
}

// WithDisableSwagger 禁用 swagger
func WithDisableSwagger() Option {
	return func(opt *option) {
		opt.disableSwagger = true
	}
}

// WithDisablePrometheus 禁用prometheus
func WithDisablePrometheus() Option {
	return func(opt *option) {
		opt.disablePrometheus = true
	}
}

// WithPanicNotify 设置panic时的通知回调
func WithPanicNotify(notify OnPanicNotify) Option {
	return func(opt *option) {
		opt.panicNotify = notify
	}
}

// WithRecordMetrics 设置记录prometheus记录指标回调
func WithRecordMetrics(record RecordMetrics) Option {
	return func(opt *option) {
		opt.recordMetrics = record
	}
}

// WithEnableCors 开启CORS
func WithEnableCors() Option {
	return func(opt *option) {
		opt.enableCors = true
	}
}

func WithEnableRate() Option {
	return func(opt *option) {
		opt.enableRate = true
	}
}

// WrapAuthHandler 用来处理 Auth 的入口，在之后的handler中只需 ctx.UserID() ctx.UserName() 即可。
// handler 是真正的处理者
func WrapAuthHandler(handler func(Context) (userID int64, userName string, err response.Error)) HandlerFunc {
	return func(ctx Context) {
		userID, userName, err := handler(ctx)
		if err != nil {
			ctx.AbortWithError(err)
			return
		}
		ctx.setUserID(userID)
		ctx.setUserName(userName)
	}
}

// RouterGroup 包装gin的RouterGroup
type RouterGroup interface {
	Group(string, ...HandlerFunc) RouterGroup
	IRoutes
}

var _ IRoutes = (*router)(nil)

// IRoutes 包装gin的IRoutes
type IRoutes interface {
	Any(string, ...HandlerFunc)
	GET(string, ...HandlerFunc)
	POST(string, ...HandlerFunc)
	DELETE(string, ...HandlerFunc)
	PATCH(string, ...HandlerFunc)
	PUT(string, ...HandlerFunc)
	OPTIONS(string, ...HandlerFunc)
	HEAD(string, ...HandlerFunc)
}

type router struct {
	group *gin.RouterGroup
}

func (r *router) Group(relativePath string, handlers ...HandlerFunc) RouterGroup {
	group := r.group.Group(relativePath, wrapHandlers(handlers...)...)
	return &router{group: group}
}

func (r *router) Any(relativePath string, handlers ...HandlerFunc) {
	r.group.Any(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) GET(relativePath string, handlers ...HandlerFunc) {
	r.group.GET(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) POST(relativePath string, handlers ...HandlerFunc) {
	r.group.POST(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) DELETE(relativePath string, handlers ...HandlerFunc) {
	r.group.DELETE(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) PATCH(relativePath string, handlers ...HandlerFunc) {
	r.group.PATCH(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) PUT(relativePath string, handlers ...HandlerFunc) {
	r.group.PUT(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) OPTIONS(relativePath string, handlers ...HandlerFunc) {
	r.group.OPTIONS(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) HEAD(relativePath string, handlers ...HandlerFunc) {
	r.group.HEAD(relativePath, wrapHandlers(handlers...)...)
}

func wrapHandlers(handlers ...HandlerFunc) []gin.HandlerFunc {
	funcs := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		handler := handler
		funcs[i] = func(c *gin.Context) {
			ctx := newContext(c)
			defer releaseContext(ctx)

			handler(ctx)
		}
	}
	return funcs
}

var _ Engine = (*engine)(nil)

// Engine http mux
type Engine interface {
	http.Handler
	Group(relativePath string, handlers ...HandlerFunc) RouterGroup
}

type engine struct {
	e         *gin.Engine
	baseGroup *gin.RouterGroup // 全局basePath
}

func (m *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m.e.ServeHTTP(w, req)
}

func (m *engine) Group(relativePath string, handlers ...HandlerFunc) RouterGroup {
	return &router{
		group: m.baseGroup.Group(relativePath, wrapHandlers(handlers...)...),
	}
}

func New(logger *zap.Logger, options ...Option) (Engine, error) {
	if logger == nil {
		return nil, errors.New("logger required")
	}

	gin.SetMode(gin.ReleaseMode)
	mux := &engine{
		e: gin.New(),
	}
	// 全部url以 serverName开头 ： /serverName/metrics
	basePath := "/" + config.Get().Server.ServerName
	mux.baseGroup = mux.e.Group(basePath)

	opt := new(option)
	for _, f := range options {
		f(opt)
	}

	if !opt.disablePProf {
		pprof.RouteRegister(mux.baseGroup) // register pprof to gin
	}

	if !opt.disableSwagger {
		mux.baseGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // register swagger
	}

	if !opt.disablePrometheus {
		mux.baseGroup.GET("/metrics", gin.WrapH(promhttp.Handler())) // register prometheus
	}

	if opt.enableCors {
		mux.baseGroup.Use(cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{
				http.MethodHead,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
			},
			AllowedHeaders:     []string{"*"},
			AllowCredentials:   true,
			OptionsPassthrough: true,
		}))
	}

	// recover两次，防止处理时发生panic，尤其是在OnPanicNotify中。
	mux.baseGroup.Use(func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("got panic", zap.String("panic", fmt.Sprintf("%+v", err)), zap.String("stack", string(debug.Stack())))
			}
		}()

		ctx.Next()
	})

	// Log \ Metrics \ Trace ( 可观察性 中间件)
	mux.baseGroup.Use(func(ctx *gin.Context) {
		ts := time.Now()

		c := newContext(ctx)
		defer releaseContext(c)

		c.init()
		c.setLogger(logger)

		tr := otel.GetTracerProvider().Tracer(config.Get().Server.ServerName + "/HTTPServer")
		jCtx := otel.GetTextMapPropagator().Extract(ctx.Request.Context(), propagation.HeaderCarrier(ctx.Request.Header))
		jCtx, span := tr.Start(jCtx, ctx.Request.URL.String(), trace.WithSpanKind(trace.SpanKindServer))

		ctx.Request.WithContext(jCtx)

		defer func() {
			if err := recover(); err != nil {
				stackInfo := string(debug.Stack())
				logger.Error("got panic", zap.String("panic", fmt.Sprintf("%+v", err)), zap.String("stack", stackInfo))
				c.AbortWithError(response.NewErrorAutoMsg(
					http.StatusInternalServerError,
					response.ServerError,
				))

				if notify := opt.panicNotify; notify != nil {
					notify(c, err, stackInfo)
				}
			}

			if ctx.Writer.Status() == http.StatusNotFound {
				return
			}

			var (
				businessCode int
				abortErr     error
				traceId      string
			)

			if ctx.IsAborted() {
				for i := range ctx.Errors { // gin error
					multierr.AppendInto(&abortErr, ctx.Errors[i])
				}

				if err := c.abortError(); err != nil { // customer err
					multierr.AppendInto(&abortErr, err.GetErr())
					resp := response.NewResponse()
					resp.Code = err.GetBusinessCode()
					resp.Message = err.GetMsg()

					ctx.JSON(err.GetHttpCode(), resp)
				}
			} else {
				payload := c.getPayload()
				resp := response.NewResponse(payload)

				if resp != nil {
					ctx.JSON(http.StatusOK, resp)
				}
			}

			withoutLogPath := []string{
				basePath + "/metrics",
				basePath + "/debug",
				basePath + "/swagger",
				basePath + "/system",
			}
			flag := false
			for _, e := range withoutLogPath {
				if strings.HasPrefix(ctx.Request.URL.Path, e) {
					flag = true
				}
			}
			if flag {
				return
			}

			if opt.recordMetrics != nil {
				uri := c.URI()

				opt.recordMetrics(
					c.RequestContext().Request.Method,
					uri,
					!ctx.IsAborted() && ctx.Writer.Status() == http.StatusOK,
					ctx.Writer.Status(),
					businessCode,
					time.Since(ts).Seconds(),
					traceId,
				)
			}
			if c.getDisableLog() {
				return
			}
			decodedURL, _ := url.QueryUnescape(ctx.Request.URL.RequestURI())

			logger.Info("core-interceptor",
				zap.Any("method", ctx.Request.Method),
				zap.Any("path", decodedURL),
				zap.Any("http_code", ctx.Writer.Status()),
				zap.Any("business_code", businessCode),
				zap.Any("success", !ctx.IsAborted() && ctx.Writer.Status() == http.StatusOK),
				zap.Any("cost_seconds", time.Since(ts).Seconds()),
				zap.Error(abortErr),
			)
		}()

		ctx.Next()
	})

	if opt.enableRate {
		limiter := rate.NewLimiter(rate.Every(time.Second*1), _MaxBurstSize)
		mux.baseGroup.Use(func(ctx *gin.Context) {
			context := newContext(ctx)
			defer releaseContext(context)

			if !limiter.Allow() {
				context.AbortWithError(response.NewErrorAutoMsg(
					http.StatusTooManyRequests,
					response.TooManyRequests,
				))
				return
			}

			ctx.Next()
		})
	}

	system := mux.Group("/system")
	{
		// 健康检查
		system.GET("/health", func(ctx Context) {
			resp := &struct {
				Timestamp time.Time `json:"timestamp"`
				Host      string    `json:"host"`
				Status    string    `json:"status"`
			}{
				Timestamp: time.Now(),
				Host:      ctx.RequestContext().Request.Host,
				Status:    "ok",
			}
			ctx.Payload(resp)
		})
	}

	return mux, nil
}
