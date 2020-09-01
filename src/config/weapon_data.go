package config

type WeaponData struct {
	ResType      int `json:"resType"`      //适用职业 (相当人物id )
	Advanced     int `json:"advanced"`     //阶级
	PHY          int `json:"PHY"`          //物攻初始值
	MAG          int `json:"MAG"`          //魔攻初始值
	DF           int `json:"DF"`           //防御初始值
	HP           int `json:"HP"`           //生命初始值
	BuffId       int `json:"buffId"`       //具备buff
	BasisSkill   int `json:"basisSkill"`   //普通消对应技能id
	SpecialSkill int `json:"specialSkill"` //特殊消对应技能id
	StrengMax    int `json:"strengMax"`    //强化上限
}

var weaponPool **WeaponDataPool

type WeaponDataPool map[string]WeaponData

func Weapon() *WeaponDataPool {
	return *weaponPool
}

func (this *WeaponDataPool) GetAllId() []string {
	ids := make([]string,0,len(*this))
	for id, _ := range *this {
		ids = append(ids, id)
	}
	return ids
}

func (this *WeaponDataPool) Get(id string) WeaponData {
	return (*this)[id]
}
