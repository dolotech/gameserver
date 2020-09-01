package config

//物品掉落数据
type DropData struct {
	DropId      int `json:"dropid"`      //掉落id
	ItemId      int `json:"itemid"`      //道具id		(正常的道具id。	如果填0.表示没有任何掉落。	还可以填其他掉落id )
	Count       int `json:"count"`       //道具数量
	Droptype    int `json:"droptype"`    //掉落类型		( 1.普通的道具，包括货币。2.不掉落。也就是说，如果填了2，表示对应的参数虽然占据掉落几率，但是它不会掉落出来。3.表示掉落id。如果次数填了3，表示对应的参数是另外一个掉落id，也就是掉落里面有掉落。)
	Probability int `json:"probability"` //掉落几率		(掉落几率按权重计算。如果填写-1，表示该道具必然掉出来。)
}

var dropPool **DropDataPool

type DropDataPool map[string]DropData

func Drop() *DropDataPool {
	return *dropPool
}

func (this *DropDataPool) Get(id string) DropData {
	return (*this)[id]
}

func (this *DropDataPool) GetByDropId(dropId int) []DropData {
	drops := make([]DropData,0,len(*this))
	for _, v := range *this {
		if v.DropId == dropId {
			drops = append(drops, v)
		}
	}
	return drops
}

func (this *DropDataPool) GetAllId() []string {
	ids :=make([]string,0,len(*this))
	for id, _ := range *this {
		ids = append(ids, id)
	}
	return ids
}
