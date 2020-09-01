package model

// 订单   gameserver
type Order struct {
	POrderId  int    `gorm:"type:int(20);column:pOrderId" json:"pOrderId" `
	Platform  string `gorm:"varchar(30);column:platform" json:"platform" `
	ProductId int    `gorm:"type:int(11);column:productId" json:"productId" `
	PlayerId  int    `gorm:"type:int(20);column:playerId" json:"playerId" `
	Notified  bool   `gorm:"int(4);column:notified" json:"notified" `
	Resolved  bool   `gorm:"int(4);column:resolved" json:"resolved" `
}
