package testing

import (
	"const/path"
	"fmt"
	"libraries/common"
	"log"
)

func DatabaseConcurrence() {
	jq := common.NewJsonQuery(path.CONFIG_PATH + "testing" + path.DS + "db_concurrence.json")

	db := common.NewMySqlInstance("testdata")

	name := jq.String("group", "name")
	parameter := jq.String("group", "parameter")

	data := map[string]interface{}{
		"Name":              name,
		"SettingParameters": parameter,
		"LogTime":           common.Date("Y-m-d H:i:s"),
	}

	schemaName := jq.String("schema", "name")
	targetSchemaDb := jq.String("schema", "db")
	amount := jq.Int("schema", "amount")
	start := jq.Int("schema", "start")
	offset := jq.Int("schema", "offset")
	max := jq.Int("schema", "max")

	//处理掉重复的记录
	sql := "select db.Id, db.GroupId from db inner join db_group on db.GroupId=db_group.Id where db.Name=? and db_group.Name=?"
	result, err := db.GetRow(sql, schemaName, name)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if !common.Empty(result) {
		db.Delete("db_group", map[string]interface{}{"Id": result["GroupId"]})
		db.Delete("db", map[string]interface{}{"Id": result["Id"]})
	}

	lastid, err := db.Insert("db_group", data)
	if err != nil {
		panic(err.Error())
	}

	c := start
	fmt.Println("max concurrence is ", max)
	for c < max {
		fmt.Println("Now, concurrence = ", c)
		dbconcurrence(lastid, schemaName, targetSchemaDb, amount, c)
		c = c + offset - (c % offset)
	}
}

//DB并发测试
//n 测试总次数
//c 并发量
func dbconcurrence(groupid int64, schemaname string, targetdbschema string, n int, c int) {
	if c > n {
		panic("error: c>n")
	}

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
