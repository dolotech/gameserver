package push

import (
	"gameserver/model"
	"gameserver/online"
	proto "gameserver/protocol"
)
//更新角色信息PUSH
type OnUpdatePartners struct {
	PlayerId uint            `json:"playerId"` //角色所属玩家的Id
	Partners []proto.Partner `json:"partners"` //所有角色信息
}

func (this *OnUpdatePartners) Push(player *model.Player, partners *model.Partners) error {
	this.PlayerId = player.PlayerId
	for _, p := range *partners {
		partner := proto.Partner{}
		partner.Exp = player.Exp
		partner.Level = player.Level
		partner.PartnerId = p.PartnerId
		partner.Train.Level = p.Level
		partner.Train.Star = p.PStar
		partner.Train.Point = p.Point
		partner.Equipment.WeaponId = p.WeaponId
		partner.Equipment.Star = p.WStar
		partner.Equipment.Strength = p.Strength
		this.Partners = append(this.Partners, partner)
	}
	online.Get().Push(player.PlayerId, this)
	return nil
}
