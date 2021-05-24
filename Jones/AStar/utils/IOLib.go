package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

func ReadBytes(reader io.Reader, num int64) ([]byte, error) {
	p := make([]byte, num)
	n, err := reader.Read(p)
	if n > 0 {
		return p[:n], nil
	}
	return p, err
}

// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func MkDir(dir string) {
	exist, err := PathExists(dir)
	if err != nil {
		Logger.Error("get dir error![%v]", err)
		return
	}

	if exist {
		Logger.Info("has dir![%v]", dir)
	} else {
		Logger.Info("no dir![%v]", dir)
		// 创建文件夹
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			Logger.Error("mkdir failed![%v]", err)
		} else {
			Logger.Info("mkdir success!")
		}
	}
}

func PrintStackTrace() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	Logger.Error("==> %s", string(buf[:n]))
}


func WriteFile(base string, fileName string, content *string) (string, error){
	MkDir(base)
	toFile := fmt.Sprintf("%s/%s",base,fileName)
	return WriteFileFullPath(toFile, content)
}

func WriteFileFullPath(toFile string, content *string) (string, error){
	parentDir := getParentDirectory(toFile)
	MkDir(parentDir)
	err := ioutil.WriteFile(toFile, []byte(*content), 0777)
	if err != nil {
		Logger.Error(err.Error())
		panic(err.Error())
	}
	return toFile,err
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}
