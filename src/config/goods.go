package config

type GoodsData struct {
	Friendly int `json:"friendly"`
	Useocca  int `json:"useocca"`
	Sell     int `json:"sell"`		//商品金额
}
var goodsPool **goodsDataPool

type goodsDataPool map[string]GoodsData

func (this *goodsDataPool) Get(id string) *GoodsData {
	s, ok := (*this)[id]
	if ok {
		return &s
	}
	return nil
}

func Goods() *goodsDataPool {
	return *goodsPool
}
