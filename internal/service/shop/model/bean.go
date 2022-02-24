package model

type GoodsListResp struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	Count      int    `json:"count"`
	CreateTime int    `json:"create_time"`
}

type GoodsAdd struct {
	Name  string `json:"name" v:"required|length:3,10#请输入商品名称|商品名称长度:min-:max"`
	Desc  string `json:"desc" v:"max-length:50#商品描述请控制在:max个字符"`
	Count uint   `json:"count" v:"min:1#请输入商品库存数"`
}

type GoodsUpdate struct {
	Id    int     `json:"id" v:"required"`
	Name  *string `json:"name" v:"length:3,10"`
	Desc  *string `json:"desc" v:"max-length:50"`
	Count *uint   `json:"count" v:"min:1"`
}
