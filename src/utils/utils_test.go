package utils

import (
	"gameserver/utils/log"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func Test_ntime(t *testing.T) {
	t.Error(Unix2Day(time.Now().Unix() +3600*10))



	t.Error(SameDay(time.Now().Unix(),time.Now().Unix()+3600*5))
}
func Test_new_id(t *testing.T) {
	id:=NewID()
	for i:=1;i<100000000;i++{

		t.Error(id.Id())
	}
}
func Test_Http_head(t *testing.T) {
	//h:=`User-Agent: Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.116 Safari/537.36 QBCore/4.0.1278.400 QQBrowser/9.0.2524.400 Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2875.116 Safari/537.36 NetType/WIFI MicroMessenger/7.0.5 WindowsWechat`
	//h := `（HTTP_USER_AGENT）：Mozilla/5.0 (iPhone; CPU iPhone OS 12_1_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/16D57 MicroMessenger/7.0.4(0x17000428) NetType/4G Language/zh_CN`

	h:=`（HTTP_USER_AGENT）：Mozilla/5.0 (Linux; Android 8.0.0; HUAWEI NXT-AL10 Build/HUAWEINXT-AL10; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/66.0.3359.126 MQQBrowser/6.2 TBS/044904 Mobile Safari/537.36 MMWEBID/2411 MicroMessenger/7.0.6.1500(0x2700063E) Process/tools NetType/4G Language/zh_CN`
	if strings.Contains(h, "iPhone") {
		log.Error("iPhone")
	}

	if strings.Contains(h, "Windows") {
		log.Error("Windows")
	}

	if strings.Contains(h, "Android") {
		log.Error("Android")
	}
	if strings.Contains(h, "4G") {
		log.Error("4G")
	}

	if strings.Contains(h, "3G") {
		log.Error("3G")
	}


	if strings.Contains(h, "WIFI") {
		log.Error("WIFI")
	}

}
func Test_group(t *testing.T) {
	var wait sync.WaitGroup

	go func() {
		log.Info(1)
	}()
	t.Error(2)
	wait.Wait()
}
func Test_11(t *testing.T) {
	t.Error(len(`{"type":6}`))
	t.Error(len(`{"socketio":"json","version":1}`))
}
func Test_1(t *testing.T) {
	payload := `{"type":3,"invocationId":"0","result":{"result":{"scrollmsg":"欢迎进入爱玩，请各位玩家认准自己的ID 防伪签名跟等级变化。谨防上当受骗，感谢支持爱玩！","fzbHint":"诚邀广大玩家与我们一起抵制盗版!\n请牢记自己的个性宣言和等级称号,如遇到宣言无法显示或与自己设置的宣言不同,便是盗版平台,请你我联手抵制盗版平台对大家进行外挂等各种侵害行为，爱玩提示。","fzbHint2":"致广大玩家：\n        打开游戏后，如果发现游戏界面和之前玩的界面不一样,就肯定是山寨平台及外挂平台,请认准游戏界面，防伪签名，等级和ID，谨防被骗，爱玩提示。","title":"葫芦鱼","appName":"新爱玩乐厅","freeGames":"","disabledGames":"","club2Enable":false,"userInfo":{"id":227033,"display_id":434172,"nickname":"李是谁","headimgurl":"http://thirdwx.qlogo.cn/mmopen/vi_32/ibe3FwU4ZxAWCIgk6u0xAiajCHSa7Fj5fCdJicTUXBxB0LiayPblOMj0KofTkbVfoFQq4ROJQ8HBk9TcbGjhAgngzQ/132","group_enable":true,"enable":true,"sign":"发个火呀","sign2":null,"card":20,"exp":28,"level":4,"level_exp":64,"phone":null},"sx2Price":"{\"12\": {\"6\": 4,\"10\": 5,\"13\": 6,\"16\": 8,\"18\": 9,\"20\": 10,\"26\": 14},\"18\": {\"6\": 6,\"10\": 7,\"13\": 9,\"16\": 12,\"18\": 14,\"20\": 16,\"26\": 20}}","sb2Price":"{\"12\": {\"6\": 4,\"10\": 5,\"13\": 6,\"16\": 8,\"18\": 9,\"20\": 10,\"26\": 14},\"18\": {\"6\": 6,\"10\": 7,\"13\": 9,\"16\": 12,\"18\": 14,\"20\": 16,\"26\": 20}}","fzPrice":"{\"12\": {\"6\": 4,\"10\": 5,\"13\": 6,\"16\": 8,\"18\": 9,\"20\": 10,\"26\": 14},\"18\": {\"6\": 6,\"10\": 7,\"13\": 9,\"16\": 12,\"18\": 14,\"20\": 16,\"26\": 20}}","navUrl":"http://t.cn/AiYDSyuI"},"errcode":0,"errmsg":null}}`
	t.Error(strings.Contains(payload, `,"userInfo":{"id":`))

}
func TestRandomGetOne(t *testing.T) {
	array := []string{"1", "2", "3", "4"}
	t.Error(RandomGetOne(array))
}

func TestAtomic(t *testing.T) {
	var gamestart int32

	t.Error(atomic.CompareAndSwapInt32(&gamestart, 0, 1))
	t.Error(atomic.CompareAndSwapInt32(&gamestart, 1, 0))
	t.Error(atomic.CompareAndSwapInt32(&gamestart, 0, 1))

	t.Error(gamestart)
}
func TestTicker(t *testing.T) {

	fangBao := make(chan struct{}, 1)

	go func() {
		select {

		case d, ok := <-fangBao:
			t.Error("close", d, ok)
		}
		t.Error("ticker")
	}()

	//ticker.Stop()

	//close(fangBao)

	fangBao <- struct{}{}

	time.Sleep(time.Second * 4)
}
func TestFun(t *testing.T) {
	var handlePong func(string) error
	t.Log(handlePong("234"))
}
func TestChan(t *testing.T) {
	var ch chan struct{}

	go func() {
		select {
		case <-ch:
			t.Error("1")
		}
	}()

	go func() {
		select {
		case <-ch:
			t.Error("2")
		}
	}()

	ch <- struct{}{}

}
func TestUserInfo3(t *testing.T) {
	/*db.InitMysql("root", "123456", "127.0.0.1:3306", "aiwan_proxy")
	user:= data.UserInfoTb{
		Id:1,
		Display_id:1,
		Nickname:"🙃",
	}


	t.Error(user.Save())*/
	Nickname := "🙃"
	//Nickname="12324"

	Nickname = UnicodeEmojiCode(Nickname)
	Nickname = UnicodeEmojiDecode(Nickname)

	t.Error(Nickname)
}

func BenchmarkRandomGetOne(b *testing.B) {
	b.ReportAllocs()
	array := []string{"1", "2", "3", "4"}
	for i := 0; i < b.N; i++ {
		RandomGetOne(array)
	}
}

func BenchmarkPoint(b *testing.B) {
	b.ReportAllocs()
	j := 0
	for i := 0; i < b.N; i++ {
		modify(&j)
	}

	b.Error(j)
}
func modify(i *int) {
	*i = 1000
}



func TestEmailRegexp(t *testing.T) {
	type args struct {
		mail string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.

	}
	for _, tt := range tests {
		if got := EmailRegexp(tt.args.mail); got != tt.want {
			t.Errorf("%q. EmailRegexp() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
