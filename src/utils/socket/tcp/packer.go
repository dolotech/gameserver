package tcp

import (
	"bufio"
	"encoding/binary"
	"gameserver/utils/log"
	"io"
	"math"
	"net"
	"time"
)

// --------------
// type | len | data |
// --------------
type packer struct {
	lenMsgLen    int    // 消息头长度默认4个字节
	maxMsgLen    uint32 // 消息体数据长度限制默认4096字节
	littleEndian bool   // 字节序
	reader *bufio.Reader
	conn net.Conn
}

func NewMsgReader(conn net.Conn) *packer {
	p := new(packer)
	p.lenMsgLen = 4    // 消息头长度默认4个字节
	p.maxMsgLen = 4096 // 接受最大包限制
	p.littleEndian = false
	p.reader = bufio.NewReaderSize(conn,4096)
	p.conn = conn
	return p
}

// It's dangerous to call the method on reading or writing
func (p *packer) SetMsgLen(minMsgLen uint32, maxMsgLen uint32) {
	if maxMsgLen != 0 {
		p.maxMsgLen = maxMsgLen
	}
	var max uint32
	switch p.lenMsgLen {
	case 1:
		max = math.MaxUint8
	case 2:
		max = math.MaxUint16
	case 4:
		max = math.MaxUint32
	}
	if p.maxMsgLen > max {
		p.maxMsgLen = max
	}
}

// It's dangerous to call the method on reading or writing
func (p *packer) SetByteOrder(littleEndian bool) {
	p.littleEndian = littleEndian
}

// 线程安全
func (p *packer) Read() ([]byte, error) {
	var b [4]byte
	bufHead := b[:p.lenMsgLen]
	p.conn.SetReadDeadline(time.Now().Add(time.Second * ReadTimeout))
	// 读取包头长度，包头长度的长度默认4个字节
	if _, err := io.ReadFull(p.reader, bufHead); err != nil {
		//log.Error("read head ", err)
		return nil, err
	}
	pkgType := int(bufHead[0])
	if pkgType < 1 || pkgType > 5 {
		return nil, ErrWrongMessageType
	}
	msgLen := binary.BigEndian.Uint32(bufHead)
	msgLen = msgLen & 0x00ffffff
	if msgLen <= 0 {
		return bufHead, nil
	}
	//log.Info("包头: %X包长%d 包类型%d ", bufHead, msgLen, int(bufHead[0]))
	// check len
	if msgLen > p.maxMsgLen {
		return nil, ErrMessageTooLong
	}
	msgData := make([]byte, msgLen+uint32(p.lenMsgLen))
	// 根据长度读指定长度的数据内容
	if n, err := io.ReadFull(p.reader, msgData[p.lenMsgLen:]); err != nil {
		log.Error("read content ", err, n)
		return nil, err
	}
	for i := 0; i < len(bufHead); i++ {
		msgData[i] = bufHead[i]
	}
	return msgData, nil
}
