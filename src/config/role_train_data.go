package config

type RoleTrainData struct {
	TrainExp        int `json:"trainExp"`        //训练度
	SoldierPhysical int `json:"soldierPhysical"` //战士物攻
	SoldierDefense  int `json:"soldierDefense"`  //战士防御
	SoldierLife     int `json:"soldierLife"`     //战士生命
	KnightPhysical  int `json:"knightPhysical"`  //骑士物攻
	KnightDefense   int `json:"knightDefense"`   //骑士防御
	KnightLife      int `json:"knightLife"`      //骑士生命
	ShooterPhysical int `json:"shooterPhysical"` //射手物攻
	ShooterDefense  int `json:"shooterDefense"`  //射手防御
	ShooterLife     int `json:"shooterLife"`     //射手生命
	MageMagic       int `json:"mageMagic"`       //法师魔攻
	MageDefense     int `json:"mageDefense"`     //法师防御
	MageLife        int `json:"mageLife"`        //法师生命
	PriestMagic     int `json:"priestMagic"`     //牧师魔攻
	PriestDefense   int `json:"priestDefense"`   //牧师防御
	PriestLife      int `json:"priestLife"`      //牧师生命
}

var roleTrainPool **RoleTrainDataPool

type RoleTrainDataPool map[string]RoleTrainData

func RoleTrain() *RoleTrainDataPool {
	return *roleTrainPool
}

func (this *RoleTrainDataPool) GetAllId() []string {
	ids := make([]string,0,len(*this))
	for id, _ := range *this {
		ids = append(ids, id)
	}
	return ids
}

func (this *RoleTrainDataPool) Get(id string) RoleTrainData {
	return (*this)[id]
}
