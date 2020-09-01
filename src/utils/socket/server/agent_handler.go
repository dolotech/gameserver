package server

import (
	"gameserver/utils/log"
	"gameserver/utils/socket/message"
	jsoniter "github.com/json-iterator/go"
	"reflect"
	"runtime/debug"
)

type agentHandler struct {
	agent *agent
}

func NewAgentHandler(a *agent) *agentHandler {
	return &agentHandler{agent: a}
}

func (h *agentHandler) Handle(pkgType byte, body []byte) {
	switch pkgType {
	case message.TYPE_HANDSHAKE:
		h.handleHandshake(body)
	case message.TYPE_HANDSHAKE_ACK:
		h.handleHandshakeAck(body)
	case message.TYPE_HEARTBEAT:
		h.handleHeartbeat(body)
	case message.TYPE_DATA:
		h.handleData(body)
	}
}

func (h *agentHandler) handleHandshake(body []byte) {
	var messageJ map[string]interface{}
	err := jsoniter.Unmarshal(body, &messageJ)
	if err != nil {
		log.Error("handshake decode error: %v", err)
		h.processError(CODE_USE_ERROR)
		return
	}
	if messageJ["user"] == nil {
		log.Error("messageJ[user]")
		h.processError(FA_TOKEN_INVALID)
		h.agent.Close()
		return
	}
	if messageJ["sys"] == nil {
		h.processError(CODE_USE_ERROR)
		return
	}
	if au, ok := messageJ["user"].(string); ok {
		if _, ok := h.agent.event.OnAuth(h.agent, au); ok {
		} else {
			h.processError(FA_TOKEN_INVALID)
			h.agent.Close()
		}
	} else {
		log.Error("FA_TOKEN_INVALID")
		h.processError(FA_TOKEN_INVALID)
		h.agent.Close()
		return
	}
	sys := make(map[string]interface{})
	res := make(map[string]interface{})
	res["code"] = CODE_OK
	res["sys"] = sys
	sys["heartbeat"] = 3 // 心跳间隔
	sys["timeout"] = 6   // 心跳超时
	sys["useDict"] = true
	ms := message.GetMsg().GetMsgMap()
	sys["dict"] = ms //发给前端的路由映射
	bin, _ := jsoniter.Marshal(res)
	pack := message.PoolGet()
	pack.Pack(message.TYPE_HANDSHAKE, bin)
	h.agent.Write(pack)
}

func (h *agentHandler) handleHandshakeAck(body []byte) {
	h.handleHeartbeat(body)
}

func (h *agentHandler) handleHeartbeat(ody []byte) {
	pack := message.PoolGet()
	pack.Pack(message.TYPE_HEARTBEAT, []byte{})
	h.agent.Write(pack)
}

func (h *agentHandler) handleData(body []byte) {
	defer func() {
		if e := recover(); e != nil {
			log.Error(string(debug.Stack()), e)
		}
	}()
	if h.agent.UId() == 0 {
		h.agent.Close()
		return
	}
	m := message.PoolGet()
	err := m.Decode(body)
	if err != nil {
		log.Error(err)
		return
	}
	msgInfo := message.GetMsg().GetMsgByRoute(m.Route)
	if msgInfo == nil {
		log.Error("路由未注册: ", m.Route)
		return
	}
	var req = reflect.New(msgInfo.MsgReqType)
	var resq = reflect.New(msgInfo.MsgRespType)
	err = m.Unmarshal(req.Interface())
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("收到:", m)

	defer h.agent.locker.Unlock()
	h.agent.locker.Lock()
	msgInfo.Cb.Call([]reflect.Value{reflect.ValueOf(h.agent), req, resq})

	resqData, err := jsoniter.Marshal(resq.Interface())
	if err != nil {
		log.Error(err)
		return
	}
	m.Data = resqData
	log.Info("发送:", m)
	m.Data = m.Compress(m.Data)
	m.Type = message.Response
	err = m.Encode()
	if err != nil {
		log.Error(err)
		return
	}
	h.agent.WriteMsg(m)
}

func (h *agentHandler) processError(code int) {
	r := make(map[string]int)
	r["code"] = code
	bin, _ := jsoniter.Marshal(r)
	pack := message.PoolGet()
	pack.Pack(message.TYPE_HANDSHAKE, bin)
	h.agent.Write(pack)
}
