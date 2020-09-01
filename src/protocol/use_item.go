package protocol

//战斗外使用物品
type UseItem struct {
	ItemId   int `json:"itemId"`   //物品标识
	Key      int `json:"key"`      // 格子ID之类的（后端不用处理，直接传给前端 ） ANJUN
	BattleId int `json:"battleId"` //战斗标识
}

//战斗外使用物品返回
type UseItemed struct {
	Key    int `json:"key"`
	Code   int `json:"code"`
	ItemId int `json:"itemId"` //物品标识
}
