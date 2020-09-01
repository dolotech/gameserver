package model

import (
	"encoding/json"
	"gameserver/config"
	"gameserver/utils"
	"github.com/spf13/viper"
	"reflect"
	"strconv"
	"testing"
)

type AA struct{
	Name string
	Age int
}



func Test_Add_test(t *testing.T) {

	aa:= &AA{Name:"Michael",Age:99}

	aaData,_:=json.Marshal(aa)
	aaType:=reflect.TypeOf(aa)


	t.Error(aaType.Elem())
	aaNew:=reflect.New(aaType.Elem())

	err:=json.Unmarshal(aaData,aaNew.Interface())

	t.Error(aaNew,err)



}


func Test_Add_Robots(t *testing.T) {
	//users 表需要执行下面语句
	//alter table users AUTO_INCREMENT=50000;
	viper.SetConfigFile("cfg.toml")
	err := viper.ReadInConfig()
	t.Error(err)
	// 初始化json表
	config.Load(viper.GetString("general.configDir"))
	//---------------添加机器人到玩家库--------------------
	robots := config.Robots()
	for k, v := range *robots {
		pid, _ := strconv.Atoi(k)
		playerId := uint(pid)
		p := &Player{}
		p.PlayerId = playerId

		par1 := &Partner{PlayerId: playerId}
		par1.InitRobot(10001, v.Partner1)

		par2 := &Partner{PlayerId: playerId}
		par2.InitRobot(10002, v.Partner2)

		par3 := &Partner{PlayerId: playerId}
		par3.InitRobot(10003, v.Partner3)

		par4 := &Partner{PlayerId: playerId}
		par4.InitRobot(10004, v.Partner4)

		par5 := &Partner{PlayerId: playerId}
		par5.InitRobot(10005, v.Partner5)
		ps := Partners{par1, par2, par3, par4, par5}
		for _, v := range ps {
			v.Add()
		}
		p.Power = ps.CalcuPower(p.Level)

		utils.StructAtoB(p, v)
		t.Error(p.AddPlayer(), p)

		for _, it := range v.Items {
			item := &Item{
				PlayerId: playerId,
				ItemId:   it,
			}
			item.Add(99999)
		}
	}
	//----------------------------------------------------
}

func Test_Select_player(t *testing.T) {

	viper.SetConfigFile("cfg.toml")
	err := viper.ReadInConfig()
	t.Error(err)
	/*
		o:=Player{
			PlayerId:1,
			Nickname:"Michael001",
			ResType:1,
		}
		o.UpdateNicknameSex()*/

	oeSkill := OwnedEnergySkill{}
	t.Error(oeSkill.IsOwnedEnergySkill(1, 39002))
}

func Test_Player_Get_Nick(t *testing.T) {
	viper.SetConfigFile("cfg.toml")
	err := viper.ReadInConfig()
	t.Error(err)

	//ps := &Players{}
	//err = ps.GetByNick("浪漫的朱迪")
	//t.Error(err)
	//for _, v := range *ps {
	//	t.Error(v)
	//}
}
func Test_Player_Get_Multi(t *testing.T) {
	viper.SetConfigFile("cfg.toml")
	err := viper.ReadInConfig()
	t.Error(err)
	/*
		ps := &Players{}
		err = ps.GetIn([]uint{1, 2, 3})
		t.Error(err)
		for _, v := range *ps {
			t.Error(v)
		}
	*/

	ps := &Players{}
	err = ps.GetNotIn([]uint{1})
	t.Error(err, len(*ps))
	for _, v := range *ps {
		t.Error(v)
	}
}
