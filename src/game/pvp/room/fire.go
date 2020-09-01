package room

import (
	"gameserver/config"
	"gameserver/game/data"
	"gameserver/game/robot"
	"gameserver/model"
	"gameserver/protocol/push"
	"gameserver/utils/log"
	"math/rand"
	"strconv"
)

func Fire(bp robot.BPlayer,r robot.Room, FighterIndex, SkillId, AttackRatio int) {
	partners := bp.GetPartners()
	SkillId = GetSkillData(partners, FighterIndex, SkillId)
	skillData, ok := config.Skill().Get(strconv.Itoa(SkillId))
	if !ok {
		log.Error("skill is null", SkillId)
		return
	}
	opp := r.Opp(bp.UId())
	msg := &push.OnUpdateFightAction{}
	targets := GenTarget(&skillData, FighterIndex)
	msg.Defenders = CalcHarm(targets, &skillData, AttackRatio, bp, opp)

	msg.FightIndex = FighterIndex
	msg.SkillId = SkillId

	r.Push(msg)
	if winner := r.CheckEnd(); winner != 0 {
		r.GameOver(winner, 0)
		robot.Get().Del(r.RId())
	}
}

//获取技能数据
//skillId为0是普攻，1是武器技能
func GetSkillData(partners *model.Partners, FighterIndex, SkillId int) int {
	if SkillId == 0 || SkillId == 1 {
		position := (FighterIndex % 6) + FighterIndex/6 // 攻击者的位置
		for _, v := range *partners {
			if (v.PartnerId % 10000) == position {
				weapon := config.Weapon().Get(strconv.Itoa(v.WeaponId))
				if SkillId == 0 {
					return weapon.BasisSkill
				} else {
					return weapon.SpecialSkill
				}
				break
			}
		}
	}
	return 0
}

//产生技能施放的目标
func GenTarget(skillData *config.SkillData, FighterIndex int) []int {
	targets := make([]int, 0, 5)
	if skillData.AtkTarget == data.ENEMY {
		for i := 0; i < skillData.AtkNums; i++ {
			r := 11 - (int(rand.Int31n(5)+1) + int(FighterIndex/6)*5) // 随机生成一个敌人
			targets = append(targets, r)
		}
	} else if skillData.AtkTarget == data.FRIEND {
		for i := 0; i < skillData.AtkNums; i++ {
			r := int(rand.Int31n(5)+1) + int(FighterIndex/6)*5 // 随机生成一个队友
			targets = append(targets, r)
		}
	} else if skillData.AtkTarget == data.ME {
		targets = append(targets, FighterIndex)
	}
	return targets
}

//todo 对每个受击者计算防御和伤害
func CalcHarm(
	targets []int,
	skillData *config.SkillData,
	AttackRatio int,
	player, opp robot.BPlayer) []push.Defender {
	defenders := make([]push.Defender, 0, len(targets))

	harm := float32(player.Power()) // 首先获取团队的攻击力
	var effect = skillData.DamageAtk[0]
	harm = effect.Value + harm*effect.Per        // 计算技能加成：先加上技能基础伤害，再乘以系数
	harm = harm * data.AttackRatios[AttackRatio] // 判断爆击情况

	for _, target := range targets {
		defender := push.Defender{AttackRatio: AttackRatio}
		var Hp, Dp float32
		if skillData.DamageType == data.PHY {
			if player.Dp() > 0 {
				Hp = harm*0.4 + 0.5
				Dp = harm*0.6 + 0.5
			} else {
				Hp = harm
			}
			if opp.HarmHp(int32(-Hp)) <= 0 {
				defender.Lethal = 1
			}
			if Hp > 0 {
				opp.HarmDp(int32(-Dp))
			}
		} else if skillData.DamageType == data.MAGIC {
			if player.Dp() > 0 {
				Hp = harm*0.6 + 0.5
				Dp = harm*0.4 + 0.5
			} else {
				Hp = harm
			}
			if opp.HarmHp(int32(-Hp)) <= 0 {
				defender.Lethal = 1
			}
			if Hp > 0 {
				opp.HarmDp(int32(-Dp))
			}
		} else if skillData.DamageType == data.ADD_HP {
			Hp = harm
			player.HarmHp(int32(Hp))
		} else if skillData.DamageType == data.ADD_DE {
			Dp = harm
			player.HarmDp(int32(Dp))
		}
		defender.Hp, defender.Dp = int(Hp), int(Dp)
		defender.FightIndex = target
		defenders = append(defenders, defender)
	}
	return defenders
}
