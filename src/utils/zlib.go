package utils

import (
	"bytes"
	"compress/zlib"
	"gameserver/utils/log"
	"io"
)

//进行zlib压缩
func DoZlibCompress(inb []byte) []byte {
	if len(inb)> 0{
		var in bytes.Buffer
		w := zlib.NewWriter(&in)
		w.Write(inb)
		w.Close()
		return in.Bytes()
	}
	return  inb
}

//进行zlib解压缩
func DoZlibUnCompress(compressSrc []byte) ([]byte, error) {
	if len(compressSrc)>0{
		b := bytes.NewReader(compressSrc)
		r, _ := zlib.NewReader(b)
		log.Info(len(compressSrc))
		var out bytes.Buffer
		io.Copy(&out, r)
		r.Close()
		return out.Bytes(), nil
	}
	return compressSrc,nil
}
