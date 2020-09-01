package config

type RoleStarData struct {
	StarLevel           int `json:"starLevel"`           //角色星级
	NeedTrain           int `json:"needTrain"`           //需求进阶等级
	NeedGold            int `json:"needGold"`            //需求金币	(填写当前值)
	NeedHonor           int `json:"needHonor"`           //需求荣誉
	SoldierStarPhysical int `json:"soldierStarPhysical"` //战士物攻
	SoldierStarDefense  int `json:"soldierStarDefense"`  //战士防御
	SoldierStarLife     int `json:"soldierStarLife"`     //战士生命	(填写累计值)
	KnightStarPhysical  int `json:"knightStarPhysical"`  //骑士物攻
	KnightStarDefense   int `json:"knightStarDefense"`   //骑士防御
	KnightStarLife      int `json:"knightStarLife"`      //骑士生命
	ShooterStarPhysical int `json:"shooterStarPhysical"` //射手物攻
	ShooterStarDefense  int `json:"shooterStarDefense"`  //射手防御
	ShooterStarLife     int `json:"shooterStarLife"`     //射手生命
	MageStarMagic       int `json:"mageStarMagic"`       //法师魔攻
	MageStarDefense     int `json:"mageStarDefense"`     //法师防御
	MageStarLife        int `json:"mageStarLife"`        //法师生命
	PriestStarMagic     int `json:"priestStarMagic"`     //牧师魔攻
	PriestStarDefense   int `json:"priestStarDefense"`   //牧师防御
	PriestStarLife      int `json:"priestStarLife"`      //牧师生命
}

var roleStarPool **RoleStarDataPool

type RoleStarDataPool map[string]RoleStarData

func RoleStar() *RoleStarDataPool {
	return *roleStarPool
}

func (this *RoleStarDataPool) GetAllId() []string {
	ids :=make( []string,0,len(*this))
	for id, _ := range *this{
		ids = append(ids, id)
	}
	return ids
}

func (this *RoleStarDataPool) Get(id string) RoleStarData {
	return (*this)[id]
}
