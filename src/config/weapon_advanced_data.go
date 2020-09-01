package config

type WeaponAdvancedData struct {
	NextId int `json:"nextId"`
	Star   int `json:"star"`
	Level  int `json:"level"`
	Golds  int `json:"golds"`
	ItemId int `json:"itemId"`
	Num    int `json:"num"`
}


type WeaponAdvancedPool map[string]WeaponAdvancedData

var WeaponAdvanced **WeaponAdvancedPool

func WeaponAdvance() *WeaponAdvancedPool {
	return *WeaponAdvanced
}

func (this *WeaponAdvancedPool) Get(id string) WeaponAdvancedData {
	return (*this)[id]
}