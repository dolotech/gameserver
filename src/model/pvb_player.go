package model

import (
	"errors"
	"gameserver/config"
	"gameserver/utils/db"
	"gameserver/utils/log"
	"time"
)

var (
	enoughChallengeTimes = errors.New("not enough challenge times")
	maxPvbNum = errors.New("challenge times is max")
)

func (this *Player) UpdateChallengeTimes(times int) error {
	timesAdd, LastChallengeTime := this.GetAddChallengeTimes()
	log.Error("-------times:", timesAdd, "time:", LastChallengeTime)
	if this.ChallengeTimes + timesAdd + times < 0{
		return enoughChallengeTimes
	}
	this.ChallengeTimes += timesAdd + times
	this.LastChallengeTime = LastChallengeTime
	return db.Get().Model(this).Updates(map[string]interface{}{
		"challengeTimes":    this.ChallengeTimes,
		"lastChallengeTime": this.LastChallengeTime,
	}).Error
}

func (this *Player) GetAddChallengeTimes() (timesAdd, lastChallengeTime int){
	now := int(time.Now().Unix())
	pvbMaxNum := config.Param().PvbMaxNum
	lastChallengeTime = now
	if this.LastChallengeTime != 0 {
		timesAdd = (now - this.LastChallengeTime) / config.Param().PvbRecover
		lastChallengeTime -= (now - this.LastChallengeTime)%config.Param().PvbRecover
	}
	if timesAdd+this.ChallengeTimes > pvbMaxNum {
		//只加到满次数为止
		timesAdd = config.Param().PvbMaxNum - this.ChallengeTimes
		lastChallengeTime = now
	}
	if this.ChallengeTimes >= pvbMaxNum {
		//原本就处于最大值，要更新一下领取时间
		timesAdd = 0
		lastChallengeTime = now
	}
	return
}

func (this *Player) UpdateChallengeId() error {
	return db.Get().Model(this).Update("maxChallengeId", this.MaxChallengeId).Error
}

func (this *Player) UpdateBuyChallenge(times int) error {
	timesAdd, LastChallengeTime := this.GetAddChallengeTimes()
	if this.ChallengeTimes + timesAdd > config.Param().PvbMaxNum {
		return maxPvbNum
	}
	return db.Get().Model(this).Updates(map[string]interface{}{
		"challengeTimes": this.ChallengeTimes + timesAdd + times,
		"diamond":        this.Diamond,
		"lastChallengeTime": LastChallengeTime,
	}).Error
}

