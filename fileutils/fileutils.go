package fileutils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func ListDir(dirPth string, suffix ...string) ([]string, error) {
	files := make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}

		//后缀参数为空，表示不作过滤，返回所有类型的文件
		if len(suffix) == 0 {
			files = append(files, fi.Name())
			continue
		}

		//文件后缀过虑
		for _, v := range suffix {
			v = strings.ToUpper(v) //忽略后缀匹配的大小写

			if strings.HasSuffix(strings.ToUpper(fi.Name()), v) { //匹配文件
				files = append(files, fi.Name())
				break
			}
		}

	}
	return files, nil
}

//获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤。
func WalkDir(dirPth string, suffix ...string) (files []string, err error) {
	files = make([]string, 0, 30)

	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		if fi.IsDir() { // 忽略目录
			return nil
		}

		//后缀参数为空，表示不作过滤，返回所有类型的文件
		if len(suffix) == 0 {
			files = append(files, fi.Name())
			return nil
		}

		//检查是匹配后缀
		for _, v := range suffix {
			v = strings.ToUpper(v) //忽略后缀匹配的大小写

			if strings.HasSuffix(strings.ToUpper(fi.Name()), v) { //匹配文件
				files = append(files, fi.Name())
				return nil
			}
		}

		return nil
	})
	return files, err
}

// 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// 检查文件是否存在
// 如果由 filename 指定的文件存在则返回 true，否则返回 false(目录也返回错误)
func ExistFile(filename string) bool {
	fi, err := os.Stat(filename)
	if err == nil {
		if fi.IsDir() {
			return false
		}
		return true
	}
	return os.IsExist(err)
}

//是否存在目录
func ExistDir(filePath string) bool {
	fi, err := os.Stat(filePath)
	if err == nil {
		if fi.IsDir() {
			return true
		}
	}
	return false
}

//比较两个路径是否相等
func IsSamePath(srcPath, dstPath string) bool {
	absSrcPath, _ := filepath.Abs(srcPath)
	absDstPath, _ := filepath.Abs(dstPath)
	ret := strings.Compare(absSrcPath, absDstPath)
	if ret == 0 {
		return true
	} else {
		return false
	}
}

//获取配置文件名
func GetCfgName() string {
	//解析命令行参数

	cfgFileName := "kviewservice.ini"
	if len(os.Args) == 2 {
		if filepath.Ext(os.Args[1]) == ".ini" {
			cfgFileName = os.Args[1]
		}
	}

	return cfgFileName
}

//拷贝文件
func CopyFile(src, des string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Println(err)
	}
	defer srcFile.Close()

	desFile, err := os.Create(des)
	if err != nil {
		fmt.Println(err)
	}
	defer desFile.Close()

	return io.Copy(desFile, srcFile)
}

//读取文件
func ReadFile(path string) (body []byte, err error) {
	return ioutil.ReadFile(path)
}

//获取程序运行文件的绝对路径
func SelfPath() string {
	path, _ := filepath.Abs(os.Args[0])
	return path
}

//获取程序运行目录
func SelfDir() string {
	return filepath.Dir(SelfPath())
}

// SearchFile Search a file in paths.
// this is often used in search config file in /etc ~/
func SearchFile(filename string, paths ...string) (fullpath string, err error) {
	for _, path := range paths {
		if fullpath = filepath.Join(path, filename); Exist(fullpath) {
			return
		}
	}
	err = errors.New(fullpath + " not found in paths")
	return
}

// GrepFile like command grep -E
// for example: GrepFile(`^hello`, "hello.txt")
// \n is striped while read
func GrepFile(patten string, filename string) (lines []string, err error) {
	re, err := regexp.Compile(patten)
	if err != nil {
		return
	}

	fd, err := os.Open(filename)
	if err != nil {
		return
	}
	lines = make([]string, 0)
	reader := bufio.NewReader(fd)
	prefix := ""
	isLongLine := false
	for {
		byteLine, isPrefix, er := reader.ReadLine()
		if er != nil && er != io.EOF {
			return nil, er
		}
		if er == io.EOF {
			break
		}
		line := string(byteLine)
		if isPrefix {
			prefix += line
			continue
		} else {
			isLongLine = true
		}

		line = prefix + line
		if isLongLine {
			prefix = ""
		}
		if re.MatchString(line) {
			lines = append(lines, line)
		}
	}
	return lines, nil
}
