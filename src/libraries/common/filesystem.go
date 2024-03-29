/**
 * @author yorkershi
 * @create on December 13, 2014
 */
package common

import (
	"bufio"
	"io"
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

//将一个文件中的内容读取到[]string中，并返回
func ReadLine(filename string) ([]string, error) {
	var retval []string = make([]string, 0)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var bufrd *bufio.Reader = bufio.NewReader(file)

	for {
		line, err := bufrd.ReadString('\n')
		if !Empty(line) {
			retval = append(retval, line)
		} else {
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			retval = append(retval, line)
		}

	}
	return retval, nil
}
