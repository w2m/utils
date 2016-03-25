package httputils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

func Test_Do(t *testing.T) {
	value := make(url.Values)
	value.Add("size", "8")
	value.Add("qr", `http://app.kuchuang.com`)
	//此二维码接口的url参数必需要size在前
	reqUrl := `http://apis.baidu.com/3023/qr/qrcode?size=8&qr=http%3A%2F%2Fapp.kuchuang.com`
	fmt.Println(reqUrl)
	req, _ := http.NewRequest("GET", reqUrl, nil)
	req.Header.Add("apikey", "f78db4344500d60edf5b971473489e61")

	resp, err := Do(req)
	if err != nil {
		fmt.Println("err")
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err")
		return
	}
	fmt.Println(resp.StatusCode)

	type DataStr struct {
		Url string `json:"url"`
	}
	fmt.Println(string(data))
	val := &DataStr{}
	json.Unmarshal(data, val)
	fmt.Println(val.Url)
}
