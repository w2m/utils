package httputils

import (
	"fmt"
	"testing"
)

func Test_GetLocalAddr(t *testing.T) {
	ip := GetLocalAddr()
	fmt.Println(ip)
}

func Test_GetOutNetWorkIp(t *testing.T) {
	ip := GetOutNetWorkIp()
	fmt.Println(ip)
}

func Test_Get(t *testing.T) {
	data, _ := Get("http://www.baidu.com")
	fmt.Println(string(data))
}
