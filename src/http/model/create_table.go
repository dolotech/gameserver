package model

import (
	"fmt"
	"gameserver/utils/db"
	"gameserver/utils/log"
)

const AUTO_INCREMENT  = 50000
type Auto_increment struct {
	Auto_increment int `gorm:"column:Auto_increment"`
}
func AutoMigrate() {

	if err := db.LDB().AutoMigrate(&User{}, &Server{}).Error; err != nil {
		log.Error(err)
	}
	//"SELECT Auto_increment FROM information_schema.tables WHERE Table_name='users'"
	var au Auto_increment

	err:=db.LDB().Raw(`SELECT Auto_increment FROM information_schema.tables WHERE Table_name='users'`).Limit(1).Find(&au).Error

	if err ==nil && au.Auto_increment < AUTO_INCREMENT{
		db.LDB().Exec(fmt.Sprintf("alter table users AUTO_INCREMENT=%d",AUTO_INCREMENT))
	}

	//db.LDB().Exec("alter table users AUTO_INCREMENT=50000")

	//alter table users AUTO_INCREMENT=50000;
}
