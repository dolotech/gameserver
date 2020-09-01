package model

import (
	"gameserver/utils/db"
	"github.com/jinzhu/gorm"
)

func (this *Players) GetIn(playerIds []uint) error {
	err := db.Get().Table("players").Select("playerId,resType,vip,level, power,avatar,nickname").Find(this, "playerId in (?)", playerIds).Error
	return err
}

func (this *Players) GetNotIn(playerIds []uint) (err error) {
	if len(playerIds) > 0 {
		err = db.Get().Table("players").Select("playerId,resType,vip,level, power,avatar,nickname").Limit(10).Find(this, "nickname !='' and robot =0 and playerId not in (?) ", playerIds).Error
	} else {
		err = db.Get().Table("players").Select("playerId,resType,vip,level, power,avatar,nickname").Limit(10).Find(this, "nickname !='' and robot =0").Error
	}
	return
}

// 通过昵称查找一个好友
func (this *Players) GetByNick(self uint, nick string) error {
	return db.Get().Table("players").Select("playerId,resType,vip,level, power,avatar,nickname").Where("robot =0 and nickname=? and playerId not in(?) ", nick, self).Find(this).Error
}

// 每天领取好友赠送的体力的次数限定是好友的最大数量次
func (this *Player) ResSetFrienRecvTimes() error {
	return db.Get().Model(this).Update("staminaRecvTimes", 0).Error
}

func (this *Player) FrienRecvTimes() error {
	return db.Get().Model(this).UpdateColumn("staminaRecvTimes", gorm.Expr("staminaRecvTimes + ?", 1)).Error
}
