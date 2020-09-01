package server

import (
	"gameserver/utils/log"
	"gameserver/utils/socket/message"
	"gameserver/utils/socket/tcp"
	jsoniter "github.com/json-iterator/go"
	"net"
	"runtime/debug"
	"sync"
	"time"
)

type NetEvent interface {
	OnClose(Session)
	OnAuth(Session, string) (uint, bool)
}

type Session interface {
	SetUid(uint)
	UId() uint
	Push(interface{}, ...bool)
	WriteMsg(*message.Message)
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	UserData() interface{}
	SetUserData(data interface{})
}

type agent struct {
	conn             tcp.Conn
	agentHandler     *agentHandler
	handshakeTimeout time.Timer
	userData         interface{}
	event            NetEvent
	uid              uint
	pushChan         chan *message.Message
	closeChan        chan struct{}
	locker           *sync.RWMutex // 推送请求反向锁，等待当前请求完成再推送
}

func NewAgent(conn tcp.Conn, event NetEvent) tcp.Agent {
	a := &agent{conn: conn}
	a.agentHandler = NewAgentHandler(a)
	a.event = event
	a.pushChan = make(chan *message.Message, 1000)
	a.closeChan = make(chan struct{})
	a.locker = &sync.RWMutex{}
	go func() {
		defer func() {
			if e := recover(); e != nil {
				log.Error(string(debug.Stack()), e)
			}
		}()

		for {
			select {
			case msg := <-a.pushChan:
				a.locker.RLock()
				a.locker.RUnlock()
				a.conn.Write(msg)

			case <-a.closeChan:
				return
			}
		}
	}()

	return a
}

// 只允许被调用一次
func (a *agent) LoopRead() {
	for {
		data, err := a.conn.Read()
		if err != nil {
			log.Info("OnClose UId: ", a.uid)
			break
		}
		pkgType, body := message.Decode(data)
		a.agentHandler.Handle(pkgType, body)
	}
}

// after 为true 时本消息会在当前的玩家请求后推送，否则立马推送
func (a *agent) Push(data interface{}, after ...bool) {
	route := message.GetMsg().GetPushMsg(data)
	if route == "" {
		return
	}
	msg := message.PoolGet()
	msg.Route = route
	msg.Type = message.Push
	d, err := jsoniter.Marshal(data)
	if err != nil {
		log.Error(err)
		return
	}
	msg.Data = d
	log.Info("发送:", msg)
	msg.Data = msg.Compress(msg.Data)

	err = msg.Encode()
	if err != nil {
		log.Error(err)
		return
	}
	if len(after) == 0 || !after[0] {
		a.conn.Write(msg)
	} else {
		select {
		case <-a.closeChan:
		default:
			select {
			case a.pushChan <- msg:
			default:
				a.Close()
			}
		}
	}
}

// 只允许被调用一次
func (a *agent) OnClose() {
	a.conn.Close()

	select {
	case <-a.closeChan:
	default:
		close(a.closeChan)
	}

	// 到此为止玩家的读写2个goroutine生命周期已终止
	// 创建新的goroutine去派发玩家掉线消息
	go func() {
		defer func() {
			if e := recover(); e != nil {
				log.Error(string(debug.Stack()), e)
			}
		}()
		a.event.OnClose(a)
	}()
}

func (a *agent) WriteMsg(msg *message.Message) {
	a.conn.Write(msg)
}

func (a *agent) Write(data *message.Message) error {
	return a.conn.Write(data)
}

func (a *agent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

// 主动关闭连接
func (a *agent) Close() {
	a.conn.Close()
}
func (a *agent) UserData() interface{} {
	return a.userData
}

func (a *agent) SetUserData(data interface{}) {
	a.userData = data
}

func (a *agent) UId() uint {
	return a.uid
}
func (a *agent) SetUid(uid uint) () {
	a.uid = uid
}
