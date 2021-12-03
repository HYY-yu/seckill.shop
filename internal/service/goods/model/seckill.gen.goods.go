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

func (obj *_GoodsMgr) Count(count *int64) (tx *gorm.DB) {
	return obj.DB.WithContext(obj.ctx).Model(Goods{}).Count(count)
}

// WithID id获取
func (obj *_GoodsMgr) WithID(id int) Option {
	return optionFunc(func(o *options) { o.query["id"] = id })
}

// WithName name获取
func (obj *_GoodsMgr) WithName(name string) Option {
	return optionFunc(func(o *options) { o.query["name"] = name })
}

// WithDesc desc获取
func (obj *_GoodsMgr) WithDesc(desc string) Option {
	return optionFunc(func(o *options) { o.query["desc"] = desc })
}

// WithCount count获取
func (obj *_GoodsMgr) WithCount(count int) Option {
	return optionFunc(func(o *options) { o.query["count"] = count })
}

// WithCreateTime create_time获取
func (obj *_GoodsMgr) WithCreateTime(createTime int) Option {
	return optionFunc(func(o *options) { o.query["create_time"] = createTime })
}

// GetFromID 通过id获取内容
func (obj *_GoodsMgr) GetFromID(id int) (result Goods, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Goods{}).Where("`id` = ?", id).Find(&result).Error

	return
}

// GetBatchFromID 批量查找
func (obj *_GoodsMgr) GetBatchFromID(ids []int) (results []*Goods, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Goods{}).Where("`id` IN (?)", ids).Find(&results).Error

	return
}

// GetFromName 通过name获取内容
func (obj *_GoodsMgr) GetFromName(name string) (results []*Goods, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Goods{}).Where("`name` = ?", name).Find(&results).Error

	return
}

// GetBatchFromName 批量查找
func (obj *_GoodsMgr) GetBatchFromName(names []string) (results []*Goods, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Goods{}).Where("`name` IN (?)", names).Find(&results).Error

	return
}

// GetFromDesc 通过desc获取内容
func (obj *_GoodsMgr) GetFromDesc(desc string) (results []*Goods, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Goods{}).Where("`desc` = ?", desc).Find(&results).Error

	return
}

// GetBatchFromDesc 批量查找
func (obj *_GoodsMgr) GetBatchFromDesc(descs []string) (results []*Goods, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Goods{}).Where("`desc` IN (?)", descs).Find(&results).Error

	return
}

// GetFromCount 通过count获取内容
func (obj *_GoodsMgr) GetFromCount(count int) (results []*Goods, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Goods{}).Where("`count` = ?", count).Find(&results).Error

	return
}

// GetBatchFromCount 批量查找
func (obj *_GoodsMgr) GetBatchFromCount(counts []int) (results []*Goods, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Goods{}).Where("`count` IN (?)", counts).Find(&results).Error

	return
}

// GetFromCreateTime 通过create_time获取内容
func (obj *_GoodsMgr) GetFromCreateTime(createTime int) (results []*Goods, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Goods{}).Where("`create_time` = ?", createTime).Find(&results).Error

	return
}

// GetBatchFromCreateTime 批量查找
func (obj *_GoodsMgr) GetBatchFromCreateTime(createTimes []int) (results []*Goods, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Goods{}).Where("`create_time` IN (?)", createTimes).Find(&results).Error

	return
}

func (obj *_GoodsMgr) CreateGoods(bean *Goods) (err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Goods{}).Create(bean).Error

	return
}

func (obj *_GoodsMgr) UpdateGoods(bean *Goods) (err error) {
	err = obj.DB.WithContext(obj.ctx).Model(bean).Updates(bean).Error

	return
}

func (obj *_GoodsMgr) DeleteGoods() (err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Goods{}).Error

	return
}
