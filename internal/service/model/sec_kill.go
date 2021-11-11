package model

// Shop 商品表
type Shop struct {
	ID         int    `gorm:"primaryKey;column:id;type:int(11);not null"`
	Name       string `gorm:"column:name;type:varchar(255);not null"`
	Desc       string `gorm:"column:desc;type:varchar(500);not null;default:''"`
	Count      int    `gorm:"column:count;type:int(11);not null;default:0"`
	CreateTime int    `gorm:"column:create_time;type:int(11);not null"`
}

// ShopColumns get sql column name.获取数据库列名
var ShopColumns = struct {
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
