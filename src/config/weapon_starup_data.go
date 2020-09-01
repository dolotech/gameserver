package config

type WeaponStarUpData struct {
	Star     int `json:"star"`     //星级
	WeaponId int `json:"weaponId"` //武器id
	PHY      int `json:"PHY"`      //物攻
	MAG      int `json:"MAG"`      //魔攻
	DF       int `json:"DF"`       //防御
	HP       int `json:"HP"`       //生命
	Honor    int `json:"honor"`    //升星花费荣誉
	Gold     int `json:"gold"`     //升星花费金币
}

func (this *WeaponStarUpPool) GetById(id int) {

}

type WeaponStarUpPool map[string]WeaponStarUpData

var weaponStarUp **WeaponStarUpPool

func WeaponStarUp() *WeaponStarUpPool {
	return *weaponStarUp
}

func (this *WeaponStarUpPool) GeyByIdAndStar(weaponId int, star int) WeaponStarUpData {
	start := WeaponStarUpData{}
	for id, v := range *this {
		if v.WeaponId == weaponId && v.Star == star {
			return (*this)[id]
		}
	}
	return start
}
