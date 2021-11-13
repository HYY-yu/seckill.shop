package core

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"

	"github.com/HYY-yu/seckill/pkg/response"
)

type HandlerFunc func(c Context)

// 内存缓存
const (
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
	// ShouldBindForm 同时反序列化 querystring 和 postform;
	// 当 querystring 和 postform 存在相同字段时，postform 优先使用。
	// tag: `form:"xxx"`
	ShouldBindForm(obj interface{}) error

	// ShouldBindJSON 反序列化 postjson
	// tag: `json:"xxx"`
	ShouldBindJSON(obj interface{}) error

	// ShouldBindURI 反序列化 path 参数(如路由路径为 /user/:name)
	// tag: `uri:"xxx"`
	ShouldBindURI(obj interface{}) error

	// Redirect 重定向
	Redirect(code int, location string)

	// Logger 获取 Logger 对象
	Logger() *zap.Logger
	setLogger(logger *zap.Logger)

	// Payload 正确返回
	Payload(payload interface{})
	getResponse() interface{}

	// HTML 返回界面
	HTML(name string, obj interface{})

	// AbortWithError 错误返回
	AbortWithError(err error)

	DisableLog(flag bool)
	getDisableLog() bool

	// Header 获取 Header 对象
	Header() http.Header
	// GetHeader 获取 Header
	GetHeader(key string) string
	// SetHeader 设置 Header
	SetHeader(key, value string)

	// UserID 获取 UserID
	UserID() int64
	setUserID(userID int64)

	// UserName 获取 UserName
	UserName() string
	setUserName(userName string)

	// RequestInputParams 获取所有参数
	RequestInputParams() url.Values
	// RequestPostFormParams  获取 PostForm 参数
	RequestPostFormParams() url.Values

	// RequestContext 获取GIN的 context
	RequestContext() *gin.Context
	// URI unescape后的uri
	URI() string
}

type context struct {
	ctx    *gin.Context
	logger *zap.Logger
}

// ShouldBindForm 同时反序列化querystring和postform;
// 当querystring和postform存在相同字段时，postform优先使用。
// tag: `form:"xxx"`
func (c *context) ShouldBindForm(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.Form)
}

// ShouldBindJSON 反序列化postjson
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
	return c.logger
}

func (c *context) setLogger(logger *zap.Logger) {
	c.logger = logger
}

func (c *context) Payload(payload interface{}) {
	resp := response.NewResponse(payload)

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
			errResp.WithErr(err)
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

// RequestInputParams 获取所有参数
func (c *context) RequestInputParams() url.Values {
	_ = c.ctx.Request.ParseForm()
	return c.ctx.Request.Form
}

// RequestPostFormParams 获取 PostForm 参数
func (c *context) RequestPostFormParams() url.Values {
	_ = c.ctx.Request.ParseForm()
	return c.ctx.Request.PostForm
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
