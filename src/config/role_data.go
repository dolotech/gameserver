package config

type RoleData struct {
	WeaponId       int `json:"weaponId"`       //初始化武器id
	Physical       int `json:"physical"`       //物理攻击
	Magic          int `json:"magic"`          //魔法攻击
	Defense        int `json:"defense"`        //防御力
	Life           int `json:"life"`           //生命值
	Speed          int `json:"speed"`          //攻击速度
	PhysicalGrowth int `json:"physicalGrowth"` //物理成长
	MagicGrowth    int `json:"magicGrowth"`    //魔法成长
	DefenseGrowth  int `json:"defenseGrowth"`  //Asis
	LifeGrowth     int `json:"lifeGrowth"`     //生命成长
	SpeedGrowth    int `json:"speedGrowth"`    //速度成长
}

var rolePool **RoleDataPool

type RoleDataPool map[string]RoleData

func Role() *RoleDataPool {
	return *rolePool
}

func (this *RoleDataPool) GetAllId() []string {
	ids := make([]string, 0, len(*this))
	for id, _ := range *this {
		ids = append(ids, id)
	}
	return ids
}

func (this *RoleDataPool) Get(id string) RoleData {
	return (*this)[id]
}
