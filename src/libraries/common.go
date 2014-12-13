package common

import (
	"math/rand"
	"runtime"
	"strings"
	"time"
)

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
