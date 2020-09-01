package config

import (
	"gameserver/utils/log"
	"strconv"
)

type SkillPowerData struct {
	Power  int `json:"power"`  //能量值
	ItemId int `json:"itemid"` //开启条件
	Count  int `json:"count"`  //花费数量
}

var skillPool **SkillPowerDataPool

type SkillPowerDataPool map[string]SkillPowerData

func SkillPower() *SkillPowerDataPool {
	return *skillPool
}

func (this *SkillPowerDataPool) Get(id string) (SkillPowerData, bool) {
	d, ok := (*this)[id]
	return d, ok
}

func (this *SkillPowerDataPool) GetAllId() []string {
	ids := make([]string, 0, len(*this))
	for id, _ := range *this {
		ids = append(ids, id)
	}
	return ids
}

func (this *SkillPowerDataPool) GetFreeSkill() int {
	for id, v := range *this {
		if v.Count == 0 {
			skillId, _ := strconv.Atoi(id)
			return skillId
		}
	}
	log.Error("no free skill in config of skill power data")
	return -1
}
