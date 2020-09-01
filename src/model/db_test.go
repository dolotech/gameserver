package model

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

func TestAddJsonFormTag(t *testing.T) {
	rs := AddJsonFormGormTag(`
type Order struct { 
	POrderID  int   
	Platform  string
	ProductId int   
	PlayerId  int   
	Notified  bool  
	Resolved  bool  
}                   
	`)
	fmt.Println(rs)
}

func AddJsonFormGormTag(in string) string {
	var result string
	scanner := bufio.NewScanner(strings.NewReader(in))
	var oldLineTmp = ""
	var lineTmp = ""
	//var propertyTmp = ""
	var seperateArr []string
	for scanner.Scan() {
		oldLineTmp = scanner.Text()
		lineTmp = strings.Trim(scanner.Text(), " ")
		if strings.Contains(lineTmp, "{") || strings.Contains(lineTmp, "}") {
			result = result + oldLineTmp + "\n"
			continue
		}
		seperateArr = Split(lineTmp, " ")
		// 接口或者父类声明不参与tag, 自带tag不参与tag
		if len(seperateArr) == 1 || len(seperateArr) == 3 {
			continue
		}
		//propertyTmp = HumpToUnderLine(seperateArr[0])
		firstToLower := FirstToLower(seperateArr[0])
		//fmt.Println(firstToLower,"==",seperateArr[0])
		//oldLineTmp = oldLineTmp + fmt.Sprintf("    `gorm:\"column:%s\" json:\"%s\" form:\"%s\"`", propertyTmp, propertyTmp, propertyTmp)
		oldLineTmp = oldLineTmp + fmt.Sprintf("    `gorm:\"column:%s\" json:\"%s\" `", firstToLower, firstToLower)
		result = result + oldLineTmp + "\n"
	}
	return result
}

// 字符首字母小写写
func FirstToLower(s string) string {
	rs := strings.Trim(s, " ")
	rs = strings.Trim(rs, "\t")
	return strings.ToLower(rs[0:1]) + rs[1:]
}

// Capitalize 字符首字母大写
func Capitalize(str string) string {
	var upperStr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 { // 后文有介绍
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				fmt.Println("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

// 增强型split，对  a,,,,,,,b,,c     以","进行切割成[a,b,c]
func Split(s string, sub string) []string {
	var rs = make([]string, 0, 20)
	tmp := ""
	Split2(s, sub, &tmp, &rs)
	return rs
}

// 附属于Split，可独立使用
func Split2(s string, sub string, tmp *string, rs *[]string) {
	s = strings.Trim(s, sub)
	if !strings.Contains(s, sub) {
		*tmp = s
		*rs = append(*rs, *tmp)
		return
	}
	for i, _ := range s {
		if string(s[i]) == sub {
			*tmp = s[:i]
			*rs = append(*rs, *tmp)
			s = s[i+1:]
			Split2(s, sub, tmp, rs)
			return
		}
	}
}

// 驼峰转下划线
func HumpToUnderLine(s string) string {
	if s == "ID" {
		return "id"
	}else if s=="IP"{
		return "ip"
	}
	var rs string
	elements := FindUpperElement(s)
	for _, e := range elements {
		s = strings.Replace(s, e, "_"+strings.ToLower(e), -1)
	}
	rs = strings.Trim(s, " ")
	rs = strings.Trim(rs, "\t")
	return strings.Trim(rs, "_")
}

// 找到字符串中大写字母的列表,附属于HumpToUnderLine
func FindUpperElement(s string) []string {
	var rs = make([]string, 0, 10)
	for i := range s {
		if s[i] >= 65 && s[i] <= 90 {
			rs = append(rs, string(s[i]))
		}
	}
	return rs
}
