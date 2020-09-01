package model

import "gameserver/utils/db"

// 邮件   gameserver
type Email struct {
	EmailId  uint `gorm:"primary_key;column:emailId;AUTO_INCREMENT;COMMENT:'邮件唯一标识'" json:"emailId" `
	PlayerId uint `gorm:"type:int(11);column:playerId;index:index_playerId" json:"playerId" `
	EndTime  int  `gorm:"type:int(11);column:endTime" json:"endTime" `
	TypeId   int  `gorm:"type:int(11);column:typeId" json:"typeId" `
	Status   int  `gorm:"type:int(4);column:status;comment:'1未发送2已发送3已领取'" json:"status" `
	Items    []byte  `gorm:"type:varchar(200);column:items;comment:'物品的json字符串，包括itemId和count'" json:"status" `
	//LeftTime int  `gorm:"type:int(11);column:leftTime" json:"leftTime" `
	//标题和文本应从配置获取
	//Head     int  `gorm:"varchar(100);column:head" json:"head" `
	//Text     int  `gorm:"type:text;column:text" json:"text" `
}
type Emails []*Email

func (this *Emails) Get(playerId uint) error {
	return db.Get().Model(this).Where("playerId = ?", playerId).Find(this).Error
}

func (this *Emails) GetNewMails(playerId uint, status int) error {
	return db.Get().Model(this).Where("playerId = ? and status <= ?", playerId, status).Find(this).Error
}

func (this *Emails) UpdateStatus(playerId uint, status int) error {
	return db.Get().Model(this).Where("playerId = ?", playerId).Update("status", status).Find(this).Error
}

func (this *Emails) GetByEmailIds(emailIds []uint, playerId uint, status int) error{

	return db.Get().Model(this).Where("emailId in (?) and playerId = ? and status = ?", emailIds, playerId, status).Find(this).Error
}

func (this *Email) AddMail() error {
	return db.Get().Create(this).Error
}

func (this *Email) UpdateOneStatus(emailIds []uint, status int) error {
	return db.Get().Model(this).Where("emailId in (?)", emailIds).Update("status", status).Find(this).Error
}
