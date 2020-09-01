package model

import (
	"gameserver/utils/db"
)

// 满能量技能(英雄)   gameserver
type  OwnedEnergySkill struct {
	PlayerId uint `gorm:"type:int(11);column:playerId;unique_index:index_playerId" json:"playerId" `
	SkillId  int  `gorm:"type:int(11);column:skillId;unique_index:index_playerId" json:"skillId" `
	Status   int  `gorm:"type:int(4);column:status;default:0" json:"status" `
}

type OwnedEnergySkills []*OwnedEnergySkill

func (this *OwnedEnergySkills) Get(playerID uint) error {
	return db.Get().Model(this).Where("playerId = ?", playerID).Find(this).Error
}

// 判断是否已经拥有该技能
func (this *OwnedEnergySkill) IsOwnedEnergySkill(playerId uint, skillId int) bool {
	return  db.Get().Model(this).Where("playerId = ? and SkillId = ?", playerId, skillId).Limit(1).Find(this).Error == nil
}

//增加技能
func (this *OwnedEnergySkill) AddEnergySkill() error {
	return db.Get().Create(this).Error
}

//设置技能
func (this *OwnedEnergySkill) SetUsedEnergySkill(playerId uint, skillId int,status int) error {
	//先取消该玩家活动技能
	if err := db.Get().Model(this).Where("playerId = ? ", playerId).Update("status", 0).Error; err != nil {
		return err
	}
	return db.Get().Model(this).Where("playerId = ? and SkillId= ?", playerId, skillId, ).Update("status", status).Error
}
