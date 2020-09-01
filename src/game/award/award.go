package award

import (
	"gameserver/config"
	"gameserver/model"
	proto "gameserver/protocol"
	"math/rand"
	"strconv"
)

type Awards struct {
	Ready    bool
	Items    map[int]int //map[ItemId]Count
	Exp      int
	Gold     int
	Diamond  int
	Stamina  int
	Honor    int
	PvpScore int
}

func (this *Awards) Drop(dropId int) {
	if this.Items == nil {
		this.Items = make(map[int]int)
	}
	currDrops := config.Drop().GetByDropId(dropId)
	isDrop := false
	total := 0
	var subDropIds []int
	constCfg := config.GetConstantCfg()
	for _, v := range currDrops {
		if v.Probability != -1 {
			total += v.Probability
		}
	}
	//根据配置表计算掉落情况
	for _, v := range currDrops {
		r := 0
		if total > 1 {
			r = rand.Intn(total-1) + 1
		}
		if v.Probability == -1 || (!isDrop && r < v.Probability) {
			if v.Probability != -1 { //-1必掉落
				isDrop = true
				total -= v.Probability
			}
			if v.Droptype == constCfg.DROPTYPE.OTHERDROP { //1正常掉落item，2空item，3其他dropid的item
				subDropIds = append(subDropIds, v.ItemId)
			} else if v.Droptype == constCfg.DROPTYPE.NORMAL {
				switch v.ItemId {
				case constCfg.ITEMTYPE.EXP:
					this.Exp += v.Count
				case constCfg.ITEMTYPE.GOLD:
					this.Gold += v.Count
				case constCfg.ITEMTYPE.DIAMOND:
					this.Diamond += v.Count
				case constCfg.ITEMTYPE.STAMINA:
					this.Stamina += v.Count
				case constCfg.ITEMTYPE.HONOR:
					this.Honor += v.Count
				case constCfg.ITEMTYPE.PVPSCORE:
					this.PvpScore += v.Count
				}
				if v.ItemId > constCfg.ITEMTYPE.ITEM {
					if _, ok := this.Items[v.ItemId]; ok {
						this.Items[v.ItemId] += v.Count
					} else {
						this.Items[v.ItemId] = v.Count
					}
				}
			}
		}
	}
	for _, id := range subDropIds {
		this.Drop(id)
	}
}

func (this *Awards) OnPlayer() bool {
	return this.Exp > 0 || this.Gold > 0 || this.Diamond > 0 || this.Stamina > 0 || this.Honor > 0 || this.PvpScore > 0
}
func (this *Awards) OnBag() bool {
	return len(this.Items) > 0
}
func (this *Awards) OnLevelUp(lv, exp int) bool {
	if this.Exp == 0 {
		return false
	}
	return lv < config.Exp().GetMaxLevel() && exp+this.Exp >= config.Exp().Get(strconv.Itoa(lv)).Exp
}

func (this *Awards) SaveAwards(player *model.Player) {
	if player.ResType == 0 || player.Power == 0 {
		player.GetAwardAttr()
	}

	if this.Exp > 0 {
		currExp := player.Exp
		currLevel := player.Level
		currExp += this.Exp
		con := config.Exp().Get(strconv.Itoa(currLevel))
		player.Exp = currExp
		if currLevel < config.Exp().GetMaxLevel() && currExp >= con.Exp {
			currLevel++
			//角色升级，需要更新战力
			partners := &model.Partners{}
			partners.Get(player.PlayerId)
			player.Power = partners.CalcuPower(currLevel)
		}
		player.Level = currLevel
	}
	if this.Gold > 0 {
		player.Gold += this.Gold
	}
	if this.Diamond > 0 {
		player.Diamond += this.Diamond
	}
	if this.Stamina > 0 {
		player.Stamina += this.Stamina
	}
	if this.Honor > 0 {
		player.Honor += this.Honor
	}
	if this.PvpScore > 0 {
		//这里关于rank可能有别的操作，在pvpUpgrade
		player.PvpScore += this.PvpScore
	}
	player.SetAwardAttr()
	for id, count := range this.Items {
		(&model.Item{PlayerId: player.PlayerId, ItemId: id}).Add(count)
	}
}

func (this *Awards) SaveAwardsById(playerId uint) {
	player := &model.Player{PlayerId: playerId}
	if player.GetAwardAttr() != nil {
		return
	}
	this.SaveAwards(player)
}

func (this *Awards) ConvertItem() (items []proto.Item) {
	constCfg := config.GetConstantCfg()
	if this.Gold > 0 {
		items = append(items, proto.Item{ID: constCfg.ITEMTYPE.GOLD, Count: this.Gold})
	}
	if this.Diamond > 0 {
		items = append(items, proto.Item{ID: constCfg.ITEMTYPE.DIAMOND, Count: this.Diamond})
	}
	if this.Stamina > 0 {
		items = append(items, proto.Item{ID: constCfg.ITEMTYPE.STAMINA, Count: this.Stamina})
	}
	if this.Honor > 0 {
		items = append(items, proto.Item{ID: constCfg.ITEMTYPE.HONOR, Count: this.Honor})
	}
	if this.PvpScore > 0 {
		items = append(items, proto.Item{ID: constCfg.ITEMTYPE.HONOR, Count: this.Honor})
	}
	for id, count := range this.Items {
		items = append(items, proto.Item{ID: id, Count: count})
	}
	return
}

func (this *Awards) ConvertAwards(items []proto.Item) {
	constCfg := config.GetConstantCfg()
	this.Items = make(map[int]int)
	for _, i := range items {
		switch i.ID {
		case constCfg.ITEMTYPE.EXP:
			this.Exp += i.Count
		case constCfg.ITEMTYPE.GOLD:
			this.Gold += i.Count
		case constCfg.ITEMTYPE.DIAMOND:
			this.Diamond += i.Count
		case constCfg.ITEMTYPE.STAMINA:
			this.Stamina += i.Count
		case constCfg.ITEMTYPE.HONOR:
			this.Honor += i.Count
		case constCfg.ITEMTYPE.PVPSCORE:
			this.PvpScore += i.Count
		}
		if i.ID > constCfg.ITEMTYPE.ITEM {
			if _, ok := this.Items[i.ID]; ok {
				this.Items[i.ID] += i.Count
			} else {
				this.Items[i.ID] = i.Count
			}
		}
	}
}
