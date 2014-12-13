package common

import (
	"encoding/json"
	//"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"strings"
	"time"
)

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
