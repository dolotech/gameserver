package mail

import (
	"encoding/json"
	"gameserver/config"
	"gameserver/game/award"
	"gameserver/model"
	proto "gameserver/protocol"
	"gameserver/protocol/push"
	"gameserver/protocol/route"
	"gameserver/utils/log"
	msg "gameserver/utils/socket/message"
	"gameserver/utils/socket/server"
	"reflect"
	"strconv"
)

const (
	NEW      = 1 //未发送
	SENDED   = 2 //已发送
	RECEIVED = 3 //已领取
)

func init(){
	msg.GetMsg().Reg(route.RecvEmailItems, &proto.RecvEmailItems{}, &proto.RecvEmailItemsed{}, recvEmailItems)
}

func recvEmailItems(sess server.Session, req *proto.RecvEmailItems, resp *proto.RecvEmailItemsed){
	resp.Code = proto.OK
	mailM := &model.Emails{}
	err := mailM.GetByEmailIds(req.EmailId, sess.UId(), SENDED)
	if err != nil || reflect.DeepEqual(mailM, &model.Email{}) {
		resp.Code = proto.FA_INVALID_EMAIL_ID
		log.Error(err)
		return
	}
	items := []proto.Item{}
	for _, m := range *mailM {
		item := []proto.Item{}
		if err := json.Unmarshal(m.Items, &item);err != nil {
			resp.Code = proto.FAIL
			log.Error(err)
			return
		}
		items = append(items, item...)
	}

	//领取物品
	player := &model.Player{PlayerId:sess.UId()}
	player.GetAwardAttr()
	awards := &award.Awards{}
	awards.ConvertAwards(items)
	awards.SaveAwards(player)
	if awards.OnBag(){
		(&push.OnUpdateBag{}).Push(player.PlayerId)
	}
	if awards.OnPlayer() {
		(&push.OnUpdatePlayer{}).Push(player)
	}
	(&model.Email{}).UpdateOneStatus(req.EmailId, RECEIVED)
}

type Mail struct {
	EmailId  uint
	PlayerId uint
	EndTime  int
	TypeId   int
	Status   int
	Head     string
	Text     string
	Items    []proto.Item
}

func (this *Mail) MailByType(typeId int, playerId uint) {
	this.PlayerId = playerId
	this.TypeId = typeId
	mailCfg := config.MailPoll().Get(strconv.Itoa(this.TypeId))
	this.Head = mailCfg.Head
	this.Text = mailCfg.Text
	this.Status = NEW
	awards := &award.Awards{}
	awards.Drop(mailCfg.DropId)
	this.Items =  awards.ConvertItem()
	this.AddMail()
}

func (this *Mail) TaskMail(tasks []config.TaskData, playerId uint) {
	this.PlayerId = playerId
	this.TypeId = config.GetConstantCfg().MAILTYPE.TASK
	mailCfg := config.MailPoll().Get(strconv.Itoa(this.TypeId))
	this.Head = mailCfg.Head
	this.Text = mailCfg.Text
	this.Status = NEW
	for _, t := range tasks {
		awards := &award.Awards{}
		awards.Drop(t.DropId)
		this.Items = append(this.Items, awards.ConvertItem()...)
	}
	this.AddMail()
}

func (this *Mail) RookieMail(playerId uint) {
	typeId := config.GetConstantCfg().MAILTYPE.ROOKIE
	this.MailByType(typeId, playerId)
}


func (this *Mail)SendNewMails(playerId uint,after ...bool){
	mails := &model.Emails{}
	if err := mails.GetNewMails(playerId, SENDED);err != nil{
		log.Error(err)
		return
	}
	if len(*mails) <= 0 {
		return
	}
	pushMails := make(push.OnUpdateEmails, 0, len(*mails))
	for _, m := range *mails {
		mailCfg := config.MailPoll().Get(strconv.Itoa(m.TypeId))
		mail := &push.OnUpdateEmail{
			EmailId:m.EmailId,
			Head:mailCfg.Head,
			Text:mailCfg.Text,
		}
		if len(m.Items) > 0{
			if err := json.Unmarshal(m.Items, &mail.Items);err != nil {
				log.Error(err)
			}
		}
		pushMails = append(pushMails, *mail)
	}
	pushMails.Push(playerId,after ...)
	mails.UpdateStatus(playerId, SENDED)
}


func (this *Mail) AddMail() {
	itemStr, _ := json.Marshal(this.Items)
	mailM := &model.Email{
		PlayerId:this.PlayerId,
		TypeId:this.TypeId,
		Status:this.Status,
		Items:itemStr,
	}
	if err := mailM.AddMail();err != nil{
		log.Error(err)
		return
	}
	this.EmailId = mailM.EmailId
	if len(this.Items) <= 0 {
		return
	}
}
