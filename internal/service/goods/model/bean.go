package model

type GoodsListResp struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	Count      int    `json:"count"`
	CreateTime int    `json:"create_time"`
}

type GoodsAdd struct {
	Name  string `json:"name" v:"required#请输入商品名称|length:3,10"`
	Desc  string `json:"desc" v:"max-length:50"`
	Count int    `json:"count" v:"required#请输入商品库存数"`
}
