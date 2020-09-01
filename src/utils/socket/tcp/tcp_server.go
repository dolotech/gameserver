package tcp

import (
	"gameserver/utils/log"
	"net"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

type TCPServer struct {
	Addr            string // 监听的端口
	MaxConnNum      int32  // 最大连接数限制
	PendingWriteNum int    // 写通道长度
	NewAgent        func(conn Conn) Agent
	ln              net.Listener
	conns           sync.Map       // 已经建立的连池
	wgLn            sync.WaitGroup // 用于关闭服务器时，等待正在建立的连接完成
	wgConns         sync.WaitGroup // 用于关闭服务器时，等待已经建立的连接关闭

	MinMsgLen    uint32
	MaxMsgLen    uint32
	LittleEndian bool
	online       int32
}

// 最线玩家数量
func (server *TCPServer) OnlineCount() int32 {
	return atomic.LoadInt32(&server.online)
}
func (server *TCPServer) Start() {
	server.init()
	server.run()
}

func (server *TCPServer) init() {
	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Error("%v", err)
	}

	if server.MaxConnNum <= 0 {
		server.MaxConnNum = 100
		log.Warning("invalid MaxConnNum, reset to %v", server.MaxConnNum)
	}
	if server.PendingWriteNum <= 0 {
		server.PendingWriteNum = 100
		log.Warning("invalid PendingWriteNum, reset to %v", server.PendingWriteNum)
	}
	if server.NewAgent == nil {
		log.Error("NewAgent must not be nil")
	}

	server.ln = ln
}

func (server *TCPServer) run() {
	server.wgLn.Add(1)
	defer server.wgLn.Done()

	var tempDelay time.Duration
	for {
		conn, err := server.ln.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Warning("accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return
		}
		tempDelay = 0
		if atomic.LoadInt32(&server.online) >= server.MaxConnNum {
			conn.Close()
			log.Warning("too many connections")
			continue
		}

		server.conns.Store(conn, struct{}{})
		atomic.AddInt32(&server.online, 1)
		server.wgConns.Add(1)

		msgParser := NewMsgReader(conn)
		msgParser.SetMsgLen(server.MinMsgLen, server.MaxMsgLen)
		msgParser.SetByteOrder(server.LittleEndian)
		tcpConn := newTCPConn(conn, msgParser, server.PendingWriteNum)
		agent := server.NewAgent(tcpConn)
		go func() {
			defer func() {
				if e := recover(); e != nil {
					log.Error(string(debug.Stack()),e)
				}
			}()
			agent.LoopRead()
			// cleanup
			tcpConn.Close()
			conn.Close()
			atomic.AddInt32(&server.online, -1)
			server.conns.Delete(conn)
			agent.OnClose()
			server.wgConns.Done()
		}()
	}
}

func (server *TCPServer) Close() {
	server.ln.Close()
	server.wgLn.Wait()

	// todo 没有关闭走tcpconn的正常关闭socket流程
	server.conns.Range(func(key, value interface{}) bool {
		key.(net.Conn).Close()
		return true
	})
	server.wgConns.Wait()
}
