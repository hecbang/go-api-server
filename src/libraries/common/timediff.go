package common

import (
	"time"
)

type TimeDiff struct {
	start int64
}

//记录并设置当前时间为开始时间点，单位nanosecond
func (this *TimeDiff) Start() {
	this.start = time.Now().UnixNano()
}

//获取从之前设置的开始时间点到当前时间用时，参数支持单位s, ms, us, ns
func (this *TimeDiff) Elapse(unit string) int64 {
	now := time.Now().UnixNano()
	switch unit {
	case "ns":
		return now - this.start
	case "us":
		return int64((now - this.start) / 1000)
	case "ms":
		return int64((now - this.start) / 1000000)
	case "s":
		return int64((now - this.start) / 1000000000)
	}
	panic("parameter invalid.")
}
