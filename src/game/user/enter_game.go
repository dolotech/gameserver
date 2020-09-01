package user

import (
	"gameserver/config"
	"gameserver/game/friend"
	"gameserver/game/mail"
	task2 "gameserver/game/task"
	"gameserver/model"
	"gameserver/online"
	proto "gameserver/protocol"
	"gameserver/protocol/push"
	"gameserver/protocol/route"
	"gameserver/utils"
	"gameserver/utils/filter"
	"gameserver/utils/log"
	msg "gameserver/utils/socket/message"
	"gameserver/utils/socket/server"
	"time"
)

// 游戏消息注册
func init() {
	msg.GetMsg().Reg(route.EnterGame, &proto.EntryGame{}, &proto.EntryGamed{}, enterGameCb)
	msg.GetMsg().Reg(route.SetNickName, &proto.SetNickname{}, &proto.SetNicknamed{}, setNicknameCb)
	msg.GetMsg().Reg(route.SetAvatar, &proto.SetAvatar{}, &proto.SetAvatared{}, setAvatarCb)
}

func createPlayer(playerId uint) (*model.Player, error) {
	player := &model.Player{
		PlayerId: playerId,
	}
	if err := player.Get(); err == nil {

	} else {
		playerCfg := config.GetPlayerCfg()
		param := config.Param()
		player.PlayerId = playerId
		player.Stamina = playerCfg.Stamina
		player.Diamond = playerCfg.Diamond
		player.Level = playerCfg.Level
		player.Exp = playerCfg.Exp
		player.Vip = playerCfg.Vip
		player.Honor = playerCfg.Honor
		player.CurThroughId = playerCfg.Dungeon.ThroughDungeonId
		player.MaxChallengeId = playerCfg.Dungeon.ChallengeDungeonId
		player.ChallengeTimes = param.PvbMaxNum
		player.DayWeeklyPKTimes = param.PvaMaxNum
		player.DayArenaPKTimes = param.PvpMaxNum
		player.PvpRank = 1
		//----------------------更新没有竞技塔排行的玩家-----------------------------------------
		player.WeeklyRank = (&model.Player{}).GetMaxRank()
		player.WeeklyRank += 1

		var partners model.Partners
		partners.InitPartner(player)
		player.Power = partners.CalcuPower(player.Level)
		if err := player.AddPlayer(); err != nil {
			log.Error(err)
			return player, err
		}
		// 赠送免费技能
		(&model.OwnedEnergySkill{PlayerId: playerId, SkillId: config.SkillPower().GetFreeSkill()}).AddEnergySkill()
		task := &task2.Task{PlayerId: playerId}
		task.InitByConfig()
	}
	return player, nil
}

// 连接建立完成客户端第一个业务协议，下发用户数据
func enterGameCb(sess server.Session, req *proto.EntryGame, resp *proto.EntryGamed) {
	resp.Code = proto.OK
	player, err := createPlayer(sess.UId())
	if err != nil {
		resp.Code = proto.FA_USER_NOT_EXIST
		return
	}
	partners, err := resp.Player.P2P(player)
	if err != nil {
		resp.Code = proto.FAIL
		return
	}
	now := time.Now().Unix()
	resp.PlayerId = player.PlayerId //ANJUN
	if player.ResType == 0 { //新手处理
		(&mail.Mail{}).RookieMail(player.PlayerId)
		resp.IsNew = 1
	}

	// 在玩家登录时重置需要每天重置的数据
	// 1,判断重置好友体力领取次数，要在更新登陆时间前处理
	// 2,签到相关
	//3,日常任务重置
	if !utils.SameDay(int64(player.LoginTime), now) {
		player.ResSetFrienRecvTimes()
		player.ResetPvpPKTimes()

		if !utils.SameMonth(int64(player.LoginTime), now) {
			player.UpdateMonthSign(0)
		}
		player.UpdateTodaySign(0)
		resp.IsNewDay = 1
		var myTasks = &task2.Task{PlayerId: player.PlayerId}
		mailTasks := myTasks.RestTask()
		if len(mailTasks) > 0 { //日常任务没领取转邮件
			(&mail.Mail{}).TaskMail(mailTasks, player.PlayerId)
		}
	}
	//----------------------更新上线时间等数据-----------------------------------------
	player.LoginUpdate()
	power := partners.CalcuPower(player.Level)
	resp.Player.Power = power
	player.Power = power
	//----------------------------保存一份玩家数据库到session--------------------------------
	userdata := sess.UserData().(*online.UserData)
	userdata.SetPartners(partners)
	userdata.SetPlayer(*player)
	// 后续的的接口不能用前端提交上来的playerID  一定要用session里验证过的playerID进行操作
	(&push.OnUpdatePlayer{}).Push(player,true)
	(&push.OnUpdatePlayerFriends{}).Push(sess.UId(),true)
	(&push.OnUpdateBag{}).Push(sess.UId(),true)
	(&push.OnUpdateThroughRecords{}).Push(player,true)
	//(&push.OnUpdateEmails{}).Push(player,true)
	(&push.OnUpdatePlayerTasks{}).Push(player,true)
	(&friend.FriendLoginAndLogout{}).OnLogin(sess)
	(&mail.Mail{}).SendNewMails(player.PlayerId,true) //统一发送新邮件
}

// 请求设置昵称
func setNicknameCb(sess server.Session, req *proto.SetNickname, resp *proto.SetNicknamed) {
	resp.Code = proto.OK
	player := &model.Player{PlayerId: sess.UId()}
	if err := player.GetById(); err != nil {
		resp.Code = proto.FA_USER_NOT_EXIST
		return
	}

	if !utils.LegalName(req.Nickname, 12) {
		resp.Code = proto.LEGAL_NAME
		return
	}

	if !filter.Valid("dic.txt", []rune(req.Nickname), []rune{}, '*') {
		resp.Code = proto.LEGAL_NAME
		return
	}

	if err := player.UpdateNicknameSex(req.Nickname, req.ResType); err != nil {
		log.Error("set nickname error, %s", err)
		resp.Code = proto.FAIL
		return
	}
	_, err := resp.Player.P2P(player)
	if err != nil {
		resp.Code = proto.FAIL
		return
	}
}

func setAvatarCb(sess server.Session, req *proto.SetAvatar, resp *proto.SetAvatared) {
	resp.Code = proto.OK
	player := &model.Player{PlayerId: sess.UId()}
	if !player.Exist() {
		resp.Code = proto.FA_USER_NOT_EXIST
		return
	}
	player.Avatar = req.Avatar
	resp.Avatar = req.Avatar
	if err := player.UpdateAvatar(req.Avatar); err != nil {
		log.Error("set UpdateAvatar error, %s", err)
		resp.Code = proto.FAIL
		return
	}
}
