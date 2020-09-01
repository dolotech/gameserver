package weapon

import (
	"gameserver/config"
	"gameserver/game/task"
	"gameserver/model"
	proto "gameserver/protocol"
	"gameserver/protocol/push"
	"gameserver/protocol/route"
	"gameserver/protocol/weapon"
	"gameserver/utils/log"
	msg "gameserver/utils/socket/message"
	"gameserver/utils/socket/server"
	"reflect"
	"strconv"
)

// 消息注册
func init() {
	msg.GetMsg().Reg(route.Strengthen, &weapon.Strengthen{}, &weapon.Strengthened{}, strengthen)
	msg.GetMsg().Reg(route.RisingStar, &weapon.RisingStar{}, &weapon.RisingStared{}, risingStar)
	msg.GetMsg().Reg(route.Advance, &weapon.Advance{}, &weapon.Advanced{}, advance)
}

//武器强化
func strengthen(sess server.Session, req *weapon.Strengthen, resp *weapon.Strengthened) {
	resp.Code = proto.OK
	player := &model.Player{PlayerId: sess.UId()}
	if err := player.Get(); err != nil {
		resp.Code = proto.FA_USER_NOT_EXIST
		return
	}
	partners := model.Partners{}
	index := -1
	partners.Get(sess.UId())
	for i, p := range partners {
		if p.PartnerId == req.PartnerId{
			index = i
		}
	}
	if index == -1 {
		resp.Code = proto.FA_INVALID_PARTNER
		return
	}
	nextStrength := getNextStrength(partners[index])
	if reflect.DeepEqual(nextStrength, &config.WeaponStrengData{}) {
		resp.Code = proto.FA_MAX_STRENGTH_LEVEL
		return
	}
	if player.Gold < nextStrength.Gold {
		resp.Code = proto.FA_NOT_ENOUGH_GOLD
		return
	}

	//数据库操作
	item := model.Item{PlayerId: sess.UId(), ItemId: config.Param().IngotsId}
	if err := item.Reduce(nextStrength.Count);err != nil{
		resp.Code = proto.FA_NOT_ENOUGH_ITEMS
		return
	}
	if err := partners[index].UpdateStrength(nextStrength.StrengLevel);err != nil {
		log.Error(err)
	}
	partners[index].Strength = nextStrength.StrengLevel
	player.Gold -= nextStrength.Gold
	player.Power = partners.CalcuPower(player.Level)
	if err := player.UpdatePowerAndGold();err != nil{
		log.Error(err)
	}

	//任务处理
	myTasks := &task.Task{PlayerId: player.PlayerId}
	myTasks.UpdatePlayerGold(nextStrength.Gold)
	myTasks.UpdateWeaponStrength()
	myTasks.PushTasks()

	(&push.OnUpdateBag{PlayerId:player.PlayerId}).Push(player.PlayerId)
	(&push.OnUpdatePartners{}).Push(player, &partners)
	//返回
	resp.Gold = player.Gold
	resp.Power = player.Power
	resp.NewEquipment.WeaponId = partners[index].WeaponId
	resp.NewEquipment.Star = partners[index].WStar
	resp.NewEquipment.Strength = partners[index].Strength
}


//武器升星
func risingStar(sess server.Session, req *weapon.RisingStar, resp *weapon.RisingStared){
	resp.Code = proto.OK
	player := &model.Player{PlayerId: sess.UId()}
	if err := player.Get(); err != nil {
		resp.Code = proto.FA_USER_NOT_EXIST
		return
	}
	partners := model.Partners{}
	index := -1
	partners.Get(sess.UId())
	for i, p := range partners {
		if p.PartnerId == req.PartnerId{
			index = i
		}
	}
	if index == -1 {
		resp.Code = proto.FA_INVALID_PARTNER
		return
	}
	starCfg := getStarCfg(partners[index])
	if reflect.DeepEqual(starCfg, &config.WeaponStarUpData{}) {
		resp.Code = proto.FA_EQUIPMENT_MAX_STAR
		return
	}
	if player.Gold < starCfg.Gold {
		resp.Code = proto.FA_NOT_ENOUGH_GOLD
		return
	}
	if player.Honor < starCfg.Honor{
		resp.Code = proto.FA_NOT_ENOUGH_HONOR
		return
	}

	//数据库操作
	player.Gold -= starCfg.Gold
	player.Honor -= starCfg.Honor
	partners[index].WStar = starCfg.Star + 1
	player.Power = partners.CalcuPower(player.Level)
	if err := player.UpdateRisingStar();err != nil{
		log.Error(err)
	}
	if err := partners[index].UpdateWStar(partners[index].WStar);err != nil {
		log.Error(err)
	}

	//任务处理
	myTasks := &task.Task{PlayerId: player.PlayerId}
	myTasks.UpdatePlayerGold(starCfg.Gold)
	//消息推送
	myTasks.PushTasks()
	(&push.OnUpdatePartners{}).Push(player, &partners)
	//返回
	resp.Gold = player.Gold
	resp.Power = player.Power
	resp.NewEquipment.WeaponId = partners[index].WeaponId
	resp.NewEquipment.Star = partners[index].WStar
	resp.NewEquipment.Strength = partners[index].Strength
}

//武器进阶
func advance(sess server.Session, req *weapon.Advance, resp *weapon.Advanced){
	resp.Code = proto.OK
	player := &model.Player{PlayerId: sess.UId()}
	if err := player.Get(); err != nil {
		resp.Code = proto.FA_USER_NOT_EXIST
		return
	}
	partners := model.Partners{}
	index := -1
	partners.Get(sess.UId())
	for i, p := range partners {
		if p.PartnerId == req.PartnerId{
			index = i
		}
	}
	if index == -1 {
		resp.Code = proto.FA_INVALID_PARTNER
		return
	}
	advance := getAdvance(partners[index])
	if reflect.DeepEqual(advance, &config.WeaponAdvancedData{}) {
		//理论上应有一个武器进阶满级的错误，后面可视情况添加
		resp.Code = proto.FA_MAX_STRENGTH_LEVEL
		return
	}
	if player.Gold < advance.Golds {
		resp.Code = proto.FA_NOT_ENOUGH_GOLD
		return
	}
	if partners[index].WStar < advance.Star{
		resp.Code = proto.FA_STAR_NOT_ENOUGH
		return
	}
	if partners[index].Strength < advance.Level {
		resp.Code = proto.FA_EQUIPMENT_NOT_ENOUGH_LEVEL
		return
	}

	//数据库操作
	item := model.Item{PlayerId: sess.UId(), ItemId: advance.ItemId}
	if err := item.Reduce(advance.Num);err != nil{
		resp.Code = proto.FA_NOT_ENOUGH_ITEMS
		return
	}
	if err := partners[index].UpdateWeapon(advance.NextId);err != nil {
		log.Error(err)
	}
	partners[index].WeaponId = advance.NextId
	partners[index].WStar = 1 //新武器，star为1
	player.Gold -= advance.Golds
	player.Power = partners.CalcuPower(player.Level)
	if err := player.UpdatePowerAndGold();err != nil{
		log.Error(err)
	}

	//任务处理
	myTasks := &task.Task{PlayerId: player.PlayerId}
	myTasks.UpdatePlayerGold(advance.Golds)
	//这里用weaponId对100取余来判断武器的阶数，因此同一类武器weaponId必须递增
	myTasks.UpdateWeaponAdvce(advance.NextId % 100)
	//消息推送
	myTasks.PushTasks()
	(&push.OnUpdateBag{PlayerId:player.PlayerId}).Push(player.PlayerId)
	(&push.OnUpdatePartners{}).Push(player, &partners)

	//返回
	resp.Gold = player.Gold
	resp.Power = player.Power
	resp.NewEquipment.WeaponId = partners[index].WeaponId
	resp.NewEquipment.Star = partners[index].WStar
	resp.NewEquipment.Strength = partners[index].Strength
}

func getNextStrength(partner *model.Partner) *config.WeaponStrengData {
	nextStrength := config.WeaponStrengData{}
	weaponCfg := config.Weapon().Get(strconv.Itoa(partner.WeaponId))
	if partner.Strength >= weaponCfg.StrengMax {
		return &nextStrength
	}
	//若Strength + 1高于最高level，会返回一个空的struct
	nextStrength = config.WeaponStreng().GetByTypeAndLevel(weaponCfg.ResType, partner.Strength+1)
	return &nextStrength
}

func getStarCfg(partner *model.Partner) *config.WeaponStarUpData{
	starCfg := config.WeaponStarUpData{}
	nextStar := config.WeaponStarUp().GeyByIdAndStar(partner.WeaponId, partner.WStar + 1)
	if reflect.DeepEqual(nextStar, config.WeaponStarUpData{}) {
		return &starCfg
	}
	starCfg = config.WeaponStarUp().GeyByIdAndStar(partner.WeaponId, partner.WStar)
	return &starCfg
}

func getAdvance(partner *model.Partner) *config.WeaponAdvancedData{
	advance := config.WeaponAdvance().Get(strconv.Itoa(partner.WeaponId))
	return &advance
}