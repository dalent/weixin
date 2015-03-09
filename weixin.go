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

func Init(key, secret string, refresh bool) {
	appKey = key
	appSecret = secret
	GetWeiXinAccess().Init(refresh)
}
