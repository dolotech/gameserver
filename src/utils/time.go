package utils

import (
	"strconv"
	"time"
)

//-----------------------------------------------------------------
const FORMAT string = "2006-01-02 15:04:05"
const FORMATDATA string = "2006-01-02 "

// 获取当前时间截
func TimestampNano() int64 {
	return time.Now().UnixNano()
}

// 获取当前时间截
func Timestamp() int64 {
	return time.Now().Unix()
}

// 获取本周六零点时间截
func TimestampSaturday() int64 {
	now := time.Now()
	unix := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Unix()
	return unix + int64(time.Saturday-now.Weekday())*86400
}

// 获取本地当天零点时间截
func TimestampToday() int64 {
	return time.Date(Year(), Month(), Day(), 0, 0, 0, 0, time.Local).Unix()
}
func TimestampTodayStr() string {
	t := time.Date(Year(), Month(), Day(), 0, 0, 0, 0, time.Local).Unix()
	return strconv.FormatInt(t, 10)
}

// 获取本地昨天零点时间截
func TimestampYesterday() int64 {
	return TimestampToday() - 86400
}

// 获取本地明天零点时间截
func TimestampTomorrow() int64 {
	return TimestampToday() + 86400
}

// 获取当前年
func Year() int {
	return time.Now().Year()
}

// 获取当前月
func Month() time.Month {
	return time.Now().Month()
}

// 获取当前天
func Day() int {
	return time.Now().Day()
}

// 获取当前周
func Weekday() time.Weekday {
	return time.Now().Weekday()
}

// 获取指定时间截的年
func Unix2Year(t int64) int {
	return time.Unix(t, 0).Year()
}

// 获取指定时间截的月
func Unix2Month(t int64) time.Month {
	return time.Unix(t, 0).Month()
}

// 两个时间戳是否是同一天
func SameDay(t1, t2 int64) bool {
	t11 := time.Unix(t1, 0)
	t22 := time.Unix(t2, 0)
	return t11.Year() == t22.Year() &&
		t11.Month() == t22.Month() &&
		t11.Day() == t22.Day()
}

// 两个时间戳是否是一个月
func SameMonth(t1, t2 int64) bool {
	t11 := time.Unix(t1, 0)
	t22 := time.Unix(t2, 0)
	return t11.Month() == t22.Month()
}

// 获取指定时间截的天
func Unix2Day(t int64) int {
	return time.Unix(t, 0).Day()
}

// 时间戳转str格式化时间
func Unix2Str(t int64) string {
	return time.Unix(t, 0).Format(FORMAT)
}

// str格式当前日期
func DateStr() string {
	return time.Now().Format(FORMATDATA)
}

// str格式化时间转时间戳
func Str2Unix(t string) (int64, error) {
	the_time, err := time.Parse(FORMAT, t)
	if err == nil {
		return the_time.Unix(), err
	}
	return 0, err
}

// 获取指定年月的天数
func MonthDays(year int, month int) (days int) {
	if month != 2 {
		if month == 4 || month == 6 || month == 9 || month == 11 {
			days = 30
		} else {
			days = 31
		}
	} else {
		if ((year%4) == 0 && (year%100) != 0) || (year%400) == 0 {
			days = 29
		} else {
			days = 28
		}
	}
	return
}
