package model

// Goods 商品表
type Goods struct {
	ID         int    `gorm:"primaryKey;column:id;type:int;not null"`
	Name       string `gorm:"unique;column:name;type:varchar(255);not null"`
	Desc       string `gorm:"column:desc;type:varchar(500);not null;default:''"`
	Count      int    `gorm:"column:count;type:int;not null;default:0"`
	CreateTime int    `gorm:"column:create_time;type:int;not null"`
	DeleteTime int    `gorm:"column:delete_time;type:int;not null;default:0"`
}

// GoodsColumns get sql column name.获取数据库列名
var GoodsColumns = struct {
	ID         string
	Name       string
	Desc       string
	Count      string
	CreateTime string
	DeleteTime string
}{
	ID:         "id",
	Name:       "name",
	Desc:       "desc",
	Count:      "count",
	CreateTime: "create_time",
	DeleteTime: "delete_time",
}
