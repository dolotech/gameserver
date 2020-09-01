package model

import (
	"gameserver/utils/db"
	"gameserver/utils/log"
)

func AutoMigrate() {
	if err := db.Get().AutoMigrate(
		&Player{},
		&Friend{},
		&Order{},
		&Partner{},
		&Item{},
		&Email{},
		&Task{},
		&OwnedEnergySkill{},
		&DungeonRecord{},
		).Error; err != nil {
		log.Error(err)
	}
}

