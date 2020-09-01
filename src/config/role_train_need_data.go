package config


type RoleTrainNeedData struct {
	ItemId int `json:"itemId"`
	Num    int `json:"num"`
}

var RoleTrainNeedPool **RoleTrainNeedDataPool

type RoleTrainNeedDataPool map[string]RoleTrainNeedData

func RoleTrainNeed() *RoleTrainNeedDataPool {
	return *RoleTrainNeedPool
}

func (this *RoleTrainNeedDataPool) GetAllId() []string {
	ids := make([]string,0,len(*this))
	for id, _ := range *this {
		ids = append(ids, id)
	}
	return ids
}

func (this *RoleTrainNeedDataPool) Get(id string) RoleTrainNeedData {
	return (*this)[id]
}