package pvp

// 好友对战
import (
	"gameserver/config"
	"gameserver/game/pvp/room"
	"gameserver/game/robot"
	task2 "gameserver/game/task"
	"gameserver/model"
	"gameserver/online"
	proto "gameserver/protocol"
	"gameserver/protocol/push"
	"gameserver/protocol/pvp"
	"gameserver/protocol/route"
	"gameserver/utils/log"
	msg "gameserver/utils/socket/message"
	"gameserver/utils/socket/server"
	"strconv"
)

// 游戏消息注册
func init() {
	msg.GetMsg().Reg(route.FindArenaOpponent, &pvp.FindArenaOpponent{}, &pvp.FindArenaOpponented{}, findArenaOpponentCb)
	msg.GetMsg().Reg(route.FinishLoadBattleRes, &pvp.FinishLoadBattleRes{}, &pvp.FinishLoadBattleResed{}, finishLoadBattleResCb)
	msg.GetMsg().Reg(route.SetFighterAttribute, &pvp.SetFighterAttribute{}, &pvp.SetFighterAttributed{}, setFighterAttributeCb)
	msg.GetMsg().Reg(route.CancelFindArenaOpponent, &pvp.CancelFindArenaOpponent{}, &pvp.CancelFindArenaOpponented{}, cancelFindArenaOpponentCb) //更新对手断线
	msg.GetMsg().Reg(route.UseItemOnBattle, &pvp.UseItemOnBattle{}, &pvp.UseItemOnBattled{}, useItemOnBattle)
	msg.GetMsg().Reg(route.ConjureSkillsOnBattle, &pvp.ConjureSkills{}, &pvp.ConjureSkillsed{}, conjureSkills)
}

// PVP使用能量技
func conjureSkills(sess server.Session, req *pvp.ConjureSkills, resp *pvp.ConjureSkillsed) {
	resp.Code = proto.OK
	userData := sess.UserData().(*online.UserData)
	if userData.GetRoom() == 0 {
		resp.Code = proto.FAIL
		return
	}
	// 判断skillId是否为能量技
	if _, ok := config.SkillPower().Get(strconv.Itoa(req.SkillId)); !ok {
		resp.Code = proto.FAIL
		log.Error("SkillId not found")
		return
	}
	// 获取发动攻击的角色，我们把能量技定为1号位发动
	var fightIndex int
	r := robot.Get().Get(uint32(req.BattleId))
	if r == nil {
		log.Error("room == nil")
		resp.Code = proto.FAIL
		return
	}

	if sess.UId() == r.GetPlayer(0).UId() {
		fightIndex = 1
	} else {
		fightIndex = 6
	}
	room.Fire(r.Get(sess.UId()), r, fightIndex, req.SkillId, 0)
}

// 战斗中使用道具
func useItemOnBattle(sess server.Session, req *pvp.UseItemOnBattle, resp *pvp.UseItemOnBattled) {
	resp.Code = proto.OK
	resp.ItemId = req.ItemId
	resp.Key = req.Key

	userData := sess.UserData().(*online.UserData)
	if userData.GetRoom() == 0 {
		resp.Code = proto.FAIL
		return
	}

	use := &push.OnUpdateUseItem{ItemId: req.ItemId}
	room := robot.Get().Get(userData.GetRoom())
	if room == nil {
		log.Error("room == nil")
		resp.Code = proto.FAIL
		return
	}

	room.PushOpp(sess.UId(), use)
}

// 不同意pk,对手断线
func cancelFindArenaOpponentCb(sess server.Session, req *pvp.CancelFindArenaOpponent, resp *pvp.CancelFindArenaOpponented) {
	resp.Code = proto.OK
	robot.Get().Cancel(sess.UId())

	// todo 通知发起方取消匹配
}

// 请求PVP设置自己的属性更改 (发起一次攻击)
func setFighterAttributeCb(sess server.Session, req *pvp.SetFighterAttribute, resp *pvp.SetFighterAttributed) {
	resp.Code = proto.OK
	userData := sess.UserData().(*online.UserData)
	if userData.GetRoom() == 0 {
		resp.Code = proto.FAIL
		log.Error("room is empty")
		return
	}
	r := robot.Get().Get(userData.GetRoom())
	if r == nil {
		log.Error("room == nil")
		resp.Code = proto.FAIL
		return
	}
	room.Fire(r.Get(sess.UId()), r, req.FighterIndex, req.SkillId, req.AttackRatio)
}

// 请求PVP资源都加载完毕
func finishLoadBattleResCb(sess server.Session, req *pvp.FinishLoadBattleRes, resp *pvp.FinishLoadBattleResed) {
	resp.Code = proto.OK
	log.Info("请求PVP资源都加载完毕", sess.UId())
	userData := sess.UserData().(*online.UserData)

	r := robot.Get().Get(userData.GetRoom())
	if r == nil {
		log.Error("room == nil")
		resp.Code = proto.FAIL
		return
	}
	log.Info("玩家准备", sess.UId())
	r.Ready()
}

// 请求对手查找或者同意PK
func findArenaOpponentCb(sess server.Session, req *pvp.FindArenaOpponent, resp *pvp.FindArenaOpponented) {
	resp.Code = proto.OK
	if req.OpponentId == sess.UId() {
		resp.Code = proto.FAIL
		return
	}
	userData := sess.UserData().(*online.UserData)
	if userData.GetRoom() > 0 { //已经在战斗中
		resp.Code = proto.FAIL
		return
	}

	// 更新到最新的用户数据库
	player := &model.Player{PlayerId: sess.UId()}
	if err := player.Get(); err != nil {
		log.Error(err)
		resp.Code = proto.FAIL
		return
	}
	// 更新到最新的用户数据库
	partners := &model.Partners{}
	if err := partners.Get(sess.UId()); err != nil {
		log.Error(err)
		resp.Code = proto.FAIL
		return
	}
	userData.SetPlayer(*player)
	userData.SetPartners(*partners)

	if req.Fightmodel == "pvp" {
		if player.DayArenaPKTimes < 1 {
			resp.Code = proto.FAIL
			return
		}

	}

	log.Info("请求对手查找", req)
	robot.Get().Match(robot.NewMatch(sess, req.Fightmodel, req.OpponentId, player.Power))

	//任务 先放这里，后面应该在游戏开始后执行任务(目前无法在开始后区分pvf和pvp)
	task := &task2.Task{PlayerId:player.PlayerId}
	if req.Fightmodel == "pvf"{
		task.UpdateFriendPVP()
	}else if req.Fightmodel == "pvp"{
		task.UpdatePvpDaily()
	}
	task.PushTasks()
}
