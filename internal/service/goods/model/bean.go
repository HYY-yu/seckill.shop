package model

type GoodsListResp struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	Count      int    `json:"count"`
	CreateTime int    `json:"create_time"`
}

type GoodsAdd struct {
	Name  string `json:"name" v:""`
	Desc  string `json:"desc"`
	Count int    `json:"count"`
}
