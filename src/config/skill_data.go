package config

type SkillData struct {
	DamageType string      `json:"damage_type"`
	IsEnergy   int         `json:"isEnergy"`
	DamageAtk  []DamageAtk `json:"damage_atk"`
	AtkTarget  string      `json:"atk_target"`
	AtkNums    int         `json:"atk_nums"`
}

type DamageAtk struct {
	Value float32 `json:"value"`
	Per   float32 `json:"per"`
}

var skillDataPool **SkillDataPool

type SkillDataPool map[string]SkillData

func (this *SkillDataPool) Get(id string) (SkillData,bool) {
	d, ok := (*this)[id]
	return d, ok
}

func Skill() *SkillDataPool {
	return *skillDataPool
}
