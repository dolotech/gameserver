/**********************************************************
 * Author        : Michael
 * Email         : dolotech@163.com
 * Last modified : 2016-04-30 09:40
 * Filename      : utils.go
 * Description   : 常用的工具方法
 * *******************************************************/
package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// Auth
func GetAuth() []rune {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var list []rune
	for i := 0; i < 6; i++ {
		ran := r.Intn(122-97+1) + 97
		list = append(list, rune(ran))
	}
	return list
}

/**
 * 截取字符串
 * @param string str
 * @param begin int
 * @param length int
 * @return int 长度
 */
func SubStr(str string, begin, length int) (substr string) {
	// 将字符串的转换成[]rune
	rs := []rune(str)
	lth := len(rs)

	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}

	// 返回子串
	return string(rs[begin:end])
}

//整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}

//整形转换成字节
func Int64ToBytes(n int64) []byte {
	x := int64(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt64(b []byte) int64 {
	bytesBuffer := bytes.NewBuffer(b)
	var x int64
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int64(x)
}

//切片中字符串第一个位置
func SliceIndexOf(arr []string, str string) int {
	var index int = -1
	arrlen := len(arr)
	for i := 0; i < arrlen; i++ {
		if arr[i] == str {
			index = i
			break
		}
	}
	return index
}

//字节转换成整形
func SliceLastIndexOf(arr []string, str string) int {
	var index int = -1
	for arrlen := len(arr) - 1; arrlen > -1; arrlen-- {
		if arr[arrlen] == str {
			index = arrlen
			break
		}
	}
	return index
}

//字节转换成整形
func SliceRemoveFormSlice(oriArr []string, removeArr []string) []string {
	endArr := oriArr[:]
	for _, value := range removeArr {
		index := SliceIndexOf(endArr, value)
		if index != -1 {
			endArr = append(endArr[:index], endArr[index+1:]...)
		}
	}
	return endArr
}

// 把时间戳转换成头像存储目录
func TimeToHeadphpoto(t int64, userid int, headname int64) (string, string) {
	var str, name string
	ti := time.Unix(t, 0)
	str = ti.Format("2006/01/02/15")

	str = "./headpic/" + str + "/" + strconv.Itoa(userid)
	if headname == 0 {
		name = "/130_" + strconv.Itoa(userid) + ".jpg"
	} else {
		name = "/" + strconv.Itoa(int(headname)) + ".jpg"
	}
	return str, name
}

// 把时间戳转换成头像存储目录
func TimeToPhpotoPath(t int64, userid int) string {
	var str string
	ti := time.Unix(t, 0)
	str = ti.Format("2006/01/02/15")
	return "./photo/" + str + "/" + strconv.Itoa(userid)
}
func UseridCovToInvate(userid string) uint32 {
	useridbyte := []byte(userid)
	useridbyte = useridbyte[len(useridbyte)-4:]
	timestr := []byte(strconv.Itoa(int(time.Now().Unix())))
	timestr = timestr[len(timestr)-5:]
	useridbyte = append(useridbyte, timestr...)
	code, _ := strconv.Atoi(string(useridbyte))
	return uint32(code)
}

var base64String = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
var flipbase = flip(base64String)
var baselen = len(base64String)

func Base62encode(num uint64) string {
	baseStr := ""
	for {
		if num <= 0 {
			break
		}

		i := num % uint64(baselen)
		baseStr += base64String[i]
		num = (num - i) / uint64(baselen)
	}
	return baseStr
}

func Base62decode(base62 string) uint64 {
	var rs uint64 = 0
	len := uint64(len(base62))
	var i uint64
	for i = 0; i < len; i++ {
		rs += flipbase[string(base62[i])] * uint64(math.Pow(float64(baselen), float64(i)))
	}
	return rs
}

func flip(s []string) map[string]uint64 {
	f := make(map[string]uint64)
	for index, value := range s {
		f[value] = uint64(index)
	}
	return f
}

// 用gob进行数据编码
func Encode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// 用gob进行数据解码
//
func Decode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}

// 对象深度拷贝
func Clone(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

//----------------------一下几个函数只对数字字符串有效----------------------------
func IsNumString(str string) bool {
	runeArr := []rune(str)
	for i := 0; i < len(runeArr); i++ {
		if runeArr[i] > 56 {
			return false
		} else if runeArr[i] < 48 {
			return false
		}
	}
	return true
}
func Between(startid, endid string) []string {
	if startid == endid {
		return []string{startid}
	}
	ids := []string{}
	if len(startid) > len(endid) {
		return ids
	}
	start := []rune(startid)
	end := []rune(endid)

	for i := 0; i < len(start); i++ {
		if int(end[i]) < int(start[i]) {
			return ids
		} else if int(end[i]) > int(start[i]) {
			break
		}
	}
	for {
		ids = append(ids, startid)
		startid = StringAdd(startid)
		if startid == endid {
			ids = append(ids, startid)
			break
		}
	}
	return ids
}

func StringAddNum(numStr string, num int) string {
	for i := 0; i < num; i++ {
		numStr = StringAdd(numStr)
	}
	return numStr
}

// 字符串加法
func StringAdd(numStr string) string {
	runeArr := []rune(numStr)
	length := len(runeArr)
	add := true
	for i := length - 1; i >= 0; i-- {
		if runeArr[i] < 57 {
			runeArr[i]++
			add = false
			break
		} else {
			runeArr[i] = 48
		}
	}
	if add {
		runeArr = append([]rune{49}, runeArr...)
	}
	return string(runeArr)
}

// md5 加密
// func Md5(text string) string {
// 	hashMd5 := md5.New()
// 	io.WriteString(hashMd5, text)
// 	return fmt.Sprintf("%x", hashMd5.Sum(nil))
// }
func Md5(text string) string {
	h := md5.New()
	h.Write([]byte(text))                 //
	return hex.EncodeToString(h.Sum(nil)) // 输出加密结果
}

// 延迟second
func Sleep(second int) {
	<-time.After(time.Duration(second) * time.Second)
}

// 延迟1~second
func SleepRand(second int) {
	<-time.After(time.Duration(rand.Intn(second)+1) * time.Second)
}

// 延迟second
func Sleep64(second int64) {
	<-time.After(time.Duration(second) * time.Second)
}

// 延迟1~second
func SleepRand64(second int64) {
	<-time.After(time.Duration(rand.Int63n(second)+1) * time.Second)
}

// 用于调式显示掩码
func BitOr(v int64) {
	var s string
	for i := 1; i <= 64; i++ {
		if v&(1<<uint(i)) > 0 {
			s = fmt.Sprintf("%s %d", s, i)
		}
	}
}

func Byte2uint32(in []byte) []uint32 {
	out := make([]uint32, 0, len(in))
	for _, v := range in {
		out = append(out, uint32(v))
	}
	return out
}
func Byte2int32(in []byte) []int32 {
	out := make([]int32, 0, len(in))
	for _, v := range in {
		out = append(out, int32(v))
	}
	return out
}

func Int642uint32(in []int64) []uint32 {
	out := make([]uint32, 0, len(in))
	for _, v := range in {
		out = append(out, uint32(v))
	}
	return out
}

func String2uint32(in []string) []uint32 {
	out := make([]uint32, 0, len(in))
	for _, v := range in {
		t, _ := strconv.Atoi(v)
		out = append(out, uint32(t))
	}
	return out
}

func String2int(in []string) []int {
	out := make([]int, 0, len(in))
	for _, v := range in {
		t, _ := strconv.Atoi(v)
		out = append(out, t)
	}
	return out
}

func Uint322string(in []uint32) []string {
	out := make([]string, 0, len(in))
	for _, v := range in {
		out = append(out, strconv.Itoa(int(v)))
	}
	return out
}

func int2string(in []int) []string {
	out := make([]string, 0, len(in))
	for _, v := range in {
		out = append(out, strconv.Itoa(v))
	}
	return out
}

// 是否在slice里面
func InSlice(ms uint32, arr []uint32) bool {
	for _, v := range arr {
		if ms == v {
			return true
		}
	}
	return false
}

func Truncate6Words(origin string) string {
	newString := origin
	nameRune := []rune(origin)
	if len(nameRune) > 6 {
		newString = string(nameRune[:6])
	}
	return newString
}

// "1.1.1"格式版本号对比,origin  =  target  :1;origin <   target  :-1;   origin =  target:0
func VersionContrast(origin, target string) (int, error) {

	originArr := strings.Split(origin, ".")
	targetArr := strings.Split(target, ".")

	for i := 0; i < len(originArr); i++ {
		originItem := originArr[i]
		originInt, err := strconv.Atoi(originItem)
		if err != nil {
			return 0, err
		}
		if len(targetArr) <= i {
			return 1, nil
		}
		targetItem := targetArr[i]

		targetInt, err := strconv.Atoi(targetItem)
		if err != nil {
			return 0, err
		}

		if originInt > targetInt {
			return 1, nil
		} else if originInt < targetInt {
			return -1, nil
		}
	}
	return 0, nil
}

func LogPrefix(uid uint32, str string) string {
	return fmt.Sprintf("玩家:%v 操作:[%v]", uid, str)
}

/*serverid helper
serverid采用6位数字，例如：123001
第1，2位是appid,对应客户端应用id
第3，4位是game类型,约局10，金币30
第5，6位是可执行文件编号,比如123001,123002分别是牛牛金币服两个game
*/

func ToServerType(sid int) int {
	return ((sid / 1000) % 100) % 10
}
