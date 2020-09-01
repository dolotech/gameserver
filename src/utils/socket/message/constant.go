package message

/**
 * ==========================
 *       消息類型 Flag
 * ==========================
 *
 * ------------------------------------------
 * |   type   |  flag  |       other        |
 * |----------|--------|--------------------|
 * | request  |----000-|<message id>|<route>|
 * | notify   |----001-|<route>             |
 * | response |----010-|<message id>        |
 * | push     |----011-|<route>             |
 * ------------------------------------------
 *
 */
const (
	Request  byte = 0x00
	Notify        = 0x01
	Response      = 0x02
	Push          = 0x03
)

const (
	msgRouteCompressMask = 0x01
	msgTypeMask          = 0x07
	msgRouteLengthMask   = 0xFF
	msgHeadLength        = 0x04			 // 消息头长度默认4个字节
)

var types = map[byte]string{
	Request:  "Request",
	Notify:   "Notify",
	Response: "Response",
	Push:     "Push",
}
