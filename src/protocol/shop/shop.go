package shop

//购买物品
type ShopBuyItem struct {
	Count int `json:"count"` //购买个数
	Type  int `json:"type"`  //货币类型，1为钻石，2为金币
	ID    int `json:"id"`    //商品代码
}

//购买物品返回
type ShopBuyItemed struct {
	Honor   int `json:"honor"`   //
	Diamond int `json:"diamond"` //钻石（玩家购买物品后的）
	Code    int `json:"code"`
	Gold    int `json:"gold"` //金币（玩家购买物品后的）
}

type BuyDiamond struct {
	Platform  string `json:"platform"`  //平台
	ProductId int    `json:"productId"` //产品标识
}

type BuyDiamonded struct {
	Code    int    `json:"code"`
	OrderID string `json:"orderId"` //订单标识
}

//补充体力
type BuyStamina struct {
}

//补充体力返回
type BuyStaminaed struct {
	Code int `json:"code"`
}

//购买能量技
type BuyEnergySkill struct {
	Skill int    `json:"skill"` //能量技
}

//购买能量技返回
type BuyEnergySkilled struct {
	Code int `json:"code"`
	Diamond int `json:"diamond"` //钻石（玩家购买物品后的）
	Gold    int `json:"gold"` //金币（玩家购买物品后的）
}

//激活能量技
type ActiveEnergySkill struct {
	Skill int    `json:"skill"` //能量技
}

//激活能量技返回
type ActiveEnergySkilled struct {
	Code int `json:"code"`
	Skill int    `json:"skill"` //能量技
}

//物品出售
type SellItem struct {
	Count    int    `json:"count"`		//物品数量
	ItemID   int    `json:"itemId"`		//物品标识
}

//物品出售返回
type SellItemed struct {
	Code    int    `json:"code"`
}