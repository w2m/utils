package stringutils

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
	"math"
	r "math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"
)

func GetAfter(src string, find string) (string, bool) {
	index := strings.Index(src, find)
	if index > -1 && index+len(find) < len(src) {
		return src[(index + len(find)):], true
	}
	return src, false
}

func GetBefore(src string, find string) (string, bool) {
	index := strings.Index(src, find)
	if index > -1 && index < len(src) {
		return src[:index], true
	}
	return src, false
}

func UnicodeDecode(text string) string {
	regex, err := regexp.Compile(`(\\u[a-fA-F0-9]{4})`)
	if err != nil {
		return text
	}

	text = regex.ReplaceAllStringFunc(text, func(match string) string {
		_txt := match[2:]
		char, err := strconv.ParseInt(_txt, 16, 32)
		if err != nil {
			return match
		}
		return string(rune(int(char)))
	})

	regex, err = regexp.Compile(`(&#[\d]{2,6})`)
	if err != nil {
		return text
	}

	text = regex.ReplaceAllStringFunc(text, func(match string) string {
		_txt := match[2:]
		char, err := strconv.ParseInt(_txt, 10, 32)
		if err != nil {
			return match
		}
		return string(rune(int(char)))
	})

	return text
}

func UnicodeEncode(str string) (js, html string) {
	rs := []rune(str)
	js = ""
	html = ""
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			js += string(r)
			html += string(r)
		} else {
			js += `\u` + strconv.FormatInt(int64(rint), 16) // json
			html += `&#` + strconv.Itoa(int(r)) + ";"       // 网页
		}
	}
	fmt.Printf("JSON: %s\n", js)
	fmt.Printf("HTML: %s\n", html)
	return
}

//比较是否在切片中存在，不区分大小写
func ExistSliceFold(srcSlice []string, elem string) bool {
	for _, v := range srcSlice {
		if strings.EqualFold(v, elem) {
			return true
		}
	}

	return false
}

//比较是否在切片中存在,区分大小写
func ExistSlice(srcSlice []string, elem string) bool {
	for _, v := range srcSlice {
		if v == elem {
			return true
		}
	}

	return false
}

//字符串转换uint16
func StringToUTF16(s string) []uint16 {
	return utf16.Encode([]rune(s + "\x00"))
}

func Float64IsZero(s float64) bool {
	if math.Abs(s) < 0.0001 {
		return true
	}
	return false
}

//取子串
func Substr(s string, start, length int) string {
	bt := []rune(s)
	if start < 0 {
		start = 0
	}
	if start > len(bt) {
		start = start % len(bt)
	}
	var end int
	if (start + length) > (len(bt) - 1) {
		end = len(bt)
	} else {
		end = start + length
	}
	return string(bt[start:end])
}

//删除slice中的元素
func RemoveSliceElement(val interface{}, index int) interface{} {

	if reflect.TypeOf(val).Kind() != reflect.Slice {
		fmt.Println("val类型非slice")
		return nil
	}

	s := reflect.ValueOf(val)
	if index < 0 || index >= s.Len() {
		fmt.Println("传入参数有误")
		return nil
	}

	prev := s.Index(index)
	for i := index + 1; i < s.Len(); i++ {
		value := s.Index(i)
		prev.Set(value)
		prev = value
	}

	return s.Slice(0, s.Len()-1).Interface()
}

// RandomCreateBytes generate random []byte by specify chars.
func RandomCreateBytes(n int, alphabets ...byte) []byte {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	var randby bool
	if num, err := rand.Read(bytes); num != n || err != nil {
		r.Seed(time.Now().UnixNano())
		randby = true
	}
	for i, b := range bytes {
		if len(alphabets) == 0 {
			if randby {
				bytes[i] = alphanum[r.Intn(len(alphanum))]
			} else {
				bytes[i] = alphanum[b%byte(len(alphanum))]
			}
		} else {
			if randby {
				bytes[i] = alphabets[r.Intn(len(alphabets))]
			} else {
				bytes[i] = alphabets[b%byte(len(alphabets))]
			}
		}
	}
	return bytes
}

func SumMd5(txtInput string) string {
	h := md5.New()
	io.WriteString(h, txtInput)
	return fmt.Sprintf("%x", h.Sum(nil))
}
