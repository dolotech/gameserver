package utils

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// 验证是否邮箱
func EmailRegexp(mail string) bool {
	b := false
	if mail != "" {
		reg := regexp.MustCompile(`^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(\.[a-zA-Z0-9_-]+)$`)
		b = reg.FindString(mail) != ""
	}
	return b
}

// 验证是否手机
func PhoneRegexp(phone string) bool {
	b := false
	if phone != "" {
		reg := regexp.MustCompile(`^(86)*0*1\d{10}$`)
		b = reg.FindString(phone) != ""
	}
	return b
}

// 验证账号是否合法
func AccountRegexp(account string) bool {
	b := false
	if account != "" {
		reg := regexp.MustCompile(`^[a-zA-Z0-9]{6,8}$`)
		b = reg.FindString(account) != ""
	}
	return b
}

// 验证只能由数字字母下划线组成的5-17位密码字符串
func AalidataPwd(name string) (b bool) {
	if name != "" {
		//reg := regexp.MustCompile(`^[a-zA-Z0-9_]*$`)
		reg := regexp.MustCompile(`^[a-zA-Z_]\w{5,17}$`)
		b = reg.FindString(name) != ""
	}
	return
}

// 不可见字符,用于用户提交的字符过滤分别对应为：,\0   \t  _  space  "  ` ctrl+z \n \r  `  %   \  ,
var IllegalNameRune = [13]rune{0x00, 0x09, 0x5f, 0x20, 0x22, 0x60, 0x1a, 0x0a, 0x0d, 0x27, 0x25, 0x5c, 0x2c}

var hasIllegalNameRune = func(c rune) bool {
	for _, v := range IllegalNameRune {
		if v == c {
			return true
		}
	}
	return false
}

// 限制最大字符数，检测不可见字符
// maxcount 限制的最大字符数，1个中文=2个英文
func LegalName(name string, maxcount int) bool {
	if !utf8.ValidString(name) {
		return false
	}

	num := len([]rune(name)) + len([]byte(name))
	result := float64(num) / 4.0
	sum := int(result + 0.99)

	if sum > maxcount*2 {
		return false
	}
	return strings.IndexFunc(name, hasIllegalNameRune) == -1
}


