package model

import (
	"gameserver/utils/db"
)

// 关卡记录   gameserver
type DungeonRecord struct {
	PlayerId  uint `gorm:"type:int(11);column:playerId;unique_index:index_playerId" json:"playerId" `
	DungeonId int  `gorm:"type:int(11);column:dungeonId;unique_index:index_playerId" json:"dungeonId" `
	Star      int  `gorm:"type:int(4);column:star" json:"star" `
}
type DungeonRecords []*DungeonRecord

func (this *DungeonRecords) Get(playerID uint) error {
	return db.Get().Model(this).Where("playerId = ?", playerID).Find(this).Error
}

func (this *DungeonRecords) GetTotalStar(playerId uint) int {
	db.Get().Model(this).Where("playerId = ?", playerId).Find(this)
	totalStar := 0
	for _, record := range *this {
		totalStar += record.Star
	}
	return totalStar
}

func (this *DungeonRecord) SaveRecord(dungeonId int, playerId uint, star int) error {
	err := db.Get().Model(this).Where("playerId = ? and dungeonId = ?", playerId, dungeonId).Find(this).Error
	if err != nil {
		this.PlayerId = playerId
		this.DungeonId = dungeonId
		this.Star = star
		return db.Get().Create(this).Error
	}
	if this.Star >= star {
		return nil
	}
	this.Star = star
	return db.Get().Model(this).Updates(map[string]interface{}{"star": star}).Error
}
