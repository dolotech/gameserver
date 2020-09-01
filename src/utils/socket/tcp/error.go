package tcp

import "errors"

var (
	ErrWrongMessageType  = errors.New("wrong message type")
	ErrMessageTooShort   = errors.New("message too short")
	ErrMessageTooLong   = errors.New("message too long")
	ErrWriteDataIsNil   = errors.New("write data is nil")
	ErrWriteChannelFull   = errors.New("write channel full")
	ErrWriteChannelClosed  = errors.New("write channel closed")
)

const (
	WriteTimeout = 2// 写网络包超时(单位/秒)
	ReadTimeout = 9// 读网络包超时(单位/秒)
)
