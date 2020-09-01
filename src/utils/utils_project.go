package utils

import (
	"io"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func GetOne(array [][]string, host string) string {
	for i := 0; i < len(array); i++ {
		if array[i][0] == host {
			return array[i][1]
		}
	}
	return ""
}

func RandomGetOne(array []string) string {
	return array[rand.Intn(len(array))]
}

type Fileutils struct {
	f *os.File
	p string
}

func NewFileUtils(filename string) *Fileutils {
	return &Fileutils{p: filename}
}

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func (this *Fileutils) Close() error {
	return this.f.Close()
}
func (this *Fileutils) Write(filedata string) error {
	if CheckFileIsExist(this.p) { //如果文件存在
		var err error
		this.f, err = os.OpenFile(this.p, os.O_APPEND|os.O_WRONLY, os.ModeAppend) //打开文件
		if err != nil {
			return err
		}
	} else {
		var err error
		this.f, err = os.Create(this.p) //创建文件
		if err != nil {
			return err
		}
	}
	_, err := io.WriteString(this.f, filedata) //写入文件(字符串)
	if err != nil {
		return err
	}
	return nil
}

//表情解码
func UnicodeEmojiDecode(s string) string {
	//emoji表情的数据表达式
	re := regexp.MustCompile("\\[[\\\\u0-9a-zA-Z]+\\]")
	//提取emoji数据表达式
	reg := regexp.MustCompile("\\[\\\\u|]")
	src := re.FindAllString(s, -1)
	for i := 0; i < len(src); i++ {
		e := reg.ReplaceAllString(src[i], "")
		p, err := strconv.ParseInt(e, 16, 32)
		if err == nil {
			s = strings.Replace(s, src[i], string(rune(p)), -1)
		}
	}
	return s
}

//表情转换
func UnicodeEmojiCode(s string) string {
	ret := ""
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if len(string(rs[i])) == 4 {
			u := `[\u` + strconv.FormatInt(int64(rs[i]), 16) + `]`
			ret += u

		} else {
			ret += string(rs[i])
		}
	}
	return ret
}
