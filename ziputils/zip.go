package ziputils

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//压缩文件
func ZipFile(srcPath, dstFileName string, bFlag bool) error {
	buf := new(bytes.Buffer)
	myzip := zip.NewWriter(buf)

	//fmt.Println("srcPath:", srcPath)
	//遍历目录下的所有文件，将文件写入到压缩包中
	err := filepath.Walk(srcPath, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return filepath.SkipDir
		}

		//fmt.Println("path:", path)

		//非目录才写入到压缩包
		if !fi.IsDir() {

			header, err := zip.FileInfoHeader(fi)
			if err != nil {
				return filepath.SkipDir
			}
			//判断是否需要把顶层目录写入到压缩包
			if bFlag {
				refFile, _ := filepath.Rel(filepath.Dir(srcPath), path)
				header.Name = strings.SplitN(refFile, `\`, 2)[1]
			} else {
				header.Name, _ = filepath.Rel(filepath.Dir(srcPath), path)
			}
			//fmt.Println("header.Name:", header.Name)

			//创建一个头
			w, err := myzip.CreateHeader(header)
			if err != nil {
				return err
			}

			fileData, err := ioutil.ReadFile(path)
			if err != nil {
				return filepath.SkipDir
			}
			//写入文件信息
			w.Write(fileData)
		}
		return nil
	})

	myzip.Close()

	//先检查目录存在不存在
	os.MkdirAll(filepath.Dir(dstFileName), 0666)
	// 建立zip文件
	retFile, err := os.Create(dstFileName)
	if err != nil {
		return err
	}
	defer retFile.Close()

	// 将buf中的数据写入文件
	_, err = buf.WriteTo(retFile)
	if err != nil {
		return err
	}
	return nil
}

//解压文件
func UnZip(srcFile, dstPath string) error {

	//创建一个目录
	err := os.MkdirAll(dstPath, 0666)
	if err != nil {
		fmt.Println("创建目录失败, path:", dstPath, ", err:", err)
		return err
	}

	//读取zip文件
	cf, err := zip.OpenReader(srcFile)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer cf.Close()

	for _, file := range cf.File {

		rc, err := file.Open()
		if err != nil {
			fmt.Println("open file failed, err:", err)
			return err
		}

		if file.FileInfo().IsDir() {
			//创建目录
			err = os.MkdirAll(filepath.Join(dstPath, file.Name), 0666)
			if err != nil && err != os.ErrExist {
				fmt.Println("创建目录失败, path:", filepath.Dir(filepath.Join(dstPath, file.Name)), ", err:", err)
				return err
			}
		} else {

			//创建文件所在目录
			err = os.MkdirAll(filepath.Dir(filepath.Join(dstPath, file.Name)), 0666)
			if err != nil && err != os.ErrExist {
				fmt.Println("创建目录失败, path:", filepath.Dir(filepath.Join(dstPath, file.Name)), ", err:", err)
				return err
			}
			//创建文件
			f, err := os.Create(filepath.Join(dstPath, file.Name))
			if err != nil {
				fmt.Println("create file failed, err:", err)
				return err
			}
			defer f.Close()
			_, err = io.Copy(f, rc)
			if err != nil {
				fmt.Println("copy file failed, err:", err)
				return err
			}
		}

	}
	return nil
}

func WriteComment(fileName string, comment ...string) error {
	f, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("打开文件失败, err:", err)
		return err
	}
	defer f.Close()

	//定位到comment长度字段，写入数据长度
	f.Seek(0, os.SEEK_END)
	commentList := make([]zipComment, 0, len(comment))
	item := zipComment{}
	for _, v := range comment {
		item.Data = v
		item.getLen()
		commentList = append(commentList, item)
	}
	byteComment := pack(commentList...)
	fmt.Println("comment bytes:", byteComment)

	commentLen := uint16(len(byteComment))
	fmt.Println("commentLen:", commentLen)

	binary.Write(f, binary.BigEndian, commentLen)
	f.Seek(0, os.SEEK_END)

	num, err := f.Write(byteComment)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = binary.Write(f, binary.BigEndian, uint16(num))
	if err != nil {
		fmt.Println("err:", err)
		return err
	}
	fmt.Println("write num:", num)
	return nil
}

func ReadComment(fileName string) []string {
	f, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("打开文件失败, err:", err)
		return make([]string, 0)
	}
	defer f.Close()

	f.Seek(-2, os.SEEK_END)
	var commentLen uint16
	err = binary.Read(f, binary.BigEndian, &commentLen)
	if err != nil {
		fmt.Println("获取comment长度错误, err:", err)
		return make([]string, 0)
	}

	//	fmt.Println("commentLen:", commentLen)

	seekLen := int64(0) - int64(commentLen) - int64(2)
	fmt.Println("seekLen:", seekLen)
	f.Seek(seekLen, os.SEEK_END)
	comment := make([]byte, commentLen)
	num, err := f.Read(comment)
	if err != nil {
		fmt.Println("读取comment数据失败, err:", err)
		return make([]string, 0)
	}

	fmt.Println("comment read len:", num)

	commentData := unPack(comment)
	fmt.Println("commentData:", commentData)

	retData := make([]string, 0, len(commentData))
	for _, v := range commentData {
		retData = append(retData, v.Data)
	}
	return retData
}

type zipComment struct {
	Data string
	Len  uint16
}

func (this *zipComment) getLen() {
	this.Len = uint16(len([]byte(this.Data)))
}

func pack(data ...zipComment) []byte {
	buf := new(bytes.Buffer)

	for _, v := range data {
		err := binary.Write(buf, binary.LittleEndian, v.Len)
		if err != nil {
			fmt.Println(err)
			return buf.Bytes()
		}
		buf.WriteString(v.Data)
	}

	return buf.Bytes()
}

func unPack(data []byte) []zipComment {
	buf := bytes.NewReader(data)

	var retData []zipComment

	buf.Seek(0, os.SEEK_SET)
	for {
		var val zipComment
		err := binary.Read(buf, binary.LittleEndian, &val.Len)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}

		item := make([]byte, val.Len)
		num, err := buf.Read(item)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("err:", err)
		}

		fmt.Println("num:", num)
		val.Data = string(item)
		retData = append(retData, val)
	}
	return retData
}
