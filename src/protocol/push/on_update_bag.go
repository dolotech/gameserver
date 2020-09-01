package push

import (
	"gameserver/model"
	"gameserver/online"
	"gameserver/protocol"
)

//更新背包PUSH
type OnUpdateBag struct {
	PlayerId uint            `json:"playerId"` //Id
	Items    []protocol.Item `json:"items"`    //背包物品信息
}

func (this *OnUpdateBag) Push(playerId uint,after ...bool) {
	p := model.Items{}
	this.PlayerId = playerId
	if p.Get(playerId) == nil {
		for _, v := range p {
			item := protocol.Item{ID: v.ItemId, Count: v.Count}
			this.Items = append(this.Items, item)
		}
		online.Get().Push(playerId, this,after...)
	}
}
