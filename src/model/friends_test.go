package model

import (
	"github.com/spf13/viper"
	"testing"
)

func Test_Friends_Get_Multi(t *testing.T) {
	viper.SetConfigFile("cfg.toml")
	err := viper.ReadInConfig()
	t.Error(err)


}

func Test_item(t *testing.T) {
	viper.SetConfigFile("cfg.toml")
	err := viper.ReadInConfig()
	t.Error(err)

	i:=&Item{PlayerId:1,ItemId:100}
	err = i.Add(20)
	t.Error(err)

	i=&Item{PlayerId:1,ItemId:100}
	count,err:= i.GetCount()
	t.Error(err,count)

	i=&Item{PlayerId:1,ItemId:100}
	i.Reduce(21)

	it:=&Items{}
	it.Get(1)

	t.Error(it)

}
