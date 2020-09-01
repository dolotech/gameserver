package config

import (
	"gameserver/utils"
)

func  Reload(pathDir string) {
	utils.Json().LoadFile(pathDir)
	utils.Json().Fill()
}
func Load(pathDir string) {
	utils.Json().LoadFile(pathDir)
	utils.Json().RegistJson("ArenaConfig", &arenaConfig)
	utils.Json().RegistJson("ShopData", &shopPool)
	utils.Json().RegistJson("DiamondData", &diamondDataPool)
	utils.Json().RegistJson("RobotData", &robotPool)
	utils.Json().RegistJson("ParamData", &paramData)
	utils.Json().RegistJson("TaskData", &taskPool)
	utils.Json().RegistJson("RoleData", &rolePool)
	utils.Json().RegistJson("PlayerConfig", &playerConfig)
	utils.Json().RegistJson("ChapterData", &chapterPool)
	utils.Json().RegistJson("SkillPowerData", &skillPool)
	utils.Json().RegistJson("WeaponData", &weaponPool)
	utils.Json().RegistJson("WeaponAdvancedData", &WeaponAdvanced)
	utils.Json().RegistJson("TollgateData", &tollgatePool)
	utils.Json().RegistJson("RoleTrainData", &roleTrainPool)
	utils.Json().RegistJson("RoleTrainNeedData", &RoleTrainNeedPool)
	utils.Json().RegistJson("RoleStarData", &roleStarPool)
	utils.Json().RegistJson("StarData", &StarPool)
	utils.Json().RegistJson("ExpData", &expPool)
	utils.Json().RegistJson("SignAwardData", &signAwardPool)
	utils.Json().RegistJson("ConstantConfig", &constantConfig)
	utils.Json().RegistJson("WeaponStrengData", &weaponStrengPool)
	utils.Json().RegistJson("DropData", &dropPool)
	utils.Json().RegistJson("WeaponStarUpData", &weaponStarUp)
	utils.Json().RegistJson("SkillData", &skillDataPool)
	utils.Json().RegistJson("GoodsData", &goodsPool)
	utils.Json().RegistJson("LdtEveryAwardData", &everyAwardPool)
	utils.Json().RegistJson("LdtLevelAwardData", &ldtLevelAwardDataPool)
	utils.Json().RegistJson("mail", &mailPool)
	utils.Json().RegistJson("PvpEveryAwardData", &pvpEveryAward)
	utils.Json().RegistJson("PvpLevelAwardData", &pvpLevelAward)
	utils.Json().RegistJson("PvpLevelData", &pvpLevelData)

	utils.Json().Fill()
}
