package weixin

import (
	"crypto/sha1"
	"fmt"
	"time"
)

//appkey appSecret 需要自己指定
var (
	appKey    = "123456"
	appSecret = "123456"
)

type SignPackage struct {
	AppId     string `json:"appId"`
	NonceStr  string `json:"nonceStr"`
	Timetamp  int64  `json:"timestamp"`
	URL       string `json:"url"`
	Signature string `json:"signature"`
	RawString string `json:"rawString"`
}

//不能是固定的url，url后缀可能有很多参数例如 http://a.b.c?a=1&b=2
func GetSignPackage(url string) (*SignPackage, error) {
	timestamp := time.Now().Unix()
	weixin := GetWeiXinAccess()
	noncestr := weixin.createNonceStr(16)
	str := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s",
		weixin.ticket.Ticket, noncestr, timestamp, url)

	signature := fmt.Sprintf("%x", sha1.Sum([]byte(str)))

	return &SignPackage{appKey, noncestr, timestamp, url, signature[:], str}, nil
}

func Init(key, secret string) {
	appKey = key
	appSecret = secret

	//先获得一次,获得失败panic
	if err := Refresh(); err != nil {
		panic(err)
	}

	go func() {
		for {
			err := Refresh()

			if err != nil {
				continue
			}
			//time.Sleep(10 * time.Second)
			//两个最少的一半时间刷新应该是够的
			time.Sleep(time.Second * time.Duration(Min(weiXinAccess.access.ExpireIn, weiXinAccess.ticket.ExpireIn)/2))
		}
	}()
}
func Refresh() error {
	if err := weiXinAccess.getAccessToken(); err != nil {
		return err
	}

	if err := weiXinAccess.getJsApiTicket(); err != nil {
		return err
	}

	return nil
}
