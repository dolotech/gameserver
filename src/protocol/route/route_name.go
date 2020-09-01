package route

const (
	//更新玩家信息
	OnUpdatePlayer = "onUpdatePlayer"
	//更新角色信息
	OnUpdatePartners = "onUpdatePartners"
	//更新好友模块信息
	OnUpdatePlayerFriends = "onUpdatePlayerFriends"
	//别人申请成为好友
	OnUpdateMakeFriend = "onUpdateMakeFriend"
	//同意对方加好友,好友上下线更新好友在线状态
	OnUpdateAgreeMakeFriend = "onUpdateAgreeMakeFriend"
	//我好友删除好友
	OnUpdateDelFriend = "onUpdateDelFriend"
	//更新背包
	OnUpdateBag = "onUpdateBag"
	//更新关卡信息
	OnUpdateThroughRecords = "onUpdateThroughRecords"
	//更新邮件信息
	OnUpdateEmails = "onUpdateEmails"
	//更新新邮件信息
	OnUpdateNewEmails = "onUpdateNewEmails"
	//更新任务
	OnUpdatePlayerTasks = "onUpdatePlayerTasks"
	//更新玩家武器列表
	OnUpdatePlayerEquipments = "onUpdatePlayerEquipments"
	//更新钻石
	OnUpdateBuyDiamond = "onUpdateBuyDiamond"
	//更新聊天信息
	OnChatReceiveMessage = "onChatReceiveMessage"
	//更新走马灯信息
	OnTipsReceiveMessage = "onTipsReceiveMessage"
	////////////////////////////////////////////////////////////////////-

	//更新好友挑战提示
	OnUpdatePvfFigthAlert = "onUpdatePvfFigthAlert"
	//更新匹配对手超时
	OnUpdateFindOverTime = "onUpdateFindOverTime"
	//更新备战状态
	OnUpdateReadyToFight = "onUpdateReadyToFight"
	//更新对手断线
	OnUpdateOpponentHadGone = "onUpdateOpponentHadGone"

	//更新战斗开始
	OnUpdateBattleStart = "onUpdateBattleStart"
	//更新战斗结束
	OnUpdateBattleEnd = "onUpdateBattleEnd"
	//更新战斗动作
	OnUpdateFightAction = "onUpdateFightAction"
	//更新角色状态
	OnUpdateSetFighterBuff = "onUpdateSetFighterBuff"

	//对方使用物品
	OnUpdateUseItem = "onUpdateUseItem"
	//////////////////////////////////////////////////////////////////////

	//登录游戏服
	EnterGame = "connector.entryHandler.entry"
	//设置昵称
	SetNickName = "logic.playerHandler.setNickname"
	//设置头像
	SetAvatar = "logic.playerHandler.setAvatar"

	//////////////////////////////////////////////////////////////-

	//换装
	ChangeFashion = "logic.playerHandler.changeFashion"
	//技能升级
	UpgradeSkill = "logic.playerHandler.upgradeSkill"

	//武器升星
	RisingStar = "logic.playerHandler.risingStar"
	//武器强化
	Strengthen = "logic.playerHandler.strengthen"
	//武器进阶
	Advance = "logic.playerHandler.advance"

	//////////////////////////////////////////////////////////////-

	//战斗外使用物品
	UseItem = "logic.playerHandler.useItem"
	//购买钻石
	BuyDiamond = "logic.playerHandler.buyDiamond"

	//////////////////////////////////////////////////////////////-

	//获取好友信息列表
	GetRelationInfo = "logic.playerHandler.getRelationInfo"
	//添加好友
	MakeFriend = "logic.playerHandler.makeFriend"
	//同意添加好友
	AgreeMakeFriend = "logic.playerHandler.agreeMakeFriend"
	//删除好友
	DelFriend = "logic.playerHandler.delFriend"
	//添加到黑名单列表
	AddBlacklist = "logic.playerHandler.addBlacklist"
	//从删除黑名单列表
	DelBlacklist = "logic.playerHandler.delBlacklist"
	//请求赠送体力
	GiveStamina = "logic.playerHandler.giveStamina"
	//回复赠送体力
	RecvStamina = "logic.playerHandler.recvStamina"
	//推荐好友
	RefreshRandomFriends = "logic.playerHandler.refreshRandomFriends"
	//查找好友
	FindFriends = "logic.playerHandler.findFriends"

	//////////////////////////////////////////////////////////////////////////

	//开始通关
	StartFightThrough = "logic.playerHandler.startFightThrough"
	//完成通关
	FinishFightThrough = "logic.playerHandler.finishFightThrough"

	//////////////////////////////////////////////////////////////////////////
	//获取PVB信息  ADD
	GetPvbInfo = "getPvbInfo"
	//重置PVB信息
	ResetChallenge = "logic.playerHandler.resetChallenge"
	//购买PVB次数
	BuyChallenge = "logic.playerHandler.buyChallenge"

	//开始PVB
	StartFightChallenge = "logic.playerHandler.startFightChallenge"
	//完成PVB
	FinishFightChallenge = "logic.playerHandler.finishFightChallenge"

	//////////////////////////////////////////////////////////////-

	//签到奖励
	RecvSignAward = "logic.playerHandler.recvSignAward"
	//等级礼包
	RecvLevelAward = "logic.playerHandler.recvLevelAward"
	//在线礼包
	RecvOnlineAward = "logic.playerHandler.recvOnlineAward"
	//开服礼包
	RecvOpenServiceAward = "logic.playerHandler.recvOpenServiceAward"

	//////////////////////////////////////////////////////////////-

	//进入PVP
	EnterArena = "logic.arenaHandler.enterArena"
	//PVP对手查找
	FindArenaOpponent = "logic.arenaHandler.findArenaOpponent"
	//取消PVP对手查找
	CancelFindArenaOpponent = "logic.arenaHandler.cancelFindArenaOpponent"
	//PVP资源都加载完毕
	FinishLoadBattleRes = "logic.arenaHandler.finishLoadBattleRes"
	//PVP奖励领取
	RecvArenaRankAward = "logic.arenaHandler.recvArenaRankAward"
	//PVP购买次数
	BuyArenaPKTimes = "logic.arenaHandler.buyArenaPKTimes"

	//PVP设置自己的BUFF添加跟删除
	SetFighterBuff = "fight.fightHandler.setFighterBuff"
	//PVP设置自己的属性更改 (发起一次攻击)
	SetFighterAttribute = "fight.fightHandler.setFighterAttribute"
	//PVP使用能量技
	ConjureSkillsOnBattle = "fight.fightHandler.conjureSkills"
	//战斗中使用物品
	UseItemOnBattle = "fight.fightHandler.useItem"

	//////////////////////////////////////////////////////////////////////////

	//获取PVA榜信息
	GetWeeklyBoardInfo = "logic.weeklyHandler.getWeeklyBoardInfo"
	//获取PVA场对手信息
	GetWeeklyOpponents = "logic.weeklyHandler.getWeeklyOpponents"
	//挑战PVA场对手
	ChallengeWeeklyOpponent = "logic.weeklyHandler.challengeWeeklyOpponent"
	//完成PVA场对手的挑战
	FinishChallengeWeeklyOpponent = "logic.weeklyHandler.finishChallengeWeeklyOpponent"
	//购买挑战PVA场次数
	BuyWeeklyChallengeTimes = "logic.weeklyHandler.buyWeeklyChallengeTimes"
	//获取挑战PVA的奖励
	RecvWeeklyAwards = "logic.weeklyHandler.recvWeeklyAwards"

	//////////////////////////////////////////////////////////////////////////////////////////////
	
	//获取章节奖励
	RecvChapterAward = "logic.playerHandler.recvChapterAward"
	//获取章节关卡数据
	GetThroughRecords = "logic.playerHandler.getThroughRecords"

	////////////////////////////////////////////////////////////////////////////////////////////

	//领取邮件附件
	RecvEmailItems = "logic.playerHandler.recvEmailItems"
	//购买物品
	ShopBuyItem = "logic.playerHandler.shopBuyItem"
	//出售物品
	SellItem = "logic.playerHandler.sellItem"
	//关卡购买物品
	DungeonBuyItem = "logic.playerHandler.dungeonBuyItem"

	////////////////////////////////////////////////////////////////////////////////////////////

	//接受任务
	RecvTask = "logic.playerHandler.recvTask"
	//完成任务
	RecvTaskAward = "logic.playerHandler.recvTaskAward"

	//获取星星奖励
	RecvCollectStarAward = "logic.playerHandler.recvCollectStarAward"

	////////////////////////////////////////////////////////-

	//角色训练
	Train = "logic.playerHandler.train"
	//角色升星
	UpgradeStar = "logic.playerHandler.upgradeStar"

	//补充体力
	BuyStamina = "logic.playerHandler.buyStamina"

	//发送信息
	ChatSendMessage = "logic.playerHandler.chatSendMessage"

	//发送走马灯
	TipsSendMessage = "logic.playerHandler.tipsSendMessage"

	////////////////////////////////////////////////////////////////////////////////////////////////-

	//购买技能卡
	BuyEnergySkill = "logic.playerHandler.buyEnergySkill"
	//激活技能卡
	ActiveEnergySkill = "logic.playerHandler.activeEnergySkill"

	////////////////////////////////////////////////////////////////////////////////////////////////-

	//异步引导记录
	SetGuideSteps = "logic.playerHandler.setGuideSteps"
	//主线引导记录
	Setnewstep = "logic.playerHandler.setNewStep"
)
