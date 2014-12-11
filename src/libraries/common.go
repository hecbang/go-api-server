package clib

import (
	"math/rand"
	"runtime"
	"time"
)

type Common struct {
}

/**
 * 获取一个min-max之间的随机数
 * @param int min
 * @param int max
 * @param int randNumber
 */
func (c *Common) Random(min, max int) int {
	//延迟1纳秒，避免不间断循环时取随机数重复的情况
	time.Sleep(1 * time.Nanosecond)
	seed := time.Now().Nanosecond()
	rand.Seed(int64(seed))
	num := rand.Intn(max)
	if num < min {
		return min + num%(max-min)
	}
	return num
}

func (c *Common) InUseMemo() uint64 {
	var mpf runtime.MemStats
	runtime.ReadMemStats(&mpf)
	return mpf.Alloc
}
