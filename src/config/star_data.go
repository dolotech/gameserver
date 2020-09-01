package config


//集星数据
type StarData struct {
	Star   int `json:"star"`
	DropID int `json:"dropId"`
}

var StarPool **StarDataPool

type StarDataPool map[string]StarData

func Star() *StarDataPool {
	return *StarPool
}

func (this *StarDataPool) Get(id string) StarData {
	return (*this)[id]
}

func (this *StarDataPool) GetAllId() []string {
	ids := []string{}
	for id, _ := range *this {
		ids = append(ids, id)
	}
	return ids
}
