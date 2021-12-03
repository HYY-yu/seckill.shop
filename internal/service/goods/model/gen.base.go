package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var globalIsRelated bool = true // 全局预加载

// prepare for other
type _BaseMgr struct {
	*gorm.DB
	ctx       context.Context
	cancel    context.CancelFunc
	timeout   time.Duration
	isRelated bool
}

// WithContext set context to db
func (obj *_BaseMgr) WithContext(c context.Context) {
	if obj.DB != nil {
		obj.DB = obj.DB.WithContext(c)
	}
}

func (obj *_BaseMgr) WithSelects(idName string, selects ...string) {
	if len(selects) > 0 {
		if len(idName) > 0 {
			selects = append(selects, idName)
		}
		// 对Select进行去重
		selectMap := make(map[string]int, len(selects))
		for _, e := range selects {
			if _, ok := selectMap[e]; !ok {
				selectMap[e] = 1
			}
		}

		newSelects := make([]string, 0, len(selects))
		for k, _ := range selectMap {
			newSelects = append(newSelects, k)
		}

		obj.DB = obj.DB.Select(newSelects)
	}
}

func (obj *_BaseMgr) WithOmit(omit ...string) {
	if len(omit) > 0 {
		obj.DB = obj.DB.Omit(omit...)
	}
}

func (obj *_BaseMgr) WithOptions(opts ...Option) {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}
	obj.DB = obj.DB.Where(options.query)
}

// SetTimeOut set timeout
func (obj *_BaseMgr) SetTimeOut(timeout time.Duration) {
	obj.ctx, obj.cancel = context.WithTimeout(context.Background(), timeout)
	obj.timeout = timeout
}

// SetCtx set context
func (obj *_BaseMgr) SetCtx(c context.Context) {
	if c != nil {
		obj.ctx = c
	}
}

// GetCtx get context
func (obj *_BaseMgr) GetCtx() context.Context {
	return obj.ctx
}

// Cancel cancel context
func (obj *_BaseMgr) Cancel(c context.Context) {
	obj.cancel()
}

// GetDB get gorm.DB info
func (obj *_BaseMgr) GetDB() *gorm.DB {
	return obj.DB
}

// UpdateDB update gorm.DB info
func (obj *_BaseMgr) UpdateDB(db *gorm.DB) {
	obj.DB = db
}

// GetIsRelated Query foreign key Association.获取是否查询外键关联(gorm.Related)
func (obj *_BaseMgr) GetIsRelated() bool {
	return obj.isRelated
}

// SetIsRelated Query foreign key Association.设置是否查询外键关联(gorm.Related)
func (obj *_BaseMgr) SetIsRelated(b bool) {
	obj.isRelated = b
}

// New new gorm.新gorm,重置条件
func (obj *_BaseMgr) New() {
	obj.DB = obj.NewDB()
}

// NewDB new gorm.新gorm
func (obj *_BaseMgr) NewDB() *gorm.DB {
	return obj.DB.Session(&gorm.Session{NewDB: true, Context: obj.ctx})
}

type options struct {
	query map[string]interface{}
}

// Option overrides behavior of Connect.
type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

// OpenRelated 打开全局预加载
func OpenRelated() {
	globalIsRelated = true
}

// CloseRelated 关闭全局预加载
func CloseRelated() {
	globalIsRelated = true
}

// -------- sql where helper ----------

type CheckWhere func(v interface{}) bool
type DoWhere func(*gorm.DB) *gorm.DB

// CheckWhere 函数 如果返回true，则表明 DoWhere 的查询条件需要加到sql中去
func (w *_BaseMgr) AddWhere(v interface{}, c CheckWhere, d DoWhere) *_BaseMgr {
	if c(v) {
		w.DB = d(w.DB)
	}
	return w
}

func (w *_BaseMgr) Sort(userSort, defaultSort string) *_BaseMgr {
	if len(userSort) > 0 {
		w.DB = w.DB.Order(userSort)
	} else {
		if len(defaultSort) > 0 {
			w.DB = w.DB.Order(defaultSort)
		}
	}
	return w
}

func (w *_BaseMgr) Build() *gorm.DB {
	return w.DB
}
