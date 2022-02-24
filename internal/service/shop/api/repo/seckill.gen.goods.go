package repo

import (
	"context"
	"fmt"

	"github.com/HYY-yu/seckill.shop/internal/service/shop/model"
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

// WithContext set context to db
func (obj *_GoodsMgr) WithContext(c context.Context) *_GoodsMgr {
	if c != nil {
		obj.ctx = c
	}
	return obj
}

func (obj *_GoodsMgr) WithSelects(idName string, selects ...string) *_GoodsMgr {
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
		for k := range selectMap {
			newSelects = append(newSelects, k)
		}

		obj.DB = obj.DB.Select(newSelects)
	}
	return obj
}

func (obj *_GoodsMgr) WithOmit(omit ...string) *_GoodsMgr {
	if len(omit) > 0 {
		obj.DB = obj.DB.Omit(omit...)
	}
	return obj
}

func (obj *_GoodsMgr) WithOptions(opts ...Option) *_GoodsMgr {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}
	obj.DB = obj.DB.Where(options.query)
	return obj
}

// GetTableName get sql table name.获取数据库名字
func (obj *_GoodsMgr) GetTableName() string {
	return "goods"
}

// Reset 重置gorm会话
func (obj *_GoodsMgr) Reset() *_GoodsMgr {
	obj.new()
	return obj
}

// Get 获取
func (obj *_GoodsMgr) Get() (result model.Goods, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Goods{}).Find(&result).Error

	return
}

// Gets 获取批量结果
func (obj *_GoodsMgr) Gets() (results []*model.Goods, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Goods{}).Find(&results).Error

	return
}

func (obj *_GoodsMgr) Count(count *int64) (tx *gorm.DB) {
	return obj.DB.WithContext(obj.ctx).Model(model.Goods{}).Count(count)
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

// WithDeleteTime delete_time获取
func (obj *_GoodsMgr) WithDeleteTime(deleteTime int) Option {
	return optionFunc(func(o *options) { o.query["delete_time"] = deleteTime })
}

// GetFromID 通过id获取内容
func (obj *_GoodsMgr) GetFromID(id int) (result model.Goods, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Goods{}).Where("`id` = ?", id).Find(&result).Error

	return
}

// GetBatchFromID 批量查找
func (obj *_GoodsMgr) GetBatchFromID(ids []int) (results []*model.Goods, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Goods{}).Where("`id` IN (?)", ids).Find(&results).Error

	return
}

// GetFromName 通过name获取内容
func (obj *_GoodsMgr) GetFromName(name string) (result model.Goods, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Goods{}).Where("`name` = ?", name).Find(&result).Error

	return
}

// GetBatchFromName 批量查找
func (obj *_GoodsMgr) GetBatchFromName(names []string) (results []*model.Goods, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Goods{}).Where("`name` IN (?)", names).Find(&results).Error

	return
}

// GetFromDesc 通过desc获取内容
func (obj *_GoodsMgr) GetFromDesc(desc string) (results []*model.Goods, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Goods{}).Where("`desc` = ?", desc).Find(&results).Error

	return
}

// GetBatchFromDesc 批量查找
func (obj *_GoodsMgr) GetBatchFromDesc(descs []string) (results []*model.Goods, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Goods{}).Where("`desc` IN (?)", descs).Find(&results).Error

	return
}

// GetFromCount 通过count获取内容
func (obj *_GoodsMgr) GetFromCount(count int) (results []*model.Goods, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Goods{}).Where("`count` = ?", count).Find(&results).Error

	return
}

// GetBatchFromCount 批量查找
func (obj *_GoodsMgr) GetBatchFromCount(counts []int) (results []*model.Goods, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Goods{}).Where("`count` IN (?)", counts).Find(&results).Error

	return
}

// GetFromCreateTime 通过create_time获取内容
func (obj *_GoodsMgr) GetFromCreateTime(createTime int) (results []*model.Goods, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Goods{}).Where("`create_time` = ?", createTime).Find(&results).Error

	return
}

// GetBatchFromCreateTime 批量查找
func (obj *_GoodsMgr) GetBatchFromCreateTime(createTimes []int) (results []*model.Goods, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Goods{}).Where("`create_time` IN (?)", createTimes).Find(&results).Error

	return
}

// GetFromDeleteTime 通过delete_time获取内容
func (obj *_GoodsMgr) GetFromDeleteTime(deleteTime int) (results []*model.Goods, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Goods{}).Where("`delete_time` = ?", deleteTime).Find(&results).Error

	return
}

// GetBatchFromDeleteTime 批量查找
func (obj *_GoodsMgr) GetBatchFromDeleteTime(deleteTimes []int) (results []*model.Goods, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Goods{}).Where("`delete_time` IN (?)", deleteTimes).Find(&results).Error

	return
}

func (obj *_GoodsMgr) CreateGoods(bean *model.Goods) (err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Goods{}).Create(bean).Error

	return
}

func (obj *_GoodsMgr) UpdateGoods(bean *model.Goods) (err error) {
	err = obj.DB.WithContext(obj.ctx).Model(bean).Updates(bean).Error

	return
}

func (obj *_GoodsMgr) DeleteGoods(bean *model.Goods) (err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Goods{}).Delete(bean).Error

	return
}
