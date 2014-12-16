/**
 * @author yorkershi
 * @create on December 15, 2014
 */
package common

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

//根据分隔符(;,|)，将字符串转为一个字符串list
func StringToList(str string) []string {
	var list []string = make([]string, 0)
	if Empty(str) {
		return list
	}
	for _, delimiter := range []string{",", "|"} {
		str = strings.Replace(str, delimiter, ";", -1)
	}
	for _, v := range strings.Split(str, ";") {
		v = strings.TrimSpace(v)
		if !Empty(v) {
			list = append(list, v)
		}
	}
	return list
}

//判断一个数据是否为空，支持int, float, string, slice, array, map的判断
func Empty(value interface{}) bool {
	if value == nil {
		return true
	}
	switch reflect.TypeOf(value).Kind() {
	case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
		if reflect.ValueOf(value).Len() == 0 {
			return true
		} else {
			return false
		}
	}
	return false
}

//判断某一个值是否在列表(支持 slice, array, map)中
func InList(needle interface{}, haystack interface{}) bool {
	//interface{}和interface{}可以进行比较，但是interface{}不可进行遍历
	hayValue := reflect.ValueOf(haystack)
	switch reflect.TypeOf(haystack).Kind() {
	case reflect.Slice, reflect.Array:
		//slice, array类型
		for i := 0; i < hayValue.Len(); i++ {
			if hayValue.Index(i).Interface() == needle {
				return true
			}
		}
	case reflect.Map:
		//map类型
		var keys []reflect.Value = hayValue.MapKeys()
		for i := 0; i < len(keys); i++ {
			if hayValue.MapIndex(keys[i]).Interface() == needle {
				return true
			}
		}
	default:
		return false
	}
	return false
}

//加载一个JSON文件，并将JSON串解析到v的结构中
func LoadJson(filename string, v interface{}) error {
	bytes, err := Read(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, &v)
}

//获取memory使用量，单位byte
func MemoryGetUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}

//获取一个从[min, max]之间的随机数
func Rand(min, max int) int {
	max += 1
	time.Sleep(time.Nanosecond * 1)
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(max)
	if min > num {
		return min + num%(max-min)
	}
	return num
}

//获取一个当前时间，如Y-m-d H:i:s
func Date(format string) string {
	var conf = map[string]string{
		"Y": "2006",
		"m": "01",
		"d": "02",
		"H": "15",
		"i": "04",
		"s": "05",
	}
	for k, v := range conf {
		format = strings.Replace(format, k, v, -1)
	}
	return time.Now().Format(format)
}

//根据输入的值，获取一个经过md5编辑的key
//该函数主要应用于对一组数据的唯一性标识下
func BuildKeyMd5(args ...interface{}) string {
	var bytes []byte = make([]byte, 0, 20)
	if len(args) == 0 {
		bytes = append(bytes, []byte("")...)
	} else {
		for _, v := range args {
			//任何类型的值都按base-16解析成字符串并转化为[]byte类型
			bytes = append(bytes, []byte(fmt.Sprintf("%x", v))...)
		}
	}
	//按base-16形式解析并返回字符串
	return fmt.Sprintf("%x", md5.Sum(bytes))
}

//设置可使用的CPU核数，默认只使用一个CPU核心
func SetCPUNum() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

//获取一个浮点数，保留小数点后n位精度后的字符串值
func Round(val float64, precision int) string {
	return fmt.Sprintf("%."+strconv.Itoa(precision)+"f", val)
}
