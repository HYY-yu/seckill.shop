package model

import (
	"context"
	"fmt"
	"gorm.io/gorm"
)

type _GoodsMgr struct {
	*_BaseMgr
}

// GoodsMgr open func
func GoodsMgr(db *gorm.DB) *_GoodsMgr {
	if db == nil {
		panic(fmt.Errorf("GoodsMgr need init by db"))
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &_GoodsMgr{_BaseMgr: &_BaseMgr{DB: db.Table("goods"), isRelated: globalIsRelated, ctx: ctx, cancel: cancel, timeout: -1}}
}

// GetTableName get sql table name.获取数据库名字
func (obj *_GoodsMgr) GetTableName() string {
	return "goods"
}

// Reset 重置gorm会话
func (obj *_GoodsMgr) Reset() *_GoodsMgr {
	obj.New()
	return obj
}

// Get 获取
func (obj *_GoodsMgr) Get() (result Goods, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Goods{}).Find(&result).Error

	return
}

// Gets 获取批量结果
func (obj *_GoodsMgr) Gets() (results []*Goods, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Goods{}).Find(&results).Error

	return
}

////////////////////////////////// gorm replace /////////////////////////////////
func (obj *_GoodsMgr) Count(count *int64) (tx *gorm.DB) {
	return obj.DB.WithContext(obj.ctx).Model(Goods{}).Count(count)
}

//////////////////////////////////////////////////////////////////////////////////

//////////////////////////option case ////////////////////////////////////////////

// GetByOption 功能选项模式获取
func (obj *_GoodsMgr) GetByOption(opts ...Option) (result Goods, err error) {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}

	err = obj.DB.WithContext(obj.ctx).Model(Goods{}).Where(options.query).Find(&result).Error

	return
}

// GetByOptions 批量功能选项模式获取
func (obj *_GoodsMgr) GetByOptions(opts ...Option) (results []*Goods, err error) {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}

	err = obj.DB.WithContext(obj.ctx).Model(Goods{}).Where(options.query).Find(&results).Error

	return
}

//////////////////////////enume case ////////////////////////////////////////////

//////////////////////////primary index case ////////////////////////////////////////////
