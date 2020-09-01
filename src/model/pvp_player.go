package model

import (
	"gameserver/config"
	"gameserver/utils/db"
	"github.com/jinzhu/gorm"
)

//对手信息
func (this *Player) GetPvpOpp() error {
	return db.Get().Model(this).Select("usedEnergySkill,playerId,resType,vip,level, power,diamond,avatar,nickname," +
		"pvpRank," +
		"pvpScore," +
		"lastArenaAwardIndex," +
		"lastArenaAwardTime," +
		"continuousArenaPKWins," +
		"dayArenaPKTimes," +
		"arenaPKLoses," +
		"arenaPKWins").Limit(1).Find(this).Error
}

func (this *Player) AddPvpScore(value int) error {
	return db.Get().Model(this).UpdateColumn("pvpScore", gorm.Expr("pvpScore + ?", value)).Error
}

func (this *Player) AddPvpWins() error {
	return db.Get().Model(this).UpdateColumn("arenaPKWins", gorm.Expr("arenaPKWins + ?", 1)).Error
}

func (this *Player) AddPvpLoses() error {
	return db.Get().Model(this).UpdateColumn("arenaPKLoses", gorm.Expr("arenaPKLoses + ?", 1)).Error
}

func (this *Player) GetPvpAwardTime() error {
	return db.Get().Model(this).Select("lastArenaAwardTime").Limit(1).Find(&this).Error
}

func (this *Player) RecvPvpAward(t int) error {
	this.LeftWeeklyAwardTime = t
	return db.Get().Model(this).Updates(map[string]interface{}{
		"lastArenaAwardTime": t,
	}).Error
}

func (this *Player) ResetPvpPKTimes() error {
	return db.Get().Model(this).Updates(map[string]interface{}{
		"dayArenaPKTimes": config.Param().PvpMaxNum,
	}).Error
}

func (this *Player) AddDayPvpPKTimes(value int) error {
	if value < 0 {
		db.Get().Model(this).Select("dayArenaPKTimes").Limit(1).Find(this)
		if this.DayArenaPKTimes < value {
			return enoughDayWeeklyPKTimes
		}
	}
	return db.Get().Model(this).UpdateColumn("dayArenaPKTimes", gorm.Expr("dayArenaPKTimes + ?", value)).Error
}

func (this *Player) UpdateContinuousArenaPKWins(value int) error {
	return db.Get().Model(this).Update("continuousArenaPKWins", value).Error
}
func (this *Player) AddContinuousArenaPKWins(value int) error {
	return db.Get().Model(this).UpdateColumn("continuousArenaPKWins", gorm.Expr("continuousArenaPKWins + ?", value)).Error
}
func (this *Player) GetDiamondAndArenaPK() error {
	return db.Get().Model(this).Select("diamond,dayArenaPKTimes").Find(this).Error
}

// 查找一个战力相近的玩家这个玩家不能是自己
/*func (this *Player) GetByPowerBetween(playerId uint, lower, upper int) error {
	return db.Get().Raw("select usedEnergySkill,playerId,resType,vip,level, power,diamond,avatar,nickname,"+
		"pvpRank,"+
		"pvpScore,"+
		"lastArenaAwardIndex,"+
		"lastArenaAwardTime,"+
		"continuousArenaPKWins,"+
		"dayArenaPKTimes,"+
		"arenaPKLoses,"+
		"arenaPKWins" +
		"from players where  (playerId != ?) and  (power between ? and ? )", playerId, lower, upper).Limit(1).Scan(this).Error
}*/

func (this *Player) GetByPowerBetween(playerId uint, lower, upper int) error {
	return db.Get().Raw("select usedEnergySkill,playerId,resType,vip,level, power,diamond,avatar,nickname,"+
		"pvpRank,"+
		"pvpScore,"+
		"lastArenaAwardIndex,"+
		"lastArenaAwardTime,"+
		"continuousArenaPKWins,"+
		"dayArenaPKTimes,"+
		"arenaPKLoses,"+
		"arenaPKWins " +
		"from players where  (playerId != ?) and  (power between ? and ? ) ", playerId, lower, upper).Limit(1).Scan(this).Error
}

/*
func (this *Players) RankNear(self, count, lower, upper int) error {
	return db.Get().Raw("select usedEnergySkill,playerId,resType,vip,level, power,avatar,nickname,"+
		"continuousWeeklyPKWins,"+
		"weeklyRank,"+
		"weeklyPKWins,"+
		"weeklyPKLoses "+
		"from players where  (weeklyRank != ?) and  (weeklyRank between ? and ? ) and playerId in (select max(playerId) from players group by  weeklyRank) ", self, lower, upper).Limit(count).Scan(this).Error
}*/