package model

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type _ShopMgr struct {
	*_BaseMgr
}

// ShopMgr open func
func ShopMgr(db *gorm.DB) *_ShopMgr {
	if db == nil {
		panic(fmt.Errorf("ShopMgr need init by db"))
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &_ShopMgr{_BaseMgr: &_BaseMgr{DB: db.Table("shop"), isRelated: globalIsRelated, ctx: ctx, cancel: cancel, timeout: -1}}
}

// GetTableName get sql table name.获取数据库名字
func (obj *_ShopMgr) GetTableName() string {
	return "shop"
}

// Reset 重置gorm会话
func (obj *_ShopMgr) Reset() *_ShopMgr {
	obj.New()
	return obj
}

// Get 获取
func (obj *_ShopMgr) Get() (result Shop, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Shop{}).Find(&result).Error

	return
}

// Gets 获取批量结果
func (obj *_ShopMgr) Gets() (results []*Shop, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Shop{}).Find(&results).Error

	return
}

////////////////////////////////// gorm replace /////////////////////////////////
func (obj *_ShopMgr) Count(count *int64) (tx *gorm.DB) {
	return obj.DB.WithContext(obj.ctx).Model(Shop{}).Count(count)
}

//////////////////////////////////////////////////////////////////////////////////

//////////////////////////option case ////////////////////////////////////////////

// WithID id获取
func (obj *_ShopMgr) WithID(id int) Option {
	return optionFunc(func(o *options) { o.query["id"] = id })
}

// WithName name获取
func (obj *_ShopMgr) WithName(name string) Option {
	return optionFunc(func(o *options) { o.query["name"] = name })
}

// WithDesc desc获取
func (obj *_ShopMgr) WithDesc(desc string) Option {
	return optionFunc(func(o *options) { o.query["desc"] = desc })
}

// WithCount count获取
func (obj *_ShopMgr) WithCount(count int) Option {
	return optionFunc(func(o *options) { o.query["count"] = count })
}

// WithCreateTime create_time获取
func (obj *_ShopMgr) WithCreateTime(createTime int) Option {
	return optionFunc(func(o *options) { o.query["create_time"] = createTime })
}

// GetByOption 功能选项模式获取
func (obj *_ShopMgr) GetByOption(opts ...Option) (result Shop, err error) {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}

	err = obj.DB.WithContext(obj.ctx).Model(Shop{}).Where(options.query).Find(&result).Error

	return
}

// GetByOptions 批量功能选项模式获取
func (obj *_ShopMgr) GetByOptions(opts ...Option) (results []*Shop, err error) {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}

	err = obj.DB.WithContext(obj.ctx).Model(Shop{}).Where(options.query).Find(&results).Error

	return
}

//////////////////////////enume case ////////////////////////////////////////////

// GetFromID 通过id获取内容
func (obj *_ShopMgr) GetFromID(id int) (result Shop, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Shop{}).Where("`id` = ?", id).Find(&result).Error

	return
}

// GetBatchFromID 批量查找
func (obj *_ShopMgr) GetBatchFromID(ids []int) (results []*Shop, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Shop{}).Where("`id` IN (?)", ids).Find(&results).Error

	return
}

// GetFromName 通过name获取内容
func (obj *_ShopMgr) GetFromName(name string) (results []*Shop, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Shop{}).Where("`name` = ?", name).Find(&results).Error

	return
}

// GetBatchFromName 批量查找
func (obj *_ShopMgr) GetBatchFromName(names []string) (results []*Shop, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Shop{}).Where("`name` IN (?)", names).Find(&results).Error

	return
}

// GetFromDesc 通过desc获取内容
func (obj *_ShopMgr) GetFromDesc(desc string) (results []*Shop, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Shop{}).Where("`desc` = ?", desc).Find(&results).Error

	return
}

// GetBatchFromDesc 批量查找
func (obj *_ShopMgr) GetBatchFromDesc(descs []string) (results []*Shop, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Shop{}).Where("`desc` IN (?)", descs).Find(&results).Error

	return
}

// GetFromCount 通过count获取内容
func (obj *_ShopMgr) GetFromCount(count int) (results []*Shop, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Shop{}).Where("`count` = ?", count).Find(&results).Error

	return
}

// GetBatchFromCount 批量查找
func (obj *_ShopMgr) GetBatchFromCount(counts []int) (results []*Shop, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Shop{}).Where("`count` IN (?)", counts).Find(&results).Error

	return
}

// GetFromCreateTime 通过create_time获取内容
func (obj *_ShopMgr) GetFromCreateTime(createTime int) (results []*Shop, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Shop{}).Where("`create_time` = ?", createTime).Find(&results).Error

	return
}

// GetBatchFromCreateTime 批量查找
func (obj *_ShopMgr) GetBatchFromCreateTime(createTimes []int) (results []*Shop, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Shop{}).Where("`create_time` IN (?)", createTimes).Find(&results).Error

	return
}

//////////////////////////primary index case ////////////////////////////////////////////

// FetchByPrimaryKey primary or index 获取唯一内容
func (obj *_ShopMgr) FetchByPrimaryKey(id int) (result Shop, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Shop{}).Where("`id` = ?", id).Find(&result).Error

	return
}
