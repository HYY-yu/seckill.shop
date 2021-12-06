package model

// Goods 商品表
type Goods struct {
	ID         int    `json:"id" gorm:"primaryKey;column:id;type:int;not null"`
	Name       string `json:"name" gorm:"column:name;type:varchar(255);not null"`
	Desc       string `json:"desc" gorm:"column:desc;type:varchar(500);not null;default:''"`
	Count      int    `json:"count" gorm:"column:count;type:int;not null;default:0"`
	CreateTime int    `json:"create_time" gorm:"column:create_time;type:int;not null"`
}

// GoodsColumns get sql column name.获取数据库列名
var GoodsColumns = struct {
	ID         string
	Name       string
	Desc       string
	Count      string
	CreateTime string
}{
	ID:         "id",
	Name:       "name",
	Desc:       "desc",
	Count:      "count",
	CreateTime: "create_time",
}
