package config

type LdtEveryAwardData struct {
	Award int `json:"award"`
}

var everyAwardPool **LdtEveryAwardPool

type LdtEveryAwardPool map[string]LdtEveryAwardData

func LdtEveryAward() *LdtEveryAwardPool {
	return *everyAwardPool
}

func (this *LdtEveryAwardPool) Get(id string) (LdtEveryAwardData,bool) {
	d,ok:= (*this)[id]
	return d, ok
}
