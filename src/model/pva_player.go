package model

import (
	"errors"
	"gameserver/config"
	"gameserver/utils/db"
	"gameserver/utils/log"
	"github.com/jinzhu/gorm"
	"time"
)

func (this *Players) Rank(self, count int) error {
	return db.Get().Raw("select usedEnergySkill,playerId,resType,vip,level, power,avatar,nickname,"+
		"continuousWeeklyPKWins,"+
		"weeklyRank,"+
		"weeklyPKWins,"+
		"weeklyPKLoses "+
		"from players  where (weeklyRank != ?) and  resType !=0 and playerId in (select max(playerId) from players group by  weeklyRank) order by weeklyRank asc", self).Limit(count).Scan(this).Error
}

func (this *Players) RankNear(self, count, lower, upper int) error {
	return db.Get().Raw("select usedEnergySkill,playerId,resType,vip,level, power,avatar,nickname,"+
		"continuousWeeklyPKWins,"+
		"weeklyRank,"+
		"weeklyPKWins,"+
		"weeklyPKLoses "+
		"from players where  (weeklyRank != ?) and  (weeklyRank between ? and ? ) and playerId in (select max(playerId) from players group by  weeklyRank) ", self, lower, upper).Limit(count).Scan(this).Error
}

//对手信息
func (this *Player) GetWeeklyOpp() error {
	return db.Get().Table("players").Select("usedEnergySkill,playerId,resType,vip,level, power,avatar,nickname," +
		"continuousWeeklyPKWins," +
		"weeklyRank," +
		"weeklyPKWins," +
		"weeklyPKLoses").Limit(1).Find(this).Error
}

// 玩家关于竞技塔的信息
func (this *Player) GetWeeklyPlayer() error {
	return db.Get().Table("players").Select("usedEnergySkill,playerId,resType,vip,level, power,avatar,nickname," +
		"continuousWeeklyPKWins," +
		"weeklyRank," +
		"weeklyPKWins," +
		"buyWeeklyPKTimes," +
		"dayWeeklyPKTimes," +
		"weeklyLeftTime," +
		"lastWeeklyRank," +
		"leftWeeklyAwardTime," +
		"weeklyPKLoses").Limit(1).Find(this).Error
}

func weeklyPKTimesCheck(pk, lastTime int) (add, next int) {
	cd := config.Param().PvaRecover
	max := config.Param().PvaMaxNum
	now := int(time.Now().Unix())
	next = now
	if lastTime > 0 {
		add = (now - lastTime) / cd
		next -= (now - lastTime) % cd
	}
	if add+pk > max { //只加到满为止
		add = max - pk
		next = now
	}
	if pk >= max { //原本就处于最大值，要更新一下领取时间
		add = 0
		next = now
	}
	return
}

var enoughDayWeeklyPKTimes = errors.New("not enough stamina")

func (this *Player) UpdateDayWeeklyPKTimes(value int) error {
	db.Get().Model(this).Select("playerId,weeklyLeftTime,dayWeeklyPKTimes").Limit(1).Find(this)
	add, next := weeklyPKTimesCheck(this.DayWeeklyPKTimes, this.WeeklyLeftTime)
	if value < 0 && this.DayWeeklyPKTimes+value+add < 0 {
		return enoughDayWeeklyPKTimes
	}

	this.DayWeeklyPKTimes += (value + add)
	this.WeeklyLeftTime = next

	return db.Get().Model(this).Updates(map[string]interface{}{
		"weeklyLeftTime":   this.WeeklyLeftTime,
		"dayWeeklyPKTimes": this.DayWeeklyPKTimes,
	}).Error
}

type maxStruct struct {
	Max int
}

func (this *Player) GetMaxRank() int {
	var result maxStruct
	err := db.Get().Model(this).Select(" MAX(players.weeklyRank) AS max ").Scan(&result).Error
	if err != nil {
		log.Error(err)
	}
	return result.Max
}

func (this *Player) GetWeeklyAwardTime() error {
	return db.Get().Model(this).Select("leftWeeklyAwardTime").Limit(1).Find(&this).Error
}

func (this *Player) RecvWeeklyAward(t int) error {
	this.LeftWeeklyAwardTime = t
	return db.Get().Model(this).Updates(map[string]interface{}{
		"leftWeeklyAwardTime": t,
	}).Error
}

func (this *Player) AddWeeklyPKWins() error {
	return db.Get().Model(this).UpdateColumn("weeklyPKWins", gorm.Expr("weeklyPKWins + ?", 1)).Error
}

func (this *Player) AddWeeklyPKLoses() error {
	return db.Get().Model(this).UpdateColumn("weeklyPKLoses", gorm.Expr("weeklyPKLoses + ?", 1)).Error
}

func (this *Player) UpdateRank() error {
	return db.Get().Model(this).Update("weeklyRank", this.WeeklyRank).Error
}
