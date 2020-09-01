package skill

import (
	"gameserver/config"
	"gameserver/game/data"
	"gameserver/game/task"
	"gameserver/model"
	proto "gameserver/protocol"
	"gameserver/protocol/route"
	"gameserver/protocol/shop"
	"gameserver/utils/log"
	msg "gameserver/utils/socket/message"
	"gameserver/utils/socket/server"
	"strconv"
)

func init() {
	msg.GetMsg().Reg(route.BuyEnergySkill, &shop.BuyEnergySkill{}, &shop.BuyEnergySkilled{}, buyEnergySkillcb)
	msg.GetMsg().Reg(route.ActiveEnergySkill, &shop.ActiveEnergySkill{}, &shop.ActiveEnergySkilled{}, activeEnergySkillcb)
}

// 请求激活能量技
func activeEnergySkillcb(sess server.Session, req *shop.ActiveEnergySkill, resp *shop.ActiveEnergySkilled) {
	resp.Code = proto.OK
	p := &model.Player{PlayerId: sess.UId()}

	//find the config
	if _, ok := config.SkillPower().Get(strconv.Itoa(req.Skill)); !ok {
		resp.Code = proto.FAIL //FA_INVALID_SKILL
		log.Error("activeEnergySkillcb, invalid skill: %d", req.Skill)
		return
	}


	oeSkill := &model.OwnedEnergySkill{}
	// 判断是否已经拥有
	if !oeSkill.IsOwnedEnergySkill(p.PlayerId, req.Skill) {
		resp.Code = proto.FAIL //FA_SKILL_NOT_OWNED
		log.Error("activeEnergySkillcb,FA_SKILL_NOT_OWNED:%d", req.Skill)
		return
	}
	// 激活技能
	if err := oeSkill.SetUsedEnergySkill(p.PlayerId,req.Skill,1); err != nil {
		log.Error("激活技能失败:", err)
		resp.Code = proto.FAIL
		return
	}
	resp.Skill = req.Skill
}

// 请求购买技能卡
func buyEnergySkillcb(sess server.Session, req *shop.BuyEnergySkill, resp *shop.BuyEnergySkilled) {
	resp.Code = proto.OK

	p := &model.Player{PlayerId: sess.UId()}
	if p.GetDiamondGoldHonor() != nil {
		resp.Code = proto.FAIL
		return
	}

	//find the config
	var configSkill config.SkillPowerData
	var ok bool
	if configSkill, ok = config.SkillPower().Get(strconv.Itoa(req.Skill)); !ok {
		resp.Code = proto.FAIL //FA_INVALID_SKILL
		log.Error("buyEnergySkill, invalid skill: %d", req.Skill)
		return
	}
	// 判断是否已经拥有
	ownedSkill := model.OwnedEnergySkill{}
	if ownedSkill.IsOwnedEnergySkill(p.PlayerId, req.Skill) {
		resp.Code = proto.FAIL //FA_SKILL_DUPLICATED
		log.Error("buyEnergySkill,FA_SKILL_DUPLICATED:%d", req.Skill)
		return
	}

	myTasks := &task.Task{PlayerId: p.PlayerId}
	// 处理消耗品
	if configSkill.ItemId == data.GOLD {
		if p.Gold < configSkill.Count {
			resp.Code = proto.FA_NOT_ENOUGH_GOLD
			return
		}
		if err := p.UpdateGold(-configSkill.Count); err != nil {
			resp.Code = proto.FAIL
			return
		}
		myTasks.UpdatePlayerGold(configSkill.Count)
	} else if configSkill.ItemId == data.DIAMOND {
		if p.Diamond < configSkill.Count {
			resp.Code = proto.FA_NOT_ENOUGH_DIAMOND
			return
		}
		if err := p.UpdateDiamond(-configSkill.Count); err != nil {
			resp.Code = proto.FAIL
			return
		}
		myTasks.UpdateSpendDiamondDaily()
	}
	myTasks.UpdateSkillPower()

	// 增加技能
	oeSkill := &model.OwnedEnergySkill{
		PlayerId: p.PlayerId,
		SkillId:  req.Skill,
	}
	if err := oeSkill.AddEnergySkill(); err != nil {
		log.Error("增加技能失败:", err)
		resp.Code = proto.FAIL
		return
	}
	resp.Diamond = p.Diamond
	resp.Gold = p.Gold

	myTasks.PushTasks()
}