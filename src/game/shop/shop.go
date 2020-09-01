package shop

import (
	"gameserver/config"
	"gameserver/game/data"
	"gameserver/game/task"
	"gameserver/model"
	proto "gameserver/protocol"
	"gameserver/protocol/push"
	"gameserver/protocol/route"
	"gameserver/protocol/shop"
	"gameserver/utils/log"
	msg "gameserver/utils/socket/message"
	"gameserver/utils/socket/server"
	"strconv"
)

// 消息注册
func init() {
	msg.GetMsg().Reg(route.ShopBuyItem, &shop.ShopBuyItem{}, &shop.ShopBuyItemed{}, shopBuyItemCb)
	msg.GetMsg().Reg(route.BuyDiamond, &shop.BuyDiamond{}, &shop.BuyDiamonded{}, buyDiamondCb)
	msg.GetMsg().Reg(route.BuyStamina, &shop.BuyStamina{}, &shop.BuyStaminaed{}, buyStaminaCb)
	msg.GetMsg().Reg(route.SellItem, &shop.SellItem{}, &shop.SellItemed{}, sellItemCb)
}

//请求物品出售
func sellItemCb(sess server.Session, req *shop.SellItem, resp *shop.SellItemed) {
	resp.Code = proto.OK
	p := &model.Player{PlayerId: sess.UId()}
	var config = config.Goods().Get( strconv.Itoa(req.ItemID))
	if config == nil {
		log.Error("sellItemCb FA_INVALID_SHOP_DATA_ID %d",req.ItemID)
		resp.Code = proto.FAIL		//FA_INVALID_SHOP_DATA_ID
		return
	}

	bag := &model.Item{
		PlayerId: sess.UId(),
		ItemId:   req.ItemID,
	}
	itmeCount, _ := bag.GetCount()
	if itmeCount < req.Count {
		log.Warning("has not enough itemId! itemId: %d", req.ItemID)
		resp.Code = proto.FA_NOT_ENOUGH_ITEMS
		return
	}
	bag.Reduce(req.Count)
	goldIncome := config.Sell * req.Count
	p.UpdateGold(goldIncome)
	//log.Info("sellItem, playerId: %d, itemId: %d, count: %d, goldIncome: %d", p.PlayerId, req.ItemID, req.Count, goldIncome)
}

// 请求购买体力
func buyStaminaCb(sess server.Session, req *shop.BuyStamina, resp *shop.BuyStaminaed) {
	resp.Code = proto.OK
	p := &model.Player{PlayerId: sess.UId()}

	if p.GetOnUpdatePlayer() != nil {
		resp.Code = proto.FAIL
		return
	}

	var paramData = config.Param()
	if p.Stamina >= paramData.StaminaMax {
		resp.Code = proto.FA_MAX_STAMINA
		return
	}

	if p.Diamond < paramData.StaminaPrice {
		resp.Code = proto.FA_NOT_ENOUGH_DIAMOND
		return
	}

	if p.UpdateDiamond(-paramData.StaminaPrice) != nil {
		resp.Code = proto.FA_NOT_ENOUGH_DIAMOND
		return
	}

	if err := p.UpdateStamina(paramData.StaminaNum); err != nil {
		resp.Code = proto.FAIL
		return
	}
	//更新任务购买体力次数
	myTasks := &task.Task{PlayerId: p.PlayerId}
	myTasks.UpdateBuyStamina()
	myTasks.UpdateSpendDiamondDaily()
	myTasks.PushTasks()
}

// 请求购买钻石
func buyDiamondCb(sess server.Session, req *shop.BuyDiamond, resp *shop.BuyDiamonded) {
	resp.Code = proto.OK
	product := config.Diamond().Get(req.Platform, req.ProductId)
	if product == nil {
		log.Error("商品不存在")
		resp.Code = proto.FAIL
		return
	}

	p := &model.Player{PlayerId: sess.UId()}

	if err:=p.UpdateDiamond(product.Count);err!=nil{
		log.Error(err)
		resp.Code = proto.FAIL
		return
	}

	if err:=p.GetOnUpdatePlayer();err!=nil{
		log.Error(err)
		return
	}
	(&push.OnUpdatePlayer{}).Push(p)
}

// 请求购买物品
func shopBuyItemCb(sess server.Session, req *shop.ShopBuyItem, resp *shop.ShopBuyItemed) {
	if req.Count == 0 || req.Type == 0 || req.ID == 0 {
		resp.Code = proto.FAIL
		return
	}
	resp.Code = proto.OK
	s := config.Shop().Get(strconv.Itoa(req.ID))
	if s == nil {
		log.Error("商品不存在")
		resp.Code = proto.FA_INVALID_TOOLS_INDEX
		return
	}

	p := &model.Player{PlayerId: sess.UId()}
	if p.GetDiamond() != nil {
		resp.Code = proto.FAIL
		return
	}
	total := s.Price * req.Count
	if p.Diamond < total {
		resp.Code = proto.FA_NOT_ENOUGH_DIAMOND
		return
	}

	if p.UpdateDiamond(-total) != nil {
		resp.Code = proto.FA_NOT_ENOUGH_DIAMOND
		return
	}
	if s.Type == 1 { //购买玩家金币，扣除玩家钻石
		p.UpdateGold(s.Count * req.Count)
	} else if s.Type == 2 { //购买道具，扣除玩家钻石
		item := &model.Item{
			PlayerId: sess.UId(),
			ItemId:   s.ShopId,
		}
		item.Add(req.Count * s.Count)
		(&push.OnUpdateBag{}).Push(sess.UId())

	} else if s.Type == 3 { //购买道具，扣除玩家钻石
		list := config.Drop().GetByDropId(s.ShopId)

		if len(list) > 0 {
			gold := 0
			honor := 0
			for _, v := range list {
				if v.ItemId == data.GOLD { // 金币
					gold += v.Count
				} else if v.ItemId == data.HONOR { // 荣誉
					honor += v.Count
				} else {
					item := &model.Item{
						PlayerId: sess.UId(),
						ItemId:   v.ItemId,
					}
					item.Add(v.Count * req.Count)
				}

			}
			if gold > 0 {
				p.UpdateGold(gold * req.Count)
			}
			if honor > 0 {
				p.UpdateHonor(honor * req.Count)
			}
			(&push.OnUpdateBag{}).Push(sess.UId())
		}
	}

	if p.GetDiamondGoldHonor() == nil {
		resp.Diamond = p.Diamond
		resp.Gold = p.Gold
		resp.Honor = p.Honor
	}
	//任务
	task := &task.Task{PlayerId:p.PlayerId}
	task.UpdateSpendDiamondDaily()
	task.PushTasks()
}
