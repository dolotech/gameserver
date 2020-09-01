package utils

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)
func SingleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func GetKey(addr string) string {
	//addr := "hgm.rgstdz.com:443"
	/*if !strings.Contains(addr,"http://"){
		addr = "http://" + addr
	}*/
	c := http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(1 * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*1)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}

	reqest, err := http.NewRequest("GET", "/", nil)

	reqest.Host = "baidu.com"

	reqest.Proto = "HTTP/1.1"
	reqest.ProtoMajor = 1
	reqest.ProtoMinor = 1
	reqest.Close = false
	//reqest.RequestURI = "/"
	//ip := req.Host
	prefixPath := ""
	httpHost := "http://" + addr
	remote, err := url.Parse(httpHost)
	if err != nil {
		fmt.Print("Parse url", err)
	}

	reqest.URL.Host = addr
	reqest.URL.Scheme = remote.Scheme
	//设置路径
	reqest.URL.Path = prefixPath + SingleJoiningSlash(remote.Path, reqest.URL.Path)
	//设置参数
	reqest.PostForm = reqest.PostForm
	reqest.URL.RawQuery = reqest.URL.RawQuery
	reqest.Form = reqest.Form

	//reqest.Header.Set("Referer", "http://xgyhbg.bca9.cn:443")
	//reqest.Header.Set("Origin", "klj.4y8wn.cn:443")
	//reqest.RemoteAddr =

	if clientIP, _, err := net.SplitHostPort(reqest.RemoteAddr); err == nil {
		if prior, ok := reqest.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		reqest.Header.Set("X-Forwarded-For", clientIP)
	}

	//Upgrade-Insecure-Requests: 1
	//Host: hgm.rgstdz.com:443
	//增加header选项
	//reqest.Header.Add("Cookie", "xxxxxx")
	//reqest.Header.Add("User-Agent", "xxx")
	reqest.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,en-US;q=0.6,en;q=0.5;q=0.4")
	reqest.Header.Set("Accept-Encoding", "gzip, deflate")
	reqest.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	reqest.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.116 Safari/537.36 QBCore/4.0.1219.400 QQBrowser/9.0.2524.400 Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36 MicroMessenger/6.5.2.501 NetType/WIFI WindowsWechat")

	resp, err := c.Do(reqest)

	//log.Info("Header: ",reqest.Header)
	//log.Info("RequestURI:",reqest.RequestURI)

	if err != nil { // 提交异常,返回错误
		fmt.Printf("[error] short err is %v", err)
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		fmt.Print("bodyerr: ",err,string(body))
		//log.Info("body:",string(body))
		fmt.Print("StatusCode:", resp.StatusCode)
		fmt.Print("Location", resp.Header.Get("Location"))
		fmt.Print("Header", resp.Header)
		return resp.Header.Get("Location")
	}
	return ""
}
