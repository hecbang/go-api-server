package testing

import (
	"fmt"
	"libraries/common"
	"log"
)

//DB并发测试
//n 测试总次数
//c 并发量
func Concurrence(groupid int64, schemaname string, targetdbschema string, n int, c int) {
	//先清空写入的目标数据库
	db := common.NewMySqlInstance(targetdbschema)
	_, err := db.UDExec("truncate table target")
	if err != nil {
		log.Fatalln(err.Error())
	}

	timer := common.NewTimer()
	timer.Start()

	//每个线程执行多少次
	segs := n / c

	var chs []chan int = make([]chan int, c)
	var cycleN int
	for i := 0; i < c; i++ {
		cycleN = segs
		if i == c-1 {
			cycleN += n % c
		}
		chs[i] = make(chan int)
		go func(cycleN int, ch chan int) {
			db := common.NewMySqlInstance(targetdbschema)
			date := common.Date("Y-m-d H:i:s")
			data := map[string]interface{}{
				"Num":     98,
				"String":  "hello,yorkershi",
				"LogTime": date,
			}
			for cyc := 0; cyc < cycleN; cyc++ {
				_, err := db.Insert("target", data)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
			ch <- 1
		}(cycleN, chs[i])
	}

	fmt.Println("wait...")
	for _, ch := range chs {
		<-ch
	}

	elapse := timer.Elapse("ms")

	fmt.Println("finished once concurrence, elapse ", elapse, "ms")

	result := map[string]interface{}{
		"GroupId":     groupid,
		"Name":        schemaname,
		"Total":       n,
		"Concurrence": c,
		"ElapseTime":  elapse,
		"QPS":         common.Round(float64(n)/(float64(elapse)/1000), 2), //每秒处理请求数
		"TPQ":         common.Round(float64(elapse)/float64(n), 2),        //平均每个请求用时多少ms
		"LogTime":     common.Date("Y-m-d H:i:s"),
	}
	tdb := common.NewMySqlInstance("testdata")
	tdb.Insert("db", result)
}
