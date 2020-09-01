package config

type DiamondData struct {
	Platform  string `json:"platform"`
	ProductId int    `json:"productId"`
	PayCode   string `json:"payCode"`
	Count     int    `json:"count"`
	Price     int    `json:"price"`
}

func (this *DiamondDataPool) Get(platform string, productId int) *DiamondData {
	for _, v := range **diamondDataPool {
		if v.Platform == platform && productId == v.ProductId {
			return &v
		}
	}
	return nil
}

var diamondDataPool **DiamondDataPool

type DiamondDataPool map[string]DiamondData

func Diamond() *DiamondDataPool {
	return *diamondDataPool
}
