package weapon

//武器强化
type Strengthen struct {
	PartnerId int `json:"partnerId"`
}


//武器强化返回
type Strengthened struct {
	NewEquipment NewEquipment `json:"newEquipment"`
	Power        int          `json:"power"`
	Code         int          `json:"code"`
	Gold         int          `json:"gold"`
}

type NewEquipment struct {
	Strength int `json:"strength"`
	Star     int `json:"star"`
	WeaponId int `json:"weaponId"`
}

//武器升星
type RisingStar struct {
	PartnerId int `json:"partnerId"`
}

//武器升星返回
type RisingStared struct {
	NewEquipment NewEquipment `json:"newEquipment"`
	Power        int          `json:"power"`
	Code         int          `json:"code"`
	Gold         int          `json:"gold"`
}

//武器进阶
type Advance struct {
	PartnerId int `json:"partnerId"`
}

//武器进阶返回
type Advanced struct {
	NewEquipment NewEquipment `json:"newEquipment"`
	Power        int          `json:"power"`
	Code         int          `json:"code"`
	Gold         int          `json:"gold"`
}