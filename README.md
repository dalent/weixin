# weixin
weixin分享api
golang版本的微信的jssdk
```
package main
import (
  "fmt"
  "github.com/dalent/weixin"
)
func main(){
  //指定key，secret，开启自动刷新
  weixin.Init("appKey","appSecret",true)
  resp,_:= weixin.GetSignPackage("url")
  fmt.Println(resp)
}
```

安装
```
go get github.com/dalent/weixin
```
