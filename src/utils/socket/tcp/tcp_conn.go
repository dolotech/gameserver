package tcp

import (
	"bufio"
	"gameserver/utils/log"
	"gameserver/utils/socket/message"
	"net"
	"runtime/debug"
	"time"
)

type Agent interface {
	LoopRead()
	OnClose()
}

type Conn interface {
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Write(*message.Message) error
	Read() ([]byte, error)
}

/*
	对已创建连接读写管理
*/
type tcpConn struct {
	conn      net.Conn
	writeChan chan *message.Message
	closeChan chan struct{}
	msgReader *packer // 消息读取处理
}

func newTCPConn(conn net.Conn, reader *packer, pendingWriteNum int) *tcpConn {
	tcpConn := new(tcpConn)
	tcpConn.conn = conn
	tcpConn.closeChan = make(chan struct{})
	tcpConn.writeChan = make(chan *message.Message, pendingWriteNum)
	tcpConn.msgReader = reader
	writer := bufio.NewWriterSize(conn, 4096)
	go func() {
		defer func() {
			if e := recover(); e != nil {
				log.Error(string(debug.Stack()), e)
			}
		}()
		for b := range tcpConn.writeChan {
			if b == nil {
				break
			}
			conn.SetWriteDeadline(time.Now().Add(time.Second * WriteTimeout))
			_, err := writer.Write(b.Bytes())
			message.PoolPut(b)
			if err != nil {
				log.Error(err)
				break
			}
			if len(tcpConn.writeChan) == 0 {
				writer.Flush()
			}
		}
		close(tcpConn.closeChan)
		tcpConn.conn.Close()
	}()
	return tcpConn
}

// 等待通道消息发送完成再关闭TPC连接,线程安全
func (tcpConn *tcpConn) Close() {
	select {
	case <-tcpConn.closeChan:
		break
	default:
		select {
		case tcpConn.writeChan <- nil:
			break
		default:
		}
	}
}
// 线程安全
func (tcpConn *tcpConn) Write(b *message.Message) error {
	if b == nil {
		return ErrWriteDataIsNil
	}
	select {
	case <-tcpConn.closeChan:
		break
	default:
		select {
		case tcpConn.writeChan <- b:
			break
		default:
			log.Error("writeChan block")
			tcpConn.Close()
			return ErrWriteChannelFull
		}
	}
	return nil
}
func (tcpConn *tcpConn) LocalAddr() net.Addr {
	return tcpConn.conn.LocalAddr()
}

func (tcpConn *tcpConn) RemoteAddr() net.Addr {
	return tcpConn.conn.RemoteAddr()
}
func (tcpConn *tcpConn) Read() ([]byte, error) {
	return tcpConn.msgReader.Read()
}
