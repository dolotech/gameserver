package config

type WeaponStrengData struct {
	StrengLevel int `json:"strengLevel"` //强化等级
	Count       int `json:"count"`       //所需强化铁锭
	Gold        int `json:"gold"`        //所需金币
	PHY         int `json:"PHY"`         //提供物攻
	MAG         int `json:"MAG"`         //提供魔攻
	DF          int `json:"DF"`          //提供防御
	HP          int `json:"HP"`          //提供生命
	ResType     int `json:"resType"`     //适用职业
}

var weaponStrengPool **WeaponStrengDataPool

type WeaponStrengDataPool map[string]WeaponStrengData

func WeaponStreng() *WeaponStrengDataPool {
	return *weaponStrengPool
}

func (this *WeaponStrengDataPool) GetAllId() []string {
	ids := make([]string,0,len(*this))
	for id, _ := range *this {
		ids = append(ids, id)
	}
	return ids
}

func (this *WeaponStrengDataPool) Get(id string) WeaponStrengData {
	return (*this)[id]
}

func (this *WeaponStrengDataPool) GetByTypeAndLevel(resType int, strengthLevel int) WeaponStrengData {
	strength := WeaponStrengData{}
	for id, v := range *this {
		if v.ResType == resType && strengthLevel == v.StrengLevel {
			return (*this)[id]
		}
	}
	return strength
}
