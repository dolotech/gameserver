package message

const (
	TYPE_NULL          byte = 0
	TYPE_HANDSHAKE     byte = 1
	TYPE_HANDSHAKE_ACK byte = 2
	TYPE_HEARTBEAT     byte = 3
	TYPE_DATA          byte = 4
	TYPE_KICK          byte = 5

	PKG_HEAD_BYTES = 4
)

/**
 * Package protocol encode.
 *
 * Pomelo package format:
 * +------+-------------+------------------+
 * | type | body length |       body       |
 * +------+-------------+------------------+
 *
 * Head: 4bytes
 *   0: package type,
 *      1 - handshake,
 *      2 - handshake ack,
 *      3 - heartbeat,
 *      4 - data
 *      5 - kick
 *   1 - 3: big-endian body length
 * Body: body length bytes
 */

func (this *Message) WriteLen(len int) {
	this.WriteByte(byte(len >> 16))
	this.WriteByte(byte(len >> 8))
	this.WriteByte(byte(len))
}
func (this *Message) Pack(h byte, b []byte) {
	length := len(b)
	this.WriteByte(h)
	this.WriteByte(byte(length >> 16))
	this.WriteByte(byte(length >> 8))
	this.WriteByte(byte(length))
	this.Write(b)
}

func Decode(buffer []byte) (pkgType byte, body []byte) {
	pkgType = buffer[0]
	body = buffer[PKG_HEAD_BYTES:]
	return
}
