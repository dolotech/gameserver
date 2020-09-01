package protocol

// 用于下发前端的错误消息
const (
	OK                            = 200  //成功
	FAIL                          = 500  //失败
	FA_USER_NOT_EXIST             = 1003 //用户名不存在
	FA_USERNAME_EXIST             = 1004 //用户名已经存在
	FA_INVALID_PARTNER            = 4003 //无效的角色id
	FA_NOT_ENOUGH_ITEMS           = 4009 //所需物品数量不足
	FA_OVER_FRIENDS_COUNT         = 4018 //已超出好友上限
	FA_NOT_ENOUGH_HONOR           = 4024 //荣誉值不够
	FA_INVALID_DUNGEON_ID         = 4032 //前置关卡未通关
	FA_CHALLENGE_MAX_TIMES        = 4035 //已到最大挑战次数
	FA_CHALLENGE_WRONG_DUNGEON_ID = 4037 //无尽挑战id错误
	FA_CHALLENGE_NO_CHANCE        = 4038 //挑战次数为空
	FA_INVALID_TASK_ID            = 4086 //无效任务id
	FA_TASK_NOT_DONE              = 4087 //任务未完成
	FA_HAD_RECV_DAY_AWARD         = 4041 //今天已经签到
	FA_LEVEL_LESS_THAN            = 4042 //等级不足
	FA_TIME_LESS_THAN             = 4043 //在线时间不足
	FA_ARENA_IS_MATCHING          = 4055 //正在匹配
	FA_INVALID_EMAIL_ID           = 4066 //emailId无效
	FA_INVALID_STAR_INDEX         = 4096 //集星坐标到顶或者错误
	FA_STAR_NOT_ENOUGH            = 4097 //武器星级不足 / 集星的星数不足
	FA_EQUIPMENT_MAX_STAR         = 4098 //武器星级已最大
	FA_EQUIPMENT_NOT_ENOUGH_LEVEL = 4099 //武器强化等级不足
	FA_NOT_ENOUGH_TRAIN_LEVEL     = 4103 //角色训练等级不足，无法升星
	FA_MAX_STRENGTH_LEVEL         = 4106 //武器已达顶级，需升阶或者已无法升级
	FA_MAX_TRAIN_LEVEL            = 4107 //角色已达训练最大等级
	FA_MAX_TRAIN_STAR             = 4108 //角色已达最大星级
	FA_MAX_STAMINA                = 4110 //体力已满
	FA_INVALID_TOOLS_INDEX        = 4111 //没这个推荐道具购买
	FA_NOT_ENOUGH_GOLD            = 4002 // 金币不足
	FA_NOT_ENOUGH_STAMINA         = 4015 // 体力不足
	FA_NOT_ENOUGH_DIAMOND         = 4022 //钻石不足

	TOO_FAST_CHAT           = 5001  // 发送聊天信息太快了
	LEGAL_NAME              = 5003  // 角色名字不合法
	WEEKLY_TIMES_NOT_ENOUGH = 60001 // 竞技场挑战次数不足

)
