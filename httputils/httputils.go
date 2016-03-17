package httputils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//获取本机IP
func GetLocalAddr() string {
	conn, err := net.Dial("udp", "baidu.com:80")
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	defer conn.Close()
	return (strings.Split(conn.LocalAddr().String(), ":")[0])
}

var (
	client *http.Client
)

func init() {
	//设置超时时间
	client = &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*5)
				if err != nil {
					fmt.Println("dail timeout", err)
					return nil, err
				}
				return c, nil
			},
			MaxIdleConnsPerHost:   1024,
			ResponseHeaderTimeout: time.Second * 10,
		},
	}
}

func PostBytes(url string, body []byte) ([]byte, error) {
	r, err := client.Post(url, "application/octet-stream", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func PostForm(url string, values url.Values) ([]byte, error) {
	r, err := client.PostForm(url, values)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func PostFormOnly(url string, values url.Values) (*http.Response, error) {
	return client.PostForm(url, values)
}

func Get(url string) ([]byte, error) {
	r, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if r.StatusCode != 200 {
		return nil, fmt.Errorf("%s: %s", url, r.Status)
	}
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Delete(url string) error {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	resp, e := client.Do(req)
	if e != nil {
		return e
	}
	defer resp.Body.Close()
	if _, err := ioutil.ReadAll(resp.Body); err != nil {
		return err
	}
	return nil
}

func DownloadUrl(fileUrl string) (filename string, content []byte, e error) {
	response, err := client.Get(fileUrl)
	if err != nil {
		return "", nil, err
	}
	defer response.Body.Close()
	contentDisposition := response.Header["Content-Disposition"]
	if len(contentDisposition) > 0 {
		if strings.HasPrefix(contentDisposition[0], "filename=") {
			filename = contentDisposition[0][len("filename="):]
			filename = strings.Trim(filename, "\"")
		}
	}
	content, e = ioutil.ReadAll(response.Body)
	return
}

func Do(req *http.Request) (resp *http.Response, err error) {
	return client.Do(req)
}
