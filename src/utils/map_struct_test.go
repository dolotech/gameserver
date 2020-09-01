package utils

import (
	"testing"
	"time"
)

type User struct {
	Name string
	Age  int8
	Date time.Time
	Logintime time.Time
}

type UserTest struct {
	Id        int64
	Username  string
	Password  string
	Name      string
	Logintime time.Time
}

func Test_StructAtoB(t *testing.T) {
	a := &User{Name: "michael"}
	b := &UserTest{Name:"Michael0000",Logintime:time.Now()}

	StructAtoB(a,b)
	b.Name = "123"
	t.Error(a)
}

func Test_Struct2Map(t *testing.T) {
	user := UserTest{Id: 5, Username: "zhangsan", Password: "pwd", Logintime: time.Now()}
	data := Struct2Map(user)
	t.Error(data)
}

func Test_Map2Struct(t *testing.T) {
	data := make(map[string]interface{})
	data["Name"] = "张三"
	data["Age"] = 26
	data["Date"] = "2015-09-29 00:00:00"

	result := &User{}

	_ = Map2Struct(data, result)

	t.Error(result)
}
