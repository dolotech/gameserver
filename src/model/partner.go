package model

import (
	"gameserver/config"
	"gameserver/utils"
	"gameserver/utils/db"
	"gameserver/utils/log"
	"strconv"
)

// 队伍角色   gameserver
// [1,1,120001,0,1],
//[40,4,110006,120,6]
//{trainLevel=1,trainStar=1,eqId=12001,eqStrength=5,eqStar=3}
type Partner struct {
	PlayerId  uint `gorm:"type:int(11);column:playerId;unique_index:index_playerId" json:"playerId" `
	PartnerId int  `gorm:"type:int(11);column:partnerId;unique_index:index_playerId" json:"partnerId" `
	Level     int  `gorm:"type:int(6);column:level;COMMENT:'角色训练等级'" json:"level" `
	Point     int  `gorm:"type:int(8);column:point" json:"point" `
	PStar     int  `gorm:"type:int(8);column:pstar;COMMENT:'角色星级'" json:"pstar" `
	WeaponId  int  `gorm:"type:int(8);column:weaponId" json:"weaponId" `
	WStar     int  `gorm:"type:int(8);column:wstar;COMMENT:'武器星级'" json:"wstar" `
	Strength  int  `gorm:"type:int(8);column:strength;COMMENT:'武器强化等级'" json:"strength" `
}

type Partners []*Partner

func (this *Partners) Clone() Partners {
	ps := make(Partners, len(*this))
	for k, v := range *this {
		p := &Partner{}
		utils.StructAtoB(p, v)
		ps[k] = p
	}
	return ps
}

//角色战力因素（非数据表）
type Power struct {
	Phy   int //物理攻击
	Mag   int //魔法攻击
	Dp    int //防御值
	Hp    int //血量
	Score int //战力综合评分值
}

//所有角色战力
type Powers []Power

func (this *Partner) InitRobot(id int, robot []int) {
	//角色属性顺序[trainLevel,trainStar,eqId,eqStrength,eqStar]
	this.PartnerId = id
	this.Level = robot[0]
	this.PStar = robot[1]
	this.WeaponId = robot[2]
	this.Strength = robot[3]
	this.WStar = robot[4]
}
func (this *Partners) InitPartner(player *Player) {
	rolePool := config.Role()
	for _, id := range rolePool.GetAllId() {
		role := rolePool.Get(id)
		partnerId, _ := strconv.Atoi(id)

		partner := &Partner{
			PlayerId:  player.PlayerId,
			PartnerId: partnerId,
			Level:     1,
			//Exp:       0,
			Point:    0,
			PStar:    1,
			WeaponId: role.WeaponId,
			WStar:    1,
			Strength: 1,
		}
		*this = append(*this, partner)
		if err := partner.Add(); err != nil {
			log.Error("init partner error,%s", err)
		}
	}
}

func (this *Partner) Add() error {
	return db.Get().Create(this).Error
}

func (this *Partner) GetSingle() error {
	return db.Get().Model(this).Where("playerID = ? and partnerId = ?", this.PlayerId, this.PartnerId).Find(this).Error
}

func (this *Partner) UpdateStrength(strengthLevel int) error {
	m := &Partner{Strength: strengthLevel}
	return db.Get().Model(this).Where("playerID = ? and partnerId = ?", this.PlayerId, this.PartnerId).Updates(m).Error
}

func (this *Partner) UpdateWStar(star int) error {
	m := &Partner{WStar: star}
	return db.Get().Model(this).Where("playerID = ? and partnerId = ?", this.PlayerId, this.PartnerId).Updates(m).Error
}

func (this *Partner) UpdateWeapon(weaponId int) error {
	m := &Partner{WeaponId: weaponId, WStar: 1} //新武器star重置为1
	return db.Get().Model(this).Where("playerID = ? and partnerId = ?", this.PlayerId, this.PartnerId).Updates(m).Error
}

func (this *Partner) UpdateTrain() error {
	m := map[string]interface{}{
		"point": this.Point,
		"level": this.Level,
	}
	return db.Get().Model(this).Where("playerID = ? and partnerId = ?", this.PlayerId, this.PartnerId).Updates(m).Error
}

func (this *Partner) UpdatePStar() error {
	m := &Partner{PStar: this.PStar}
	return db.Get().Model(this).Where("playerID = ? and partnerId = ?", this.PlayerId, this.PartnerId).Updates(m).Error
}

func (this *Partners) Get(playerID uint) error {
	return db.Get().Order("partnerId ASC").Model(this).Where("playerID = ?", playerID).Find(this).Error
}

func (this *Partners) TotalHP(level int) int {
	basic := this.GetBasicPower(level)
	role := this.RolesPower()
	weapon := this.GetWeaponPower()
	return basic.Hp + role.Hp + weapon.Hp
}

func (this *Partners) TotalDP(level int) int {
	basic := this.GetBasicPower(level)
	role := this.RolesPower()
	weapon := this.GetWeaponPower()
	return basic.Dp + role.Dp + weapon.Dp
}

func (this *Partners) CalcuPower(level int) int {
	basic := this.GetBasicPower(level)
	role := this.RolesPower()
	weapon := this.GetWeaponPower()
	totalPower := &Power{
		Phy: basic.Phy + role.Phy + weapon.Phy,
		Mag: basic.Mag + role.Mag + weapon.Mag,
		Dp:  basic.Dp + role.Dp + weapon.Dp,
		Hp:  basic.Hp + role.Hp + weapon.Hp,
	}
	totalPower.Score = totalPower.Mag + totalPower.Phy + int(float64(totalPower.Dp+totalPower.Hp)*0.2)
	return totalPower.Score
}

//角色基础数据以及自然等级加成
func (this *Partners) GetBasicPower(level int) *Power {
	totalBasic := &Power{Phy: 0, Mag: 0, Dp: 0, Hp: 0}
	for _, partner := range *this {
		basic := getBasicPowerByPartner(partner.PartnerId, level)
		totalBasic.Phy += basic.Phy
		totalBasic.Mag += basic.Mag
		totalBasic.Dp += basic.Dp
		totalBasic.Hp += basic.Hp
	}
	return totalBasic
}

func getBasicPowerByPartner(partnerId int, level int) *Power {
	rolePool := config.Role()
	basic := &Power{}
	for _, id := range rolePool.GetAllId() {
		rid, _ := strconv.Atoi(id)
		if partnerId == rid {
			role := rolePool.Get(id)
			basic.Phy = getAttribute(level, role.Physical, role.PhysicalGrowth)
			basic.Mag = getAttribute(level, role.Magic, role.MagicGrowth)
			basic.Dp = getAttribute(level, role.Defense, role.DefenseGrowth)
			basic.Hp = getAttribute(level, role.Life, role.LifeGrowth)
		}
	}
	return basic
}

func getAttribute(level int, init int, growth int) int {
	return init + (level-1)*growth
}

//角色训练与升星加成
func (this *Partners) RolesPower() *Power {
	roleTotal := &Power{Phy: 0, Mag: 0, Dp: 0, Hp: 0}
	for _, partner := range *this {
		role := partner.RolePower()
		roleTotal.Phy += role.Phy
		roleTotal.Mag += role.Mag
		roleTotal.Dp += role.Dp
		roleTotal.Hp += role.Hp
	}
	return roleTotal
}

func (partner *Partner) RolePower() *Power {
	role := &Power{}
	constCfg := config.GetConstantCfg()
	trainData := config.RoleTrain().Get(strconv.Itoa(partner.Level))
	starData := config.RoleStar().Get(strconv.Itoa(partner.PStar))
	switch partner.PartnerId {
	case constCfg.PARTNERTYPE.WARRIOR:
		role.Hp = trainData.SoldierLife + starData.SoldierStarLife
		role.Dp = trainData.SoldierDefense + starData.SoldierStarDefense
		role.Phy = trainData.SoldierPhysical + starData.SoldierStarPhysical
		role.Mag = 0
	case constCfg.PARTNERTYPE.KNIGHT:
		role.Hp = trainData.KnightLife + starData.KnightStarLife
		role.Dp = trainData.KnightDefense + starData.KnightStarDefense
		role.Phy = trainData.KnightPhysical + starData.KnightStarPhysical
		role.Mag = 0
	case constCfg.PARTNERTYPE.MAGICIAN:
		role.Hp = trainData.MageLife + starData.MageStarLife
		role.Dp = trainData.MageDefense + starData.MageStarDefense
		role.Phy = 0
		role.Mag = trainData.MageMagic + starData.MageStarMagic
	case constCfg.PARTNERTYPE.PRIEST:
		role.Hp = trainData.PriestLife + starData.PriestStarLife
		role.Dp = trainData.PriestDefense + starData.PriestStarDefense
		role.Phy = 0
		role.Mag = trainData.PriestMagic + starData.PriestStarMagic
	case constCfg.PARTNERTYPE.ARCHER:
		role.Hp = trainData.ShooterLife + starData.ShooterStarLife
		role.Dp = trainData.ShooterDefense + starData.ShooterStarDefense
		role.Phy = trainData.ShooterPhysical + starData.ShooterStarPhysical
		role.Mag = 0
	}
	return role
}

//武器基础属性以及升级加成
func (this *Partners) GetWeaponPower() *Power {
	weaponTotal := &Power{Phy: 0, Mag: 0, Dp: 0, Hp: 0}
	for _, partner := range *this {
		weapon := partner.WeaPonPower()
		weaponTotal.Phy += weapon.Phy
		weaponTotal.Mag += weapon.Mag
		weaponTotal.Dp += weapon.Dp
		weaponTotal.Hp += weapon.Hp
	}
	return weaponTotal
}

func (partner *Partner) WeaPonPower() *Power {
	weapon := &Power{}
	basic := config.Weapon().Get(strconv.Itoa(partner.WeaponId))
	strength := config.WeaponStreng().GetByTypeAndLevel(basic.ResType, partner.Strength)
	star := config.WeaponStarUp().GeyByIdAndStar(partner.WeaponId, partner.PStar)
	weapon.Hp = basic.HP + strength.HP + star.HP
	weapon.Dp = basic.DF + strength.DF + star.DF
	weapon.Phy = basic.PHY + strength.PHY + star.PHY
	weapon.Mag = basic.MAG + strength.MAG + star.MAG
	return weapon
}
