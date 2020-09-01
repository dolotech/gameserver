package message

import (
	"bytes"
	"compress/zlib"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
)

// Message represents a unmarshaled message or a message which to be marshaled
type Message struct {
	Type  byte   // message type
	ID    uint   // unique id, zero while notify mode
	Route string // route for locating service
	Data  []byte // payload
	bytes.Buffer

	zlib zliber
}

func (m *Message) Unmarshal(i interface{}) (err error) {
	m.Data = m.UnCompress(m.Data)
	return jsoniter.Unmarshal(m.Data, i)
}
func (m *Message) String() string {
	return fmt.Sprintf("Type: %s, ID: %d, Route: %s,  BodyLength: %d  ,Body:%s",
		types[m.Type],
		m.ID,
		m.Route,
		len(m.Data),
		string(m.Data))
}

type zliber struct {
	in     bytes.Buffer
	out    bytes.Buffer
	reader bytes.Reader
}

//进行zlib解压缩
func (m *Message) UnCompress(compressSrc []byte) []byte {
	return m.zlib.unCompress(compressSrc)
}

//进行zlib压缩
func (m *Message) Compress(inb []byte) []byte {
	return m.zlib.compress(inb)
}

func (m *zliber) compress(inb []byte) []byte {
	if len(inb) > 0 {
		m.in.Reset()
		w := zlib.NewWriter(&m.in)
		w.Write(inb)
		w.Close()
		return m.in.Bytes()
	}
	return inb
}

func (m *zliber) unCompress(compressSrc []byte) []byte {
	if len(compressSrc) > 0 {
		m.out.Reset()
		m.reader.Reset(compressSrc)
		r, _ := zlib.NewReader(&m.reader)
		io.Copy(&m.out, r)
		r.Close()
		return m.out.Bytes()
	}
	return compressSrc
}
