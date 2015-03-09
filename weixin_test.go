package weixin

import (
	"fmt"
	"testing"
)

func TestWeixin(_ *testing.T) {
	Init("", "", true)
	resp, err := GetSignPackage("http://12.23.45")
	fmt.Println(resp, err)
}
