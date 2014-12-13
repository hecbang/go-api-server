package common

import (
	"io/ioutil"
	"os"
)

//判断一个文件或目录是否存在
//如果存在返回true，否则返回false
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil && !os.IsExist(err) {
		return false
	}
	return true
}

//从一个文件中读取内容
//@param string filename 要读取的文件
//@return []byte content 读取的文件内容
func Read(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}
