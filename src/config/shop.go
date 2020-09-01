package config

type ShopData struct {
	Sid      int     `json:"sid"`      //商城表id  (sid 不填代表没有了 )
	Type     int     `json:"type"`     //商品类型 (1 金币商城,用购买金币 2 道具商城。花费钻石买各种道具 3 热卖 4 虚拟物品)
	ShopId   int     `json:"shopId"`   //商品id
	Price    int     `json:"price"`    //商品钻石价格
	Count    int     `json:"count"`    //数量
	Statue   int     `json:"statue"`   //状态	1热卖 2限买 3赠送
	Discount float64 `json:"discount"` //打折  1为100%
}

var shopPool **shopDataPool

type shopDataPool map[string]ShopData

func (this *shopDataPool) Get(id string) *ShopData {
	s, ok := (*this)[id]
	if ok {
		return &s
	}
	return nil
}

func Shop() *shopDataPool {
	return *shopPool
}
