package weixin

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var (
	chars              = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	jsApiPattern       = `https://api.weixin.qq.com/cgi-bin/ticket/getticket?type=jsapi&access_token=%s`
	AccessTokenPattern = `https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s`
	defaultRand        = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type WeiXinToken struct {
	AccessToken string `json:"access_token"`
	ExpireIn    int    `json:"expires_in"`
}

type WeiXinJs struct {
	Ticket   string `json:"access_token"`
	ExpireIn int    `json:"expires_in"`
}

type WeiXinAccess struct {
	access WeiXinToken
	ticket WeiXinJs
	//更新信息需要锁
	mutex sync.Mutex
}

//这个结构基本上全局唯一不变的,不需要返回多个副本
var weiXinAccess WeiXinAccess

func GetWeiXinAccess() *WeiXinAccess {
	return &weiXinAccess
}

func (p *WeiXinAccess) createNonceStr(length int) string {
	var str string
	for i := 0; i < length; i++ {
		tmpI := defaultRand.Intn(len(chars) - 1)
		str += chars[tmpI : tmpI+1]
	}
	return str
}

func (p *WeiXinAccess) getJsApiTicket() error {
	response := struct {
		Code      int    `json:"errcode"`
		Msg       string `json:"errmsg"`
		Ticket    string `json:"ticket"`
		ExpiresIn int    `json:"expires_in"`
	}{}
	resp, err := http.Get(fmt.Sprintf(jsApiPattern, p.access.AccessToken))
	if err != nil {
		return err
	}

	json.NewDecoder(resp.Body).Decode(&response)
	if response.Code != 0 {
		return errors.New(response.Msg)
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.ticket.ExpireIn = response.ExpiresIn
	p.ticket.Ticket = response.Ticket
	return nil
}

func (p *WeiXinAccess) getAccessToken() error {
	resp, err := http.Get(fmt.Sprintf(AccessTokenPattern, AppKey, AppSecret))
	if err != nil {
		return err
	}

	var tmpAccess WeiXinToken
	err = json.NewDecoder(resp.Body).Decode(&tmpAccess)
	if err != nil {
		return err
	}

	if tmpAccess.AccessToken == "" {
		return errors.New("access token get failed")
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.access = tmpAccess
	return nil
}

//下面的这一部分负责刷新token
func Min(a, b int) int {
	if a > b {
		return b
	}

	return a
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

func init() {
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
