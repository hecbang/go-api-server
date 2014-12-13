package common

import (
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
