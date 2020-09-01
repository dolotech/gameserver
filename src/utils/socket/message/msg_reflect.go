package message

import (
	"gameserver/utils/log"
	"reflect"
	"strconv"
	"sync"
)

var msgMgrSingleton  *msgMgr
	var oneMmsg sync.Once
func GetMsg() *msgMgr {
	// 保证对象一定会被创建，并且在返回前第二个调用阻塞等待第第一返回
	oneMmsg.Do(func() {
		msgMgrSingleton = &msgMgr{
			msgMap:     make(map[string]*MsgInfo, 160),
			pushMsgMap: make(map[reflect.Type]string, 160),
		}
	})
	return msgMgrSingleton
}

type msgMgr struct {
	msgMap     map[string]*MsgInfo
	pushMsgMap map[reflect.Type]string
}

type MsgInfo struct {
	Route       string
	RouteId     uint16
	Cb          reflect.Value
	MsgReqType  reflect.Type
	MsgRespType reflect.Type
}

func (m *msgMgr) Reg(route string, msgReq interface{}, msgResp interface{}, f interface{}) {
	msgType := reflect.TypeOf(msgReq)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Error("message request pointer required")
		return
	}

	if reflect.TypeOf(f).Kind() != reflect.Func {
		log.Error("must be function")
		return
	}

	if _, ok := m.msgMap[route]; ok {
		log.Warning("route %s is already registered", route)
		return
	}

	msgRespType := reflect.TypeOf(msgResp)
	if msgRespType == nil || msgRespType.Kind() != reflect.Ptr {
		log.Error("message response pointer required")
		return
	}

	i := &MsgInfo{}
	i.RouteId = uint16(len(m.msgMap) + 1)
	i.Route = route
	i.MsgReqType = msgType.Elem()
	i.MsgRespType = msgRespType.Elem()
	i.Cb = reflect.ValueOf(f)
	m.msgMap[route] = i
}

func (m *msgMgr) GetPushMsg(msg interface{}) string {
	msgRespType := reflect.TypeOf(msg)
	if msgRespType == nil || msgRespType.Kind() != reflect.Ptr {
		log.Error("message msg pointer required")
		return ""
	}
	respType:=m.pushMsgMap[msgRespType]
	if respType == "" {
		log.Error("push route %s not registered",msgRespType)
	}
	return respType
}
func (m *msgMgr) RegPush(route string, msg interface{}) {
	if _, ok := m.msgMap[route]; ok {
		log.Warning("route %s is already registered", route)
		return
	}
	msgRespType := reflect.TypeOf(msg)
	if msgRespType == nil || msgRespType.Kind() != reflect.Ptr {
		log.Error("message msg pointer required")
		return
	}
	i := &MsgInfo{}
	i.Route = route
	i.MsgRespType = msgRespType
	i.RouteId = uint16(len(m.msgMap) + 1)
	m.msgMap[route] = i
	m.pushMsgMap[msgRespType] = i.Route
}

func (m *msgMgr) GetMsgByRoute(route string) *MsgInfo {
	return m.msgMap[route]
}

func (m *msgMgr) GetMsgMap() map[string]string {
	ms := make(map[string]string, len(m.msgMap))
	for _, v := range m.msgMap {
		ms[strconv.Itoa(int(v.RouteId))] = v.Route
	}
	return ms
}
