package core

import (
	"bytes"
	stdContext "context"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/spf13/cast"
	"go.uber.org/zap"

	"github.com/HYY-yu/seckill.pkg/pkg/response"
)

type HandlerFunc func(c Context)

// 内存缓存
const (
	_Logger     = "_logger_"
	_Response   = "_response_"
	_UserID     = "_user_id_"
	_UserName   = "_user_name_"
	_DisableLog = "_disable_log_"
)

var contextPool = &sync.Pool{
	New: func() interface{} {
		return new(context)
	},
}

func newContext(ctx *gin.Context) Context {
	context := contextPool.Get().(*context)
	context.ctx = ctx
	return context
}

func releaseContext(ctx Context) {
	c := ctx.(*context)
	c.ctx = nil
	contextPool.Put(c)
}

var _ Context = (*context)(nil)

// Context 封装 gin.Context
type Context interface {
	// ShouldBindForm 同时反序列化 querystring 和 postForm;
	// 当 querystring 和 postForm 存在相同字段时，postForm 优先使用。
	// tag: `form:"xxx"`
	ShouldBindForm(obj interface{}) error

	// ShouldBindJSON 反序列化 postJson
	// tag: `json:"xxx"`
	ShouldBindJSON(obj interface{}) error

	// ShouldBindURI 反序列化 path 参数(如路由路径为 /user/:name)
	// tag: `uri:"xxx"`
	ShouldBindURI(obj interface{}) error

	// Header 获取 Header 对象
	Header() http.Header
	// GetHeader 获取 Header
	GetHeader(key string) string
	// SetHeader 设置 Header
	SetHeader(key, value string)

	// URI unescape后的uri
	URI() string
	// RequestData 获取请求体（可多次读取）
	RequestData() []byte

	// Redirect 重定向
	Redirect(code int, location string)

	// Payload 正确返回
	Payload(payload interface{})
	getResponse() interface{}

	// HTML 返回界面
	HTML(name string, obj interface{})

	// AbortWithError 错误返回
	AbortWithError(err error)

	// Logger 获取 Logger 对象
	Logger() *zap.Logger
	setLogger(logger *zap.Logger)

	DisableLog(flag bool)
	getDisableLog() bool

	// UserID 获取 UserID
	UserID() int64
	setUserID(userID int64)

	// UserName 获取 UserName
	UserName() string
	setUserName(userName string)

	// RequestContext 获取GIN的 context
	RequestContext() *gin.Context
	// SvcContext 给下层用的context
	SvcContext() SvcContext
}

type context struct {
	ctx *gin.Context
}

// ShouldBindForm 同时反序列化 querystring 和 postForm;
// 当 querystring 和 postForm 存在相同字段时，postForm 优先使用。
// tag: `form:"xxx"`
func (c *context) ShouldBindForm(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.Form)
}

// ShouldBindJSON 反序列化postJson
// tag: `json:"xxx"`
func (c *context) ShouldBindJSON(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.JSON)
}

// ShouldBindURI 反序列化path参数(如路由路径为 /user/:name)
// tag: `uri:"xxx"`
func (c *context) ShouldBindURI(obj interface{}) error {
	return c.ctx.ShouldBindUri(obj)
}

// Redirect 重定向
func (c *context) Redirect(code int, location string) {
	c.ctx.Redirect(code, location)
}

func (c *context) Logger() *zap.Logger {
	logger, ok := c.ctx.Get(_Logger)
	if !ok {
		return nil
	}
	return logger.(*zap.Logger)
}

func (c *context) setLogger(logger *zap.Logger) {
	c.ctx.Set(_Logger, logger)
}

func (c *context) Payload(payload interface{}) {
	resp := response.NewResponse(payload)

	if _, exist := c.ctx.Get(_Response); exist {
		return
	}
	c.ctx.JSON(http.StatusOK, resp)
	c.ctx.Set(_Response, resp)
}

func (c *context) getResponse() interface{} {
	if resp, ok := c.ctx.Get(_Response); ok != false {
		return resp
	}
	return nil
}

func (c *context) HTML(name string, obj interface{}) {
	c.ctx.HTML(200, name+".html", obj)
}

func (c *context) AbortWithError(err error) {
	if err != nil {
		errResp := response.NewErrorAutoMsg(http.StatusInternalServerError, response.ServerError)
		if v, ok := err.(response.Error); ok {
			errResp = v
		} else {
			_ = errResp.WithErr(err)
		}

		if errResp.GetErr() != nil {
			c.Logger().Error("server error ...", zap.Error(errResp.GetErr()))
		}

		httpCode := errResp.GetHttpCode()
		if httpCode == 0 {
			httpCode = http.StatusInternalServerError
		}

		resp := response.NewResponse()
		resp.Code = errResp.GetBusinessCode()
		resp.Message = errResp.GetMsg()

		c.ctx.AbortWithStatus(httpCode)
		c.ctx.Set(_Response, resp)
		c.ctx.JSON(httpCode, resp)
	}
}

func (c *context) Header() http.Header {
	header := c.ctx.Request.Header

	clone := make(http.Header, len(header))
	for k, v := range header {
		value := make([]string, len(v))
		copy(value, v)

		clone[k] = value
	}
	return clone
}

func (c *context) GetHeader(key string) string {
	return c.ctx.GetHeader(key)
}

func (c *context) SetHeader(key, value string) {
	c.ctx.Header(key, value)
}

func (c *context) UserID() int64 {
	val, ok := c.ctx.Get(_UserID)
	if !ok {
		return 0
	}

	return val.(int64)
}

func (c *context) setUserID(userID int64) {
	c.ctx.Set(_UserID, userID)
}

func (c *context) UserName() string {
	val, ok := c.ctx.Get(_UserName)
	if !ok {
		return ""
	}

	return val.(string)
}

func (c *context) setUserName(userName string) {
	c.ctx.Set(_UserName, userName)
}

func (c *context) RequestData() []byte {
	rawData, _ := c.ctx.GetRawData()
	c.ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawData)) // re-construct req body
	return rawData
}

// RequestContext 获取GIN的Context
func (c *context) RequestContext() *gin.Context {
	return c.ctx
}

// SvcContext 传给下层用的Context，精简去掉Request、Response等信息
// 只保留以下信息
type SvcContext interface {
	UserId() int64
	UserName() string
	Context() stdContext.Context
	Logger() *zap.Logger
}

type svcContext struct {
	ctx    stdContext.Context
	logger *zap.Logger
}

func (s *svcContext) UserId() int64 {
	return cast.ToInt64(s.ctx.Value(_UserID))
}

func (s *svcContext) UserName() string {
	return cast.ToString(s.ctx.Value(_UserName))
}

func (s *svcContext) Context() stdContext.Context {
	return s.ctx
}

func (s *svcContext) Logger() *zap.Logger {
	return s.logger
}

func (c *context) SvcContext() SvcContext {
	ctx := c.RequestContext().Request.Context()

	// 用户信息设置进去
	ctx = stdContext.WithValue(ctx, _UserID, c.UserID())
	ctx = stdContext.WithValue(ctx, _UserName, c.UserName())

	return &svcContext{
		ctx:    ctx,
		logger: c.Logger(),
	}
}

// URI unescape后的uri
func (c *context) URI() string {
	uri, _ := url.QueryUnescape(c.ctx.Request.URL.RequestURI())
	return uri
}

func (c *context) DisableLog(flag bool) {
	c.ctx.Set(_DisableLog, flag)
}

func (c *context) getDisableLog() bool {
	val, ok := c.ctx.Get(_DisableLog)
	if !ok {
		return false
	}

	return val.(bool)
}
