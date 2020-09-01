package model

import "gameserver/utils/db"

//
type Server struct {
	ID         uint   `gorm:"primary_key;column:id;AUTO_INCREMENT;COMMENT:''" json:"id"`
	Addr       string `gorm:"type:varchar(200);column:addr" json:"addr"`
	Port       int    `gorm:"type:int(11);column:port" json:"port"`
	CreateTime int    `gorm:"type:int(11);column:createTime" json:"createTime"`
	Name       string `gorm:"type:varchar(30);column:name" json:"name"`
	Status     int    `gorm:"type:int(4);column:status" json:"status"`
}
type Servers []*Server

func (this *Servers) Get() error {
	return db.LDB().Model(this).Where("status > 0").Scan(this).Error
}
