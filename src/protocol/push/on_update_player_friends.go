package push

import (
	"gameserver/model"
	"gameserver/online"
	"gameserver/protocol"
	"gameserver/protocol/friend"
	"gameserver/utils"
	"gameserver/utils/log"
)

//更新好友模块信息PUSH
type OnUpdatePlayerFriends struct {
	Code       int             `json:"code"`
	Friends    []friend.Friend `json:"friends"`    //所有好友信息
	Applicants []friend.Friend `json:"applicants"` //申请人信息
}

func (this *OnUpdatePlayerFriends) Push(playerid uint,after ...bool) error {
	this.Code = protocol.OK
	fs := model.Friends{}
	if err := fs.Get(playerid); err != nil {
		log.Error(err)
		return err
	}
	if len(fs) == 0 {
		log.Warning("OnUpdatePlayerFriends Push len(fs) %d",len(fs))
		return nil
	}
	ids := make([]uint, len(fs))
	for k, v := range fs {
		ids[k] = v.FriendId
	}
	ps := model.Players{}
	if err := ps.GetIn(ids); err != nil {
		log.Error(err)
		return err
	}

	for _, v := range ps {
		friend := friend.Friend{}
		utils.StructAtoB(&friend, v)
		isOn := online.Get().Online(v.PlayerId)
		if isOn {
			friend.IsOnline = 1
		} else {
			friend.IsOnline = 0
		}
		for _, v1 := range fs {
			if v1.FriendId == v.PlayerId {
				friend.Statue = v1.Status
			}
		}

		if friend.Statue == model.Apply {
			this.Applicants = append(this.Applicants, friend)
		} else {
			this.Friends = append(this.Friends, friend)
		}
	}
	online.Get().Push(playerid, this,after...)
	return nil
}
