package ziputils

import (
	"archive/zip"
	"bytes"
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
		//创建目录
		err = os.MkdirAll(filepath.Dir(filepath.Join(dstPath, file.Name)), 0666)
		if err != nil && err != os.ErrExist {
			fmt.Println("创建目录失败, path:", filepath.Dir(filepath.Join(dstPath, file.Name)), ", err:", err)
			return err
		}
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
	return nil

}
