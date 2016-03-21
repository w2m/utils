package fileutils

import (
	"fmt"
	"testing"
)

func Test_ListDir(t *testing.T) {
	path := `./test`
	data, err := ListDir(path)
	if err != nil {
		fmt.Println("列出目录的文件失败, err:", err)
	} else {
		fmt.Println("目录下的文件列表")
		fmt.Println(data)
	}

	data, err = ListDir(path, "jpg")
	if err != nil {
		fmt.Println("列出目录的文件失败, err:", err)
	} else {
		fmt.Println("目录下的文件列表")
		fmt.Println(data)
	}
}

func Test_WalkDir(t *testing.T) {
	files, err := WalkDir(`./test`)
	if err != nil {
		fmt.Println("列出目录的文件失败, err:", err)
	} else {
		fmt.Println("目录下的文件列表")
		fmt.Println(files)
		if len(files) != 5 {
			t.Error("获取文件数量有误")
		}
	}

	files, err = WalkDir(`./test`, "txt")
	if err != nil {
		fmt.Println("列出目录的文件失败, err:", err)
	} else {
		fmt.Println("目录下的文件列表")
		fmt.Println(files)
		if len(files) != 4 {
			t.Error("获取文件数量有误")
		}
	}
}

func Test_Exist(t *testing.T) {
	isExist := Exist(`./test`)
	if !isExist {
		t.Error("判断目录或文件是否存在错误")
	}

	isExist = Exist(`./test/1.txt`)
	if !isExist {
		t.Error("判断目录或文件是否存在错误")
	}

	isExist = Exist(`./test/22.txt`)
	if isExist {
		t.Error("判断目录或文件是否存在错误")
	}

	isExist = Exist(`./test/test1`)
	if isExist {
		t.Error("判断目录或文件是否存在错误")
	}
}

func Test_ExistDir(t *testing.T) {
	isExist := ExistDir(`./test`)
	if !isExist {
		t.Error("判断目录是否存在错误")
	}

	isExist = ExistDir(`./test/1.txt`)
	if isExist {
		t.Error("判断目录是否存在错误")
	}

	isExist = ExistDir(`./test/22.txt`)
	if isExist {
		t.Error("判断目录是否存在错误")
	}

	isExist = ExistDir(`./test11`)
	if isExist {
		t.Error("判断目录是否存在错误")
	}
}

func Test_ExistFile(t *testing.T) {
	isExist := ExistFile(`./test`)
	if isExist {
		t.Error("判断文件是否存在错误")
	}

	isExist = ExistFile(`./test/1.txt`)
	if !isExist {
		t.Error("判断文件是否存在错误")
	}

	isExist = ExistFile(`./test/22.txt`)
	if isExist {
		t.Error("判断文件是否存在错误")
	}

	isExist = ExistFile(`./test11`)
	if isExist {
		t.Error("判断文件是否存在错误")
	}
}

func Test_IsSamePath(t *testing.T) {
	bFlag := IsSamePath(`./test`, `.\test`)
	if !bFlag {
		t.Error("判断是否为同目录失败")
	}

	bFlag = IsSamePath(`./test`, `./test/test`)
	if bFlag {
		t.Error("判断是否为同目录失败")
	}

	bFlag = IsSamePath(`E:\go\path\src\github.com\w2m\utils\fileutils\test`, `.\test`)
	if !bFlag {
		t.Error("判断是否为同目录失败")
	}

}
