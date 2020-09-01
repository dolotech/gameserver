package model

import (
	"errors"
	"gameserver/utils/db"
	"gameserver/utils/log"
	"github.com/jinzhu/gorm"
)

// 背包物品   gameserver
type Item struct {
	PlayerId uint `gorm:"type:int(11);column:playerId;unique_index:index_playerId_itemId" json:"playerId" `
	ItemId   int  `gorm:"column:itemId;type:int unsigned not null;default:0;unique_index:index_playerId_itemId" json:"itemId" `
	Count    int  `gorm:"column:count;type:int unsigned not null;default:0;" json:"count" `
}

type Items []*Item

// 获取玩家的所有物品
func (this *Items) Get(playerId uint) error {
	return db.Get().Model(this).Where("count > 0 and playerId= ?", playerId).Find(this).Error
}

func (this *Item) UpdateItemCount(playerId uint, itemId int, count int) error {
	return db.Get().Model(this).Where("playerId=? and itemId=?", playerId, itemId).UpdateColumn("count", gorm.Expr("count + ?", count)).Error
}

// 传正值
func (this *Item) Reduce(count int) error {
	if count <= 0 {
		return errors.New("value error!")
	}
	c, err := this.GetCount()
	if err != nil {
		return err
	}
	if count > c {
		err := errors.New("item not enough")
		log.Error(err)
		return err
	}
	return db.Get().Model(this).Where("playerId=? and itemId =?", this.PlayerId, this.ItemId).UpdateColumn("count", gorm.Expr("count + ?", -count)).Error
}

// 传正值
func (this *Item) Add(count int) error {
	if count <= 0 {
		return errors.New("value error!")
	}
	if !this.Exist(this.PlayerId,this.ItemId){
		this.Count = count
		return db.Get().Create(this).Error
	}
	return db.Get().Model(this).Where("playerId=? and itemId =?", this.PlayerId, this.ItemId).UpdateColumn("count", gorm.Expr("count + ?", count)).Error
}

// 获取玩家物品数量
func (this *Item) GetCount() (int, error) {
	err := db.Get().Model(this).Where("playerId= ?", this.PlayerId).Where("itemId= ?", this.ItemId).Limit(1).Find(this).Error
	return this.Count, err
}

func (this *Item) Exist(playerId uint, id int) bool {
	return db.Get().Model(this).Select("count").Where("playerId= ? and itemId=?", playerId, id).Limit(1).Find(this).Error == nil
}
