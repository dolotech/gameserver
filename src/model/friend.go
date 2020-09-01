package model

import (
	"errors"
	"gameserver/utils"
	"gameserver/utils/db"
	"time"
)

const (
	Apply = 9 //1:请求添加
	START = 1 //可以赠送别人状态
	ING   = 2 //可以领取
	CLOSE = 3 //CLOSE
)

var (
	FriendAlready = errors.New("is friend already")
	AddedAlready  = errors.New("add friend already")
)

// 好友   gameserver
type Friend struct {
	PlayerId uint `gorm:"type:int(11);column:playerId;unique_index:index_playerId" json:"playerId" `
	FriendId uint `gorm:"type:int(11);column:friendId;unique_index:index_playerId" json:"friendId" `
	Status   int  `gorm:"type:int(4);column:status;default:0" json:"status" ` //
	GiftTime int  `gorm:"type:int(11);column:giftTime;default:0" json:"_" `   //  礼物赠送时间
}

type Friends []*Friend

// 获取玩家的好友列表
func (this *Friends) Get(playerID uint) error {
	err := db.Get().Model(this).Where("playerId = ?", playerID).Find(this).Error
	this.filterGiftTime()
	return err
}
func (v *Friend) GiftStatus(now int64) {
	if v.Status != ING && v.Status != Apply {
		giftt := int64(v.GiftTime)
		if !utils.SameDay(giftt, now) {
			v.Status = START
		} else {
			v.Status = CLOSE
		}
	}
}
func (this *Friends) filterGiftTime() {
	now := time.Now().Unix()
	for _, v := range *this {
		v.GiftStatus(now)
	}
}
func (this *Friends) GetFriend(playerID uint) error {
	err := db.Get().Model(this).Where("playerId = ?", playerID).Where("status <> ?", Apply).Find(this).Error
	this.filterGiftTime()
	return err
}

// 获取玩家的好友申请列表
func (this *Friends) GetApply(playerID uint) error {
	return db.Get().Model(this).Where("playerId = ?", playerID).Where("status=?", Apply).Find(this).Error
}

// 申请好友
func (this *Friend) ApplyFriend(playerID, friendID uint) error {
	// 判断是不是好友
	if err := this.Get(playerID, friendID); err == nil {
		return FriendAlready
	}
	// 已经申请等待对方通过
	if err := this.Get(friendID, playerID); err == nil {
		return AddedAlready
	}

	this.PlayerId = friendID
	this.FriendId = playerID
	this.Status = Apply
	return db.Get().Create(this).Error
}

func (this *Friend) UpateStatus(Status int) error {
	return db.Get().Model(this).Where("playerId = ? and friendId= ?", this.PlayerId, this.FriendId).Update("status", Status).Error
}

func (this *Friend) UpateGiftTime() error {
	return db.Get().Model(this).Where("playerId = ? and friendId= ?", this.PlayerId, this.FriendId).Update("giftTime", time.Now().Unix()).Error
}

// 好友数量
func (this *Friend) Count(playerID uint) int {
	var count int
	db.Get().Model(this).Where("playerId = ? ", playerID).Count(&count)
	return count
}

// 获取指定一个好友，包括已经在申请的
func (this *Friend) Get(playerID, friendID uint) error {
	err := db.Get().Model(this).Where("playerId = ?", playerID).Where("friendId=?", friendID).Find(this).Limit(1).Error
	now := time.Now().Unix()
	this.GiftStatus(now)
	return err
}

// 对方同意添加申请
func (this *Friend) AgreeFriend(playerID, friendID uint) error {
	this.PlayerId = playerID
	this.FriendId = friendID
	if err := this.UpateStatus(START); err != nil {
		return err
	}

	f := &Friend{}
	f.PlayerId = friendID
	f.FriendId = playerID
	f.Status = START
	return f.Add()
}

func (this *Friend) Add() error {
	return db.Get().Create(this).Error
}

func (this *Friend) Delete(playerID, friendID uint) error {
	if this.Get(playerID, friendID) == nil {
		if err := db.Get().Model(this).Where("playerId = ?", playerID).Where("friendId=?", friendID).Delete(&Friend{}).Error; err != nil {
			return err
		}
	}
	if this.Get(friendID, playerID) == nil {
		return db.Get().Model(this).Where("friendId = ?", playerID).Where("playerId=?", friendID).Delete(&Friend{}).Error
	}
	return nil
}
