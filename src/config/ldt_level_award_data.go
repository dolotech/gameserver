package config

//LdtLevelAwardData
type LdtLevelAwardData struct {
	Level int `json:"level"`
	Award int `json:"award"`
}

var ldtLevelAwardDataPool **LdtLevelAwardDataPool

type LdtLevelAwardDataPool map[string]LdtLevelAwardData

func LdtLevelAward() *LdtLevelAwardDataPool {
	return *ldtLevelAwardDataPool
}

func (this *LdtLevelAwardDataPool) Get(rank int) int {
	var key string
	for k, v := range *this {
		if rank == v.Level {
			return v.Award
		} else if v.Level > rank {
			if key == "" {
				key = k
			} else {
				if v.Level < (*this)[key].Level {
					key = k
				}
			}
		}
	}
	if key == "" {
		return 0
	}
	return (*this)[key].Award

}
