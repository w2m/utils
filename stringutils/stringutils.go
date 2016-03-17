package stringutils

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
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
func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}

//取子串
func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}
