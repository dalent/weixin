package weixin

import (
	"fmt"
	"testing"
)

func TestWeixin(_ *testing.T) {
	resp, err := GetSignPackage("http://12.23.45")
	fmt.Println(resp, err)
}
