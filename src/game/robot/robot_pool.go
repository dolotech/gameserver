package robot

import (
	"gameserver/config"
	"gameserver/model"
	"gameserver/online"
	"gameserver/utils"
	"gameserver/utils/log"
	"gameserver/utils/socket/server"
	"math/rand"
	"sort"
	"strconv"
)

var robotPool *RobotPool

type RobotPool struct {
	pool map[uint]server.Session
}

func Pool() *RobotPool {
	return robotPool
}

func (this *RobotPool) Get(uid uint) server.Session {
	return this.pool[uid]
}

// 提供竞技塔使用
func (this *RobotPool) Rank(rank int) *model.Player {
	for _, v := range this.pool {
		player := v.UserData().(*online.UserData).GetPlayer()
		r := player.WeeklyRank
		if r == rank {
			return &player
		}
	}
	return nil
}

// 提供竞技塔使用
func (this *RobotPool) Match(rank []int, power, count int) []*model.Player {
	arr := make([]*model.Player, 0, 5)
	var score float32 = 1.1
	for ; score < 2; {
		for _, v := range this.pool {
			player := v.UserData().(*online.UserData).GetPlayer()
			exist := false
			for _, ex := range rank {
				if player.WeeklyRank == ex {
					exist = true
				}
			}
			if exist {
				continue
			}
			var t = false
			if power > player.Power && int(float32(player.Power)*score) >= power {
				t = true
			} else if power < player.Power && int(float32(power)*score) >= player.Power {
				t = true
			} else if power == player.Power {
				t = true
			}

			if  t {
				rank = append(rank, player.WeeklyRank)
				arr = append(arr, &player)
			}

			if len(arr) == count {
				return arr
			}
		}
		score += 0.1
	}
	return arr
}

// 周赛场虚拟玩家使用
func (this *RobotPool) Search(playerId uint, power int) server.Session {
	player := model.Player{}
	err := player.GetByPowerBetween(playerId, power-power/10, power+power/5)
	log.Error(err)
	if err == nil {
		ps := model.Partners{}
		err := ps.Get(player.PlayerId)
		if err == nil {
			player.PlayerId = uint(rand.Intn(49998) + 1)

			sess := &RobotSession{}
			sess.SetUid(player.PlayerId)
			u := online.NewUserData(player.PlayerId)
			sess.SetUserData(u)
			player.Nickname = utils.GetName()
			player.Avatar = rand.Intn(20)
			u.SetPlayer(player)
			u.SetPartners(ps)
			return sess
		} else {
			log.Error(err)
		}
	} else {
		log.Error(err)
	}

	var score float32 = 1.1
	for ; score < 1.6; {
		for _, v := range this.pool {
			player := v.UserData().(*online.UserData).GetPlayer()
			if power > player.Power && int(float32(player.Power)*score) >= power {
				return cloneSession(v)
			} else if power < player.Power && int(float32(power)*score) >= player.Power {
				return cloneSession(v)
			} else if power == player.Power {
				return cloneSession(v)
			}
		}
		score += 0.1
	}

	return nil
}

func cloneSession(s server.Session) server.Session {
	user := s.UserData().(*online.UserData)
	player := user.GetPlayer()
	ps := user.GetPartners()

	sess := &RobotSession{}
	sess.SetUid(player.PlayerId)
	u := online.NewUserData(player.PlayerId)
	sess.SetUserData(u)
	player.Nickname = utils.GetName()
	player.Avatar = rand.Intn(20)
	u.SetPlayer(player)
	u.SetPartners(ps)
	return sess
}

func Robot() {
	robotPool = &RobotPool{
		pool: make(map[uint]server.Session, 1000),
	}

	robots := config.Robots()
	setPartners := make(map[uint]*model.Partners, len(*robots))
	playerArr := make(SortPlayerList, 0, len(*robots))
	for k, v := range *robots {
		pid, _ := strconv.Atoi(k)
		playerId := uint(pid)
		p := &model.Player{}

		p.PlayerId = playerId
		p.Robot = 1
		par1 := &model.Partner{PlayerId: playerId}
		par1.InitRobot(10001, v.Partner1)

		par2 := &model.Partner{PlayerId: playerId}
		par2.InitRobot(10002, v.Partner2)

		par3 := &model.Partner{PlayerId: playerId}
		par3.InitRobot(10003, v.Partner3)

		par4 := &model.Partner{PlayerId: playerId}
		par4.InitRobot(10004, v.Partner4)

		par5 := &model.Partner{PlayerId: playerId}
		par5.InitRobot(10005, v.Partner5)
		ps := &model.Partners{par1, par2, par3, par4, par5}
		p.Power = ps.CalcuPower(p.Level)
		utils.StructAtoB(p, v)
		playerArr = append(playerArr, p)
		setPartners[playerId] = ps
	}

	sort.Sort(&playerArr)
	for k, v := range playerArr {
		sess := &RobotSession{}
		sess.SetUid(v.PlayerId)
		v.WeeklyRank = k + 1
		u := online.NewUserData(v.PlayerId)
		sess.SetUserData(u)
		u.SetPlayer(*v)
		u.SetPartners(*setPartners[v.PlayerId])
		robotPool.pool[v.PlayerId] = sess
	}
}
