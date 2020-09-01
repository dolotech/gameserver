package push

//更新钻石PUSH
type OnUpdateBuyDiamond struct {
	Diamond int    `json:"diamond"` //玩家当前拥有钻石数
	OrderId string `json:"orderId"` //订单标识
}
