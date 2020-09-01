package config

type buffevent struct {
	BECURE   int `json:"BE_CURE"`
	BEATTACK int `json:"BE_ATTACK"`
	LAPCOUNT int `json:"LAP_COUNT"`
}

type contact struct {
	STRANGER int `json:"STRANGER"`
	FRIEND   int `json:"FRIEND"`
	BLACK    int `json:"BLACK"`
}

type recstate struct {
	RECED  int `json:"REC_ED"`
	RECING int `json:"REC_ING"`
}

type tasktype struct {
	MAIN    int `json:"MAIN"`
	DAILY   int `json:"DAILY"`
	ACHIEVE int `json:"ACHIEVE"`
}

type arenastate struct {
	WAIT  int `json:"WAIT"`
	CLOSE int `json:"CLOSE"`
	OPEN  int `json:"OPEN"`
}

type attacktype struct {
	PHY   int `json:"PHY"`
	MAGIC int `json:"MAGIC"`
}

type gridtype struct {
	JELLYFISH int `json:"JELLYFISH"`
	FISH      int `json:"FISH"`
	GOLD      int `json:"GOLD"`
	SEAHORSE  int `json:"SEA_HORSE"`
	SEASTAR   int `json:"SEA_STAR"`
	CONCH     int `json:"CONCH"`
}

type taskjumptype struct {
	SPENDDIAMOND    int `json:"SPEND_DIAMOND"`
	GOLD            int `json:"GOLD"`
	STAMINA         int `json:"STAMINA"`
	DUNGEONALL      int `json:"DUNGEON_ALL"`
	PVPITMES        int `json:"PVP_ITMES"`
	PVAITMES        int `json:"PVA_ITMES"`
	FRIENDPVP       int `json:"FRIEND_PVP"`
	ROLESTAR        int `json:"ROLE_STAR"`
	PVAACH          int `json:"PVA_ACH"`
	PLAYERLEVEL     int `json:"PLAYER_LEVEL"`
	EQUIMENTSTRENTH int `json:"EQUIMENT_STRENTH"`
	FRIEND          int `json:"FRIEND"`
	ROLETRIAN       int `json:"ROLE_TRIAN"`
	ITEM            int `json:"ITEM"`
	PARTNER         int `json:"PARTNER"`
	PVPACH          int `json:"PVP_ACH"`
	ROLEADVCE       int `json:"ROLE_ADVCE"`
	PVB             int `json:"PVB"`
	POWERSKILL      int `json:"POWER_SKILL"`
	STAR            int `json:"STAR"`
	EQUIMENTADVCE   int `json:"EQUIMENT_ADVCE"`
	SIGN            int `json:"SIGN"`
	DUNGEON         int `json:"DUNGEON"`
	FRIENDGIVE      int `json:"FRIEND_GIVE"`
	SHOP            int `json:"SHOP"`
}

type buffeventobj struct {
	MYSELF int `json:"MYSELF"`
	ALLY   int `json:"ALLY"`
	BUFF   int `json:"BUFF"`
	ENEMY  int `json:"ENEMY"`
}

type damagetype struct {
	NONE  string `json:"NONE"`
	MAGIC string `json:"MAGIC"`
	PHY   string `json:"PHY"`
	ADDDE string `json:"ADD_DE"`
	ADDHP string `json:"ADD_HP"`
}

type buffeffectobj struct {
	MYSELF int `json:"MYSELF"`
	ALLY   int `json:"ALLY"`
	RANDOM int `json:"RANDOM"`
	ENEMY  int `json:"ENEMY"`
}

type giftstate struct {
	GIFTING int `json:"GIFT_ING"`
	GIFTED  int `json:"GIFT_ED"`
}

type getgiftstate struct {
	ING   int `json:"ING"`
	START int `json:"START"`
	CLOSE int `json:"CLOSE"`
}

type global struct {
	WEEKLYRANK int `json:"WEEKLY_RANK"`
	ARENA      int `json:"ARENA"`
}

type shoptype struct {
	ITEM    int `json:"ITEM"`
	HONOR   int `json:"HONOR"`
	GOLD    int `json:"GOLD"`
	PACKAGE int `json:"PACKAGE"`
}

type operator struct {
	DIANDIAN string `json:"DIAN_DIAN"`
}

type invitestate struct {
	CLOSE int `json:"CLOSE"`
	START int `json:"START"`
}

type pvprank struct {
	PLATINUM int `json:"PLATINUM"`
	BRONZE   int `json:"BRONZE"`
	DIAMOND  int `json:"DIAMOND"`
	SILVER   int `json:"SILVER"`
	GOLD     int `json:"GOLD"`
}

type fightmodel struct {
	PVB string `json:"PVB"`
	PVP string `json:"PVP"`
	PVA string `json:"PVA"`
	PVE string `json:"PVE"`
	PVF string `json:"PVF"`
}

type actioncode struct {
	DEL int `json:"DEL"`
	ADD int `json:"ADD"`
}

type result struct {
	WIN  int `json:"WIN"`
	LOSS int `json:"LOSS"`
	DRAW int `json:"DRAW"`
}

type itemtype struct {
	HONOR    int `json:"HONOR"`
	GOLD     int `json:"GOLD"`
	EXP      int `json:"EXP"`
	STAMINA  int `json:"STAMINA"`
	DIAMOND  int `json:"DIAMOND"`
	ITEM     int `json:"ITEM"`
	PVPSCORE int `json:"PVP_SCORE"`
}

type fighttarget struct {
	ENEMY  string `json:"ENEMY"`
	ME     string `json:"ME"`
	FRIEND string `json:"FRIEND"`
}
type attackratio struct {
	CRIT   int `json:"CRIT"`
	NORMAL int `json:"NORMAL"`
	FATAL  int `json:"FATAL"`
}

type taskstate struct {
	OPENED int `json:"OPENED"`
	CLOSED int `json:"CLOSED"`
	DOING  int `json:"DOING"`
	DONE   int `json:"DONE"`
}

type useitemids struct {
	HPDEITEMID int `json:"HPDE_ITEM_ID"`
}

type partnertype struct {
	WARRIOR  int `json:"WARRIOR"`
	KNIGHT   int `json:"KNIGHT"`
	MAGICIAN int `json:"MAGICIAN"`
	PRIEST   int `json:"PRIEST"`
	ARCHER   int `json:"ARCHER"`
}

type droptype struct {
	NORMAL    int `json:"NORMAL"`
	EMPTY     int `json:"EMPTY"`
	OTHERDROP int `json:"OTHERDROP"`
}

type starrate struct {
	STATUSTHREE float64 `json:"STATUSTHREE"`
	STATUSTWO   float64 `json:"STATUSTWO"`
	STEPTHREE   float64 `json:"STEPTHREE"`
	STEPTWO     float64 `json:"STEPTWO"`
}

type mailtype struct {
	ROOKIE int `json:"ROOKIE"`
	TASK   int `json:"TASK"`
}

type ConstantConfig struct {
	BUFFEVENT         buffevent     `json:"BUFF_EVENT"`
	CONTACT           contact       `json:"CONTACT"`
	RECSTATE          recstate      `json:"REC_STATE"`
	TASKTYPE          tasktype      `json:"TASK_TYPE"`
	ARENASTATE        arenastate    `json:"ARENA_STATE"`
	ATTACKTYPE        attacktype    `json:"ATTACK_TYPE"`
	GRIDTYPE          gridtype      `json:"GRID_TYPE"`
	TASKJUMPTYPE      taskjumptype  `json:"TASK_JUMP_TYPE"`
	MAXWEEKLYBUYTIMES int           `json:"MAX_WEEKLY_BUY_TIMES"`
	BUFFEVENTOBJ      buffeventobj  `json:"BUFF_EVENT_OBJ"`
	DAMAGETYPE        damagetype    `json:"DAMAGE_TYPE"`
	BUFFEFFECTOBJ     buffeffectobj `json:"BUFF_EFFECT_OBJ"`
	GIFTSTATE         giftstate     `json:"GIFT_STATE"`
	GETGIFTSTATE      getgiftstate  `json:"GET_GIFT_STATE"`
	GLOBAL            global        `json:"GLOBAL"`
	SHOPTYPE          shoptype      `json:"SHOP_TYPE"`
	OPERATOR          operator      `json:"OPERATOR"`
	MAXWEEKLYPKTIMES  int           `json:"MAX_WEEKLY_PK_TIMES"`
	INVITESTATE       invitestate   `json:"INVITE_STATE"`
	PVPRANK           pvprank       `json:"PVP_RANK"`
	FightModel        fightmodel    `json:"FightModel"`

	ACTIONCODE  actioncode  `json:"ACTION_CODE"`
	RESULT      result      `json:"RESULT"`
	ITEMTYPE    itemtype    `json:"ITEM_TYPE"`
	FIGHTTARGET fighttarget `json:"FIGHT_TARGET"`
	ATTACKRATIO attackratio `json:"ATTACK_RATIO"`
	TASKSTATE   taskstate   `json:"TASK_STATE"`
	USEITEMIDS  useitemids  `json:"USE_ITEM_IDS"`
	PARTNERTYPE partnertype `json:"PARTNER_TYPE"`
	DROPTYPE    droptype    `json:"DROP_TYPE"`
	STARRATE    starrate    `json:"STAR_RATE"`
	ROOKIEGUIDE int         `json:"ROOKIE_GUIDE"`
	MAILTYPE    mailtype    `json:"MAIL_TYPE""`
}

var constantConfig **ConstantConfig

func GetConstantCfg() *ConstantConfig {
	return *constantConfig
}
