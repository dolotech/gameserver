package message

// ------------------------------------------
// |   type   |  flag  |       other        |
// |----------|--------|--------------------|
// | request  |----000-|<message id>|<route>|
// | notify   |----001-|<route>             |
// | response |----010-|<message id>        |
// | push     |----011-|<route>             |
// ------------------------------------------
func (m *Message) Encode() error {
	if m.Type < Request || m.Type > Push { //   消息类型是否合法
		return ErrWrongMessageType
	}
	flag := m.Type << 1
	m.WriteByte(TYPE_DATA)
	length := len(m.Data) + 1
	var msgid []byte
	if m.Type == Request || m.Type == Response {
		n := m.ID
		// variant length encode
		for {
			b := byte(n % 128)
			n >>= 7
			if n != 0 {
				msgid = append(msgid, b+128)
				length += 1
			} else {
				msgid = append(msgid, b)
				length += 1
				break
			}
		}
	}
	if m.Type == Request || m.Type == Notify || m.Type == Push {
		length += (len(m.Route) + 1)
	}
	m.WriteLen(length)
	m.WriteByte(flag)
	if len(msgid) > 0{
		m.Write(msgid)
	}
	if m.Type == Request || m.Type == Notify || m.Type == Push {
		m.WriteByte(byte(len(m.Route)))
		m.Write([]byte(m.Route))

	}
	m.Write(m.Data)
	return nil
}

// Decode unmarshal the bytes slice to a message
func (m *Message) Decode(data []byte) error {
	if len(data) < msgHeadLength {
		return ErrInvalidMessage
	}
	//m := New()
	flag := data[0]
	offset := 1
	m.Type = byte((flag >> 1) & msgTypeMask)

	if m.Type < Request || m.Type > Push { //   消息类型是否合法
		return ErrWrongMessageType
	}

	if m.Type == Request || m.Type == Response {
		id := uint(0)
		for i := offset; i < len(data); i++ {
			b := data[i]
			id += uint(b&0x7F) << uint(7*(i-offset))
			if b < 128 {
				offset = i + 1
				break
			}
		}
		m.ID = id
	}

	if m.Type == Request || m.Type == Notify || m.Type == Push {
		rl := data[offset]
		offset++
		m.Route = string(data[offset:(offset + int(rl))])
		offset += int(rl)
	}

	m.Data = data[offset:]
	return nil
}
