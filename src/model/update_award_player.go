package model

import "gameserver/utils/db"

// 奖励相关的玩家属性
func (this *Player) GetAwardAttr() error {
	return db.Get().Model(this).Select("playerId,level,pvpScore,diamond,honor,power,exp,gold,stamina").Limit(1).Find(this).Error
}

// 保存奖励时需要涉及到的玩家属性
func (this *Player) SetAwardAttr() error {
	return db.Get().Model(this).Updates(map[string]interface{}{
		"exp":      this.Exp,
		"gold":     this.Gold,
		"diamond":  this.Diamond,
		"stamina":  this.Stamina,
		"honor":    this.Honor,
		"pvpScore": this.PvpScore,
		"level":    this.Level,
		"power":    this.Power,
	}).Error
}

// 推送OnUpdatePlayer 协议用的的玩家字段
func (this *Player) GetOnUpdatePlayer() error {
	return db.Get().Table("players").Select("lastRecvStaminaTime,gold,exp,pvpRank,honor,pvpScore,stamina,diamond,usedEnergySkill,playerId,resType,vip,level, power,avatar,nickname").Limit(1).Find(this).Error
}
