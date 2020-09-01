package config

//签到奖励
type SignAwardData struct {
	DropId int `json:"dropId"`		//掉落id
}

var signAwardPool **SignAwardPool
type SignAwardPool map[string]SignAwardData

func SignAward() *SignAwardPool {
	return *signAwardPool
}

func (this *SignAwardPool) Get(id string) SignAwardData{
	return (*this)[id]
}
